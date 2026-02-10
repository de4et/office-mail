package app

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/de4et/office-mail/pkg/logger"
	"github.com/de4et/office-mail/pkg/tracer"
	"github.com/de4et/office-mail/services/delivery/internal/adapters/kafka"
	"github.com/de4et/office-mail/services/delivery/internal/adapters/redis"
	"github.com/de4et/office-mail/services/delivery/internal/adapters/stub"
	"github.com/de4et/office-mail/services/delivery/internal/usecase"
)

func Run(config Config) {
	ctx := context.Background()
	logger.SetupLog("", slog.LevelInfo)

	tp, err := tracer.NewTraceProvider("delivery")
	if err != nil {
		panic(err)
	}
	defer tp.Shutdown(ctx)

	kafkaConsumer := kafka.MustGetKafkaTaskConsumer(kafka.Config{
		Addresses:     strings.Split(config.KAFKA_ADDRESSES, ","),
		MailTopic:     config.KAFKA_MAIL_TOPIC,
		ConsumerGroup: config.KAFKA_CONSUMER_GROUP,
	})

	mailDelivery := stub.NewStubMailDelivery()
	spamChecker := stub.NewStubSpamChecker()

	idempotenceChecker := redis.MustGetRedisIdempotenceChecker(redis.Config{
		Addr:     fmt.Sprintf("%s:%s", config.REDIS_HOST, config.REDIS_PORT),
		Password: config.REDIS_PASSWORD,
	})

	slog.InfoContext(ctx, "starting!")

	outboxWorkerUsecase := usecase.NewOutboxWorkerUsecase(kafkaConsumer, spamChecker, idempotenceChecker, mailDelivery)
	outboxWorkerUsecase.Run(ctx)
}
