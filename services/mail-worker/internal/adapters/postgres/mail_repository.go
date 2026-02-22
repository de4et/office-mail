package postgres

import (
	"context"
	_ "embed"

	"github.com/de4et/office-mail/pkg/postgres"
	"github.com/de4et/office-mail/services/mail-worker/internal/adapters/postgres/dto"
	"github.com/de4et/office-mail/services/mail-worker/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
)

//go:embed queries/get_mail.sql
var getMailQuery string

//go:embed queries/get_mails.sql
var getMailsByIDsQuery string

type PostgresqlMailRepository struct {
	client *postgres.TxClient
	tr     trace.Tracer
}

func NewPostgresqlMailRepository(client *sqlx.DB, tr trace.Tracer) *PostgresqlMailRepository {
	return &PostgresqlMailRepository{
		client: postgres.NewTxClient(client),
		tr:     tr,
	}
}

func (rep *PostgresqlMailRepository) GetMail(ctx context.Context, mailID int) (domain.Mail, error) {
	ctx, span := rep.tr.Start(ctx, "mailRepository.GetMail")
	defer span.End()

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

func (rep *PostgresqlMailRepository) GetMailsByIDs(ctx context.Context, ids []int) (map[int]domain.Mail, error) {
	ctx, span := rep.tr.Start(ctx, "mailRepository.GetMailsByIDs")
	defer span.End()

	if len(ids) == 0 {
		return map[int]domain.Mail{}, nil
	}

	rows, err := rep.client.QueryContext(ctx, getMailsByIDsQuery, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]domain.Mail, len(ids))

	for rows.Next() {
		var m dto.GetMailResponse
		if err := rows.Scan(&m.ID, &m.From, &m.To, &m.Body); err != nil {
			return nil, err
		}
		result[m.ID] = domain.Mail{
			ID:   m.ID,
			To:   domain.Address(m.To),
			From: domain.Address(m.From),
			Body: m.Body,
		}
	}

	return result, rows.Err()
}
