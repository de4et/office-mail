package postgres

import (
	"context"
	_ "embed"

	"github.com/de4et/office-mail/pkg/postgres"
	"github.com/de4et/office-mail/services/mail-gateway/internal/domain"
	"github.com/jmoiron/sqlx"
)

//go:embed queries/create_mail.sql
var createMailQuery string

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
