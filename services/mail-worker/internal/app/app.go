package app

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/de4et/office-mail/pkg/logger"
	pg "github.com/de4et/office-mail/pkg/postgres"
	"github.com/de4et/office-mail/pkg/tracer"
	"github.com/de4et/office-mail/services/mail-worker/internal/adapters/kafka"
	"github.com/de4et/office-mail/services/mail-worker/internal/adapters/postgres"
	"github.com/de4et/office-mail/services/mail-worker/internal/usecase"
	"go.opentelemetry.io/otel"
)

func Run(config Config) {
	ctx := context.Background()
	logger.SetupLog("", slog.LevelInfo)

	tp, err := tracer.NewTraceProvider(os.Getenv("SERVICE_NAME"))
	if err != nil {
		panic(err)
	}
	defer tp.Shutdown(ctx)

	tr := otel.Tracer("mail-worker")

	pgClient := pg.MustGetPostgresqlClient(pg.Config{
		Host:     config.DB_HOST,
		Port:     config.DB_PORT,
		Username: config.DB_USERNAME,
		Password: config.DB_PASSWORD,
		DbName:   config.DB_DATABASE,
	})

	if err := pgClient.Ping(); err != nil {
		panic("couldn't ping pgClient")
	}

	slog.InfoContext(ctx, "Successfully set up pg client!")

	transactor := pg.NewPostgresqlTransactor(pgClient)
	mailRep := postgres.NewPostgresqlMailRepository(pgClient, tr)
	outboxRep := postgres.NewPostgresqlOutboxRepository(pgClient, tr)
	kafkaPublisher := kafka.MustGetKafkaTaskPublisher(kafka.Config{
		Addresses: strings.Split(config.KAFKA_ADDRESSES, ","),
		MailTopic: config.KAFKA_MAIL_TOPIC,
	}, tr)

	outboxWorkerUsecase := usecase.NewOutboxWorkerUsecase(transactor, mailRep, outboxRep, kafkaPublisher, tr)
	outboxWorkerUsecase.Run(ctx)
}
