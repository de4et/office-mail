package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/de4et/office-mail/pkg/logger"
	pg "github.com/de4et/office-mail/pkg/postgres"
	"github.com/de4et/office-mail/pkg/tracer"
	"github.com/de4et/office-mail/services/mail-gateway/internal/adapters/postgres"
	"github.com/de4et/office-mail/services/mail-gateway/internal/controller"
	"github.com/de4et/office-mail/services/mail-gateway/internal/usecase"
)

func Run(config Config) {
	ctx := context.Background()

	logger.SetupLog("", slog.LevelDebug)

	tp, err := tracer.NewTraceProvider(os.Getenv("SERVICE_NAME"))
	if err != nil {
		panic(err)
	}
	defer tp.Shutdown(ctx)

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
	mailRep := postgres.NewPostgresqlMailRepository(pgClient)
	outboxRep := postgres.NewPostgresqlOutboxRepository(pgClient)

	sendUC := usecase.NewSendMailUsecase(transactor, mailRep, outboxRep)
	controller := controller.SetupRoutes(ctx, sendUC)

	listenAddr := fmt.Sprintf(":%d", config.PORT)
	controller.Listen(listenAddr)
}
