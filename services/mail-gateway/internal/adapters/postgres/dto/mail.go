package dto

import "github.com/de4et/office-mail/services/mail-gateway/internal/domain"

type Mail struct {
	ID   int    `db:"id"`
	To   string `db:"to_addr"`
	From string `db:"from_addr"`
	Body string `db:"body"`
}

func (m *Mail) ToDomain() domain.Mail {
	return domain.Mail{
		ID:   m.ID,
		To:   domain.Address(m.To),
		From: domain.Address(m.From),
		Body: m.Body,
	}
}
