package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/de4et/office-mail/pkg/logger"
	"github.com/de4et/office-mail/services/delivery/internal/domain"
	"github.com/de4et/office-mail/services/delivery/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	sessionTimeout = 8000
)

var (
	tr = otel.Tracer("delivery")
)

type KafkaTaskConsumer struct {
	consumer *kafka.Consumer
	topic    string
}

func MustGetKafkaTaskConsumer(config Config) *KafkaTaskConsumer {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(config.Addresses, ","),
		"group.id":                 config.ConsumerGroup,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": true,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"auto.offset.reset":        "earliest",
	}

	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		panic(err)
	}

	c.Subscribe(config.MailTopic, func(c *kafka.Consumer, ev kafka.Event) error {
		switch e := ev.(type) {
		case kafka.AssignedPartitions:
			log.Printf("assigned: %v", e.Partitions)
			c.Assign(e.Partitions)
		case kafka.RevokedPartitions:
			log.Printf("revoked: %v", e.Partitions)
			c.Unassign()
		default:
			log.Println(ev)
		}
		return nil
	})

	return &KafkaTaskConsumer{
		consumer: c,
		topic:    config.MailTopic,
	}
}

type payload struct {
	TaskID int    `json:"task_id"`
	MailID int    `json:"mail_id"`
	From   string `json:"from"`
	To     string `json:"to"`
	Body   string `json:"body"`
}

func serializeFromMessage(msg string) (domain.Mail, domain.OutboxTask, error) {
	var payload payload
	err := json.Unmarshal([]byte(msg), &payload)
	if err != nil {
		return domain.Mail{}, domain.OutboxTask{}, err
	}

	mail := domain.Mail{
		ID:   payload.MailID,
		To:   domain.Address(payload.To),
		From: domain.Address(payload.From),
		Body: payload.Body,
	}

	task := domain.OutboxTask{
		ID:          payload.TaskID,
		AggregateID: payload.MailID,
	}

	return mail, task, nil
}

func (c *KafkaTaskConsumer) GetNextMailTask(ctx context.Context) (domain.Mail, domain.OutboxTask, error) {
	kafkaMsg, err := c.consumer.ReadMessage(-1)
	if err != nil {
		return domain.Mail{}, domain.OutboxTask{}, err
	}

	mail, task, err := serializeFromMessage(string(kafkaMsg.Value))
	if err != nil {
		return mail, task, err
	}

	carrier := propagation.MapCarrier{}

	for _, h := range kafkaMsg.Headers {
		carrier[h.Key] = string(h.Value)
	}
	fmt.Printf("carrier: %v\n", carrier)

	ctx = otel.GetTextMapPropagator().
		Extract(ctx, carrier)

	span := trace.SpanFromContext(ctx)
	fmt.Printf("span.SpanContext().TraceID(): %v\n", span.SpanContext().TraceID())

	span.SetAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination", c.topic),
		attribute.Int("outbox.task_id", task.ID),
	)

	return mail, task, nil
}

func (c *KafkaTaskConsumer) Consume(ctx context.Context, handler usecase.HandlerFunc) {
	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			slog.ErrorContext(ctx, "couldn't retrive task")
		}

		mail, task, err := serializeFromMessage(string(msg.Value))
		if err != nil {
			continue
		}

		carrier := propagation.MapCarrier{}
		for _, h := range msg.Headers {
			carrier[h.Key] = string(h.Value)
		}
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)

		ctx, span := tr.Start(ctx, "delivery.process_mail")
		span.SetAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", c.topic),
			attribute.Int("outbox.task_id", task.ID),
		)

		ctx = logger.WithContext(ctx, "trace_id", span.SpanContext().TraceID())

		err = handler(ctx, mail, task)
		if err != nil {
			span.RecordError(err)
		}

		span.End()
	}
}

func (c *KafkaTaskConsumer) Close() {
	c.consumer.Commit()
	c.consumer.Close()
}
