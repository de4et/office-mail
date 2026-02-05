package usecase

import (
	"context"
	"log/slog"

	"github.com/de4et/office-mail/services/mail-gateway/internal/domain"
)

type mailRepository interface {
	CreateMail(context.Context, domain.Mail) (int, error)
}

type outboxRepository interface {
	CreateOutboxDeliveryTask(context.Context, domain.Mail) error
}

type SendMailUsecase struct {
	tx        transactor
	mailRep   mailRepository
	outboxRep outboxRepository
}

func NewSendMailUsecase(tx transactor, mailRep mailRepository, outboxRep outboxRepository) *SendMailUsecase {
	return &SendMailUsecase{
		tx:        tx,
		mailRep:   mailRep,
		outboxRep: outboxRep,
	}
}

func (uc *SendMailUsecase) Send(ctx context.Context, mail domain.Mail) error {
	slog.InfoContext(ctx, "Sending", "mail", mail)
	return uc.tx.WithTx(ctx, func(ctx context.Context) error {
		id, err := uc.mailRep.CreateMail(ctx, mail)
		if err != nil {
			return err
		}
		mail.ID = id

		err = uc.outboxRep.CreateOutboxDeliveryTask(ctx, mail)
		if err != nil {
			return err
		}

		return nil
	})
}
