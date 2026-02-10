package postgres

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/de4et/office-mail/pkg/postgres"
	"github.com/de4et/office-mail/services/mail-worker/internal/adapters/postgres/dto"
	"github.com/de4et/office-mail/services/mail-worker/internal/domain"
	"github.com/de4et/office-mail/services/mail-worker/internal/usecase"
	"github.com/jmoiron/sqlx"
)

//go:embed queries/get_first_available_task.sql
var getFirstAvailableTaskQuery string

//go:embed queries/mark_outbox_task_as_done.sql
var markTaskAsDoneQuery string

type PostgresqlOutboxRepository struct {
	client *postgres.TxClient
}

func NewPostgresqlOutboxRepository(client *sqlx.DB) *PostgresqlOutboxRepository {
	return &PostgresqlOutboxRepository{
		client: postgres.NewTxClient(client),
	}
}

func (rep *PostgresqlOutboxRepository) GetFirstAvailableTask(ctx context.Context) (domain.OutboxTask, error) {
	var resp dto.GetFirstAvailableTaskResponse
	err := rep.client.GetContext(
		ctx,
		&resp,
		getFirstAvailableTaskQuery,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.OutboxTask{}, usecase.ErrNoTasks
		}
		return domain.OutboxTask{}, err
	}

	return domain.OutboxTask{
		ID:           resp.ID,
		AggregateID:  resp.AggregateID,
		TraceContext: resp.TraceContext,
	}, nil
}

func (rep *PostgresqlOutboxRepository) MarkTaskAsDone(ctx context.Context, task domain.OutboxTask) error {
	_, err := rep.client.ExecContext(
		ctx,
		markTaskAsDoneQuery,
		task.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
