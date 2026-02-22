package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/de4et/office-mail/pkg/logger"
	"github.com/de4et/office-mail/services/mail-worker/internal/domain"
)

var (
	ErrNoTasks = errors.New("no tasks")
)

const (
	checkInterval   = time.Millisecond * 10
	tasksBatchLimit = 100
)

type mailRepository interface {
	GetMailsByIDs(context.Context, []int) (map[int]domain.Mail, error)
}

type outboxRepository interface {
	GetFirstAvailableTasks(context.Context, int) ([]domain.OutboxTask, error)
	MarkTasksAsDone(context.Context, []domain.OutboxTask) error
}

type mailTaskPublisher interface {
	PublishMailTask(context.Context, domain.OutboxTask, domain.Mail) error
}

type OutboxWorkerUsecase struct {
	tx        transactor
	mailRep   mailRepository
	outboxRep outboxRepository
	mailTP    mailTaskPublisher
	tr        trace.Tracer
}

func NewOutboxWorkerUsecase(tx transactor, mailRep mailRepository, outboxRep outboxRepository, mailTP mailTaskPublisher, tr trace.Tracer) *OutboxWorkerUsecase {
	return &OutboxWorkerUsecase{
		tx:        tx,
		mailRep:   mailRep,
		outboxRep: outboxRep,
		mailTP:    mailTP,
		tr:        tr,
	}
}

func (uc *OutboxWorkerUsecase) Run(ctx context.Context) error {
	for {
		err := uc.proccessFirstAvailableTasks(ctx)
		if err != nil {
			if !errors.Is(err, ErrNoTasks) {
				slog.ErrorContext(ctx, "couldn't make transaction", "err", err)
			}
			time.Sleep(checkInterval)
		}
	}
}

func (uc *OutboxWorkerUsecase) proccessFirstAvailableTasks(ctx context.Context) error {
	err := uc.tx.WithTx(ctx, func(ctx context.Context) error {
		tasks, err := uc.outboxRep.GetFirstAvailableTasks(ctx, tasksBatchLimit)
		if err != nil {
			return err
		}

		var links []trace.Link
		for _, task := range tasks {
			carrier := propagation.MapCarrier{}
			_ = json.Unmarshal(task.TraceContext, &carrier)
			taskCtx := otel.GetTextMapPropagator().Extract(ctx, carrier)
			links = append(links, trace.LinkFromContext(taskCtx,
				attribute.Int("outbox.task_id", task.ID),
			))
		}
		ctx, span := uc.tr.Start(ctx, "outbox.process_tasks",
			trace.WithLinks(links...),
		)
		defer span.End()

		ids := make([]int, 0, len(tasks))
		for _, t := range tasks {
			ids = append(ids, t.AggregateID)
		}
		mails, err := uc.mailRep.GetMailsByIDs(ctx, ids)
		if err != nil {
			return err
		}

		for _, task := range tasks {
			carrier := propagation.MapCarrier{}
			_ = json.Unmarshal(task.TraceContext, &carrier)
			taskCtx := otel.GetTextMapPropagator().Extract(ctx, carrier)
			taskCtx = logger.WithContext(taskCtx, "trace_id", trace.SpanContextFromContext(taskCtx).TraceID())
			slog.InfoContext(taskCtx, "publishing new mail")

			mail := mails[task.AggregateID]
			if err := uc.mailTP.PublishMailTask(taskCtx, task, mail); err != nil {
				return err
			}
		}

		if err := uc.outboxRep.MarkTasksAsDone(ctx, tasks); err != nil {
			return err
		}
		return nil
	})

	return err
}
