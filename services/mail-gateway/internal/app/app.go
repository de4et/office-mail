package app

import (
	"context"
	"fmt"
	"log/slog"

	pg "github.com/de4et/office-mail/pkg/postgres"
	"github.com/de4et/office-mail/services/mail-gateway/internal/adapters/postgres"
	"github.com/de4et/office-mail/services/mail-gateway/internal/controller"
	"github.com/de4et/office-mail/services/mail-gateway/internal/usecase"
)

func Run(config Config) {
	ctx := context.Background()

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

	fmt.Printf("config: %v\n", config)
	listenAddr := fmt.Sprintf(":%d", config.PORT)
	controller.Listen(listenAddr)
}
