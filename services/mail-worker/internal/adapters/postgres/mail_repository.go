package postgres

import (
	"context"
	_ "embed"

	"github.com/de4et/office-mail/pkg/postgres"
	"github.com/de4et/office-mail/services/mail-worker/internal/adapters/postgres/dto"
	"github.com/de4et/office-mail/services/mail-worker/internal/domain"
	"github.com/jmoiron/sqlx"
)

//go:embed queries/get_mail.sql
var getMailQuery string

type PostgresqlMailRepository struct {
	client *postgres.TxClient
}

func NewPostgresqlMailRepository(client *sqlx.DB) *PostgresqlMailRepository {
	return &PostgresqlMailRepository{
		client: postgres.NewTxClient(client),
	}
}

func (rep *PostgresqlMailRepository) GetMail(ctx context.Context, mailID int) (domain.Mail, error) {
	var resp dto.GetMailResponse
	err := rep.client.GetContext(
		ctx,
		&resp,
		getMailQuery,
		mailID,
	)
	if err != nil {
		return domain.Mail{}, err
	}

	return domain.Mail{
		ID:   resp.ID,
		To:   domain.Address(resp.To),
		From: domain.Address(resp.From),
		Body: resp.Body,
	}, nil
}
