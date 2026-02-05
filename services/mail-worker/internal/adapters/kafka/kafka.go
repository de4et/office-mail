package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/de4et/office-mail/services/mail-worker/internal/domain"
)

const (
	flushTimeout = 5000
)

var errUnknownType = errors.New("unknown event type")

type KafkaTaskPublisher struct {
	producer *kafka.Producer
	topic    string
}

func MustGetKafkaTaskPublisher(config Config) *KafkaTaskPublisher {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(config.Addresses, ","),
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		panic(err)
	}

	return &KafkaTaskPublisher{
		producer: p,
		topic:    config.MailTopic,
	}
}

func (p *KafkaTaskPublisher) PublishMailTask(ctx context.Context, task domain.OutboxTask, mail domain.Mail) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value:     []byte(p.prepareMessage(task, mail)),
		Timestamp: time.Now(),
	}

	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return nil
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return ev
	default:
		return errUnknownType
	}
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
