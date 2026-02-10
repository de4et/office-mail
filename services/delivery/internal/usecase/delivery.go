package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/de4et/office-mail/services/delivery/internal/domain"
)

var (
	ErrNoTasks = errors.New("no tasks")
)

type HandlerFunc func(ctx context.Context, mail domain.Mail, task domain.OutboxTask) error

type mailTaskConsumer interface {
	Consume(ctx context.Context, handler HandlerFunc)
}

type spamChecker interface {
	CheckMailForSpam(context.Context, domain.Mail) error
}

type idempotenceChecker interface {
	CheckTaskForIdempotence(context.Context, domain.OutboxTask) (bool, error)
}

type mailDelivery interface {
	DeliverMail(context.Context, domain.Mail) error
}

type OutboxWorkerUsecase struct {
	mailTC             mailTaskConsumer
	spamChecker        spamChecker
	idempotenceChecker idempotenceChecker
	mailDelivery       mailDelivery
}

func NewOutboxWorkerUsecase(mailTC mailTaskConsumer, spamChecker spamChecker, idempotenceChecker idempotenceChecker, mailDelivery mailDelivery) *OutboxWorkerUsecase {
	return &OutboxWorkerUsecase{
		mailTC:             mailTC,
		spamChecker:        spamChecker,
		idempotenceChecker: idempotenceChecker,
		mailDelivery:       mailDelivery,
	}
}

func (uc *OutboxWorkerUsecase) Run(ctx context.Context) error {
	uc.mailTC.Consume(ctx, func(ctx context.Context, mail domain.Mail, task domain.OutboxTask) error {

		slog.InfoContext(ctx, "retrieved new task", "task", task)

		ok, err := uc.idempotenceChecker.CheckTaskForIdempotence(ctx, task)
		if err != nil {
			slog.ErrorContext(ctx, "error during checking for idempotence", "err", err, "task.ID", task.ID)
			return err
		}
		if !ok {
			slog.InfoContext(ctx, "idempotence error", "task", task)
			return err
		}

		err = uc.spamChecker.CheckMailForSpam(ctx, mail)
		if err != nil {
			slog.ErrorContext(ctx, "It is spam!", "mail", mail)
			return err
		}

		err = uc.mailDelivery.DeliverMail(ctx, mail)
		if err != nil {
			slog.ErrorContext(ctx, "couldn't deliver mail", "mail", mail)
			return err
		}
		return nil
	})
	return nil
}
