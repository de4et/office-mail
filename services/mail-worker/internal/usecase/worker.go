package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/de4et/office-mail/pkg/logger"
	"github.com/de4et/office-mail/services/mail-worker/internal/domain"
)

var (
	ErrNoTasks = errors.New("no tasks")

	tr = otel.Tracer("mail-worker")
)

type mailRepository interface {
	GetMail(context.Context, int) (domain.Mail, error)
}

type outboxRepository interface {
	GetFirstAvailableTask(context.Context) (domain.OutboxTask, error)
	MarkTaskAsDone(context.Context, domain.OutboxTask) error
}

type mailTaskPublisher interface {
	PublishMailTask(context.Context, domain.OutboxTask, domain.Mail) error
}

type OutboxWorkerUsecase struct {
	tx        transactor
	mailRep   mailRepository
	outboxRep outboxRepository
	mailTP    mailTaskPublisher
}

func NewOutboxWorkerUsecase(tx transactor, mailRep mailRepository, outboxRep outboxRepository, mailTP mailTaskPublisher) *OutboxWorkerUsecase {
	return &OutboxWorkerUsecase{
		tx:        tx,
		mailRep:   mailRep,
		outboxRep: outboxRep,
		mailTP:    mailTP,
	}
}

func (uc *OutboxWorkerUsecase) Run(ctx context.Context) error {
	for {
		err := uc.proccessFirstAvailableTask(ctx)
		if err != nil && !errors.Is(err, ErrNoTasks) {
			slog.ErrorContext(ctx, "couldn't make transaction", "err", err)
		}
	}
}

func (uc *OutboxWorkerUsecase) proccessFirstAvailableTask(ctx context.Context) error {
	err := uc.tx.WithTx(ctx, func(ctx context.Context) error {
		task, err := uc.outboxRep.GetFirstAvailableTask(ctx)
		if err != nil {
			return err
		}

		carrier := propagation.MapCarrier{}
		_ = json.Unmarshal(task.TraceContext, &carrier)

		ctx = otel.GetTextMapPropagator().
			Extract(ctx, carrier)

		// span.SpanContext().TraceID

		ctx, span := tr.Start(ctx, "outbox.process_first_available")
		defer span.End()

		mail, err := uc.mailRep.GetMail(ctx, task.AggregateID)
		if err != nil {
			return err
		}

		ctx = logger.WithContext(ctx, "trace_id", span.SpanContext().TraceID())

		slog.InfoContext(ctx, "publishing new mail")

		err = uc.mailTP.PublishMailTask(ctx, task, mail)
		if err != nil {
			return err
		}

		err = uc.outboxRep.MarkTaskAsDone(ctx, task)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
