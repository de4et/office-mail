package postgres

import (
	"context"
	_ "embed"

	"github.com/de4et/office-mail/pkg/postgres"
	"github.com/de4et/office-mail/services/mail-gateway/internal/adapters/postgres/dto"
	"github.com/de4et/office-mail/services/mail-gateway/internal/domain"
	"github.com/jmoiron/sqlx"
)

//go:embed queries/create_mail.sql
var createMailQuery string

//go:embed queries/get_last.sql
var getLastQuery string

type PostgresqlMailRepository struct {
	client *postgres.TxClient
}

func NewPostgresqlMailRepository(client *sqlx.DB) *PostgresqlMailRepository {
	return &PostgresqlMailRepository{
		client: postgres.NewTxClient(client),
	}
}

func (rep *PostgresqlMailRepository) CreateMail(ctx context.Context, mail domain.Mail) (int, error) {
	var id int
	err := rep.client.GetContext(
		ctx,
		&id,
		createMailQuery,
		mail.Body,
		mail.From,
		mail.To,
	)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (rep *PostgresqlMailRepository) GetWithLimitAndOffset(ctx context.Context, from string, limit int, offset int) ([]domain.Mail, error) {
	var res []dto.Mail
	err := rep.client.SelectContext(
		ctx,
		&res,
		getLastQuery,
		from,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	mails := make([]domain.Mail, len(res))
	for i := range res {
		mails[i] = res[i].ToDomain()
	}
	return mails, nil
}
