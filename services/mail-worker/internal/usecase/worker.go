package usecase

import (
	"context"

	"github.com/de4et/office-mail/services/mail-worker/internal/domain"
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
	return uc.tx.WithTx(ctx, func(ctx context.Context) error {
		task, err := uc.outboxRep.GetFirstAvailableTask(ctx)
		if err != nil {
			return err
		}

		mail, err := uc.mailRep.GetMail(ctx, task.AggregateID)
		if err != nil {
			return err
		}

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
}
