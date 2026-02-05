package postgres

import (
	"context"
	_ "embed"

	"github.com/de4et/office-mail/pkg/postgres"
	"github.com/de4et/office-mail/services/mail-gateway/internal/domain"
	"github.com/jmoiron/sqlx"
)

//go:embed queries/create_outbox_delivery_task.sql
var createOutboxDeliveryTaskQuery string

type PostgresqlOutboxRepository struct {
	client *postgres.TxClient
}

func NewPostgresqlOutboxRepository(client *sqlx.DB) *PostgresqlOutboxRepository {
	return &PostgresqlOutboxRepository{
		client: postgres.NewTxClient(client),
	}
}

func (rep *PostgresqlOutboxRepository) CreateOutboxDeliveryTask(ctx context.Context, mail domain.Mail) error {
	_, err := rep.client.ExecContext(
		ctx,
		createOutboxDeliveryTaskQuery,
		mail.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
