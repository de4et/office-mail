package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/de4et/office-mail/services/mail-worker/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	flushTimeout = 5000
)

var errUnknownType = errors.New("unknown event type")

type KafkaTaskPublisher struct {
	producer *kafka.Producer
	topic    string
	tr       trace.Tracer
}

func MustGetKafkaTaskPublisher(config Config, tr trace.Tracer) *KafkaTaskPublisher {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(config.Addresses, ","),
		// "linger.ms":         1,
		// "acks":              1,
		// "compression.type":  "none",
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		panic(err)
	}

	return &KafkaTaskPublisher{
		producer: p,
		topic:    config.MailTopic,
		tr:       tr,
	}
}

func (p *KafkaTaskPublisher) PublishMailTask(ctx context.Context, task domain.OutboxTask, mail domain.Mail) error {
	ctx, span := p.tr.Start(ctx, "taskPublisher.PublishMailTask")
	defer span.End()

	headers := make([]kafka.Header, 0)

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	for k, v := range carrier {
		headers = append(headers, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value:     []byte(p.prepareMessage(task, mail)),
		Timestamp: time.Now(),
		Headers:   headers,
	}

	// kafkaChan := make(chan kafka.Event)
	// if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
	if err := p.producer.Produce(kafkaMsg, nil); err != nil {
		return err
	}

	// e := <-kafkaChan
	// switch ev := e.(type) {
	// case *kafka.Message:
	// 	return nil
	// case kafka.Error:
	// 	return ev
	// default:
	// 	return errUnknownType
	// }
	return nil
}

type payload struct {
	TaskID int    `json:"task_id"`
	MailID int    `json:"mail_id"`
	From   string `json:"from"`
	To     string `json:"to"`
	Body   string `json:"body"`
}

func (p *KafkaTaskPublisher) prepareMessage(task domain.OutboxTask, mail domain.Mail) []byte {
	payload := payload{
		TaskID: task.ID,
		MailID: task.AggregateID,
		From:   string(mail.From),
		To:     string(mail.To),
		Body:   mail.Body,
	}

	b, _ := json.Marshal(&payload)
	return b
}

func (p *KafkaTaskPublisher) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
