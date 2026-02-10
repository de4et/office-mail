package stub

import (
	"context"

	"github.com/de4et/office-mail/services/delivery/internal/domain"
)

type StubMailDelivery struct {
}

func NewStubMailDelivery() *StubMailDelivery {
	return &StubMailDelivery{}
}

func (md *StubMailDelivery) DeliverMail(ctx context.Context, mail domain.Mail) error {
	return nil
}
