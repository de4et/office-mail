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
	"github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
)

//go:embed queries/get_first_available_task.sql
var getFirstAvailableTaskQuery string

//go:embed queries/mark_outbox_task_as_done.sql
var markTaskAsDoneQuery string

type PostgresqlOutboxRepository struct {
	client *postgres.TxClient
	tr     trace.Tracer
}

func NewPostgresqlOutboxRepository(client *sqlx.DB, tr trace.Tracer) *PostgresqlOutboxRepository {
	return &PostgresqlOutboxRepository{
		client: postgres.NewTxClient(client),
		tr:     tr,
	}
}

func (rep *PostgresqlOutboxRepository) GetFirstAvailableTasks(ctx context.Context, limit int) ([]domain.OutboxTask, error) {
	var resp []dto.GetFirstAvailableTaskResponse
	err := rep.client.SelectContext(
		ctx,
		&resp,
		getFirstAvailableTaskQuery,
		limit,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, usecase.ErrNoTasks
		}
		return nil, err
	}
	if len(resp) == 0 {
		return nil, usecase.ErrNoTasks
	}

	tasks := make([]domain.OutboxTask, len(resp))
	for i := range resp {
		tasks[i] = domain.OutboxTask{
			ID:           resp[i].ID,
			AggregateID:  resp[i].AggregateID,
			TraceContext: resp[i].TraceContext,
		}
	}

	return tasks, nil
}

func (rep *PostgresqlOutboxRepository) MarkTaskAsDone(ctx context.Context, task domain.OutboxTask) error {
	ctx, span := rep.tr.Start(ctx, "outboxRepository.MarkTaskAsDone")
	defer span.End()

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

func (rep *PostgresqlOutboxRepository) MarkTasksAsDone(ctx context.Context, tasks []domain.OutboxTask) error {
	ctx, span := rep.tr.Start(ctx, "outboxRepository.MarkTasksAsDone")
	defer span.End()

	tasksIDs := make([]int, len(tasks))
	for i := range tasks {
		tasksIDs[i] = tasks[i].ID
	}

	_, err := rep.client.ExecContext(
		ctx,
		markTaskAsDoneQuery,
		pq.Array(tasksIDs),
	)
	if err != nil {
		return err
	}

	return nil
}
