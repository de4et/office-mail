package postgres

import (
	"context"
	_ "embed"
	"encoding/json"

	"github.com/de4et/office-mail/pkg/postgres"
	"github.com/de4et/office-mail/services/mail-gateway/internal/domain"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	traceCtx, err := json.Marshal(carrier)
	if err != nil {
		return err
	}

	_, err = rep.client.ExecContext(
		ctx,
		createOutboxDeliveryTaskQuery,
		mail.ID,
		traceCtx,
	)
	if err != nil {
		return err
	}

	return nil
}
