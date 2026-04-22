package usecase

import (
	"context"
	"log/slog"

	"github.com/de4et/office-mail/services/mail-gateway/internal/domain"
)

type mailRepository interface {
	CreateMail(context.Context, domain.Mail) (int, error)
	GetWithLimitAndOffset(context.Context, string, int, int) ([]domain.Mail, error)
}

type outboxRepository interface {
	CreateOutboxDeliveryTask(context.Context, domain.Mail) error
}

type MailUsecase struct {
	tx        transactor
	mailRep   mailRepository
	outboxRep outboxRepository
}

func NewMailUsecase(tx transactor, mailRep mailRepository, outboxRep outboxRepository) *MailUsecase {
	return &MailUsecase{
		tx:        tx,
		mailRep:   mailRep,
		outboxRep: outboxRep,
	}
}

func (uc *MailUsecase) Send(ctx context.Context, mail domain.Mail) error {
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

func (uc *MailUsecase) GetLast(ctx context.Context, from string, limit, page int) ([]domain.Mail, error) {
	slog.InfoContext(ctx, "Getting last")
	return uc.mailRep.GetWithLimitAndOffset(ctx, from, limit, page*limit)
}
