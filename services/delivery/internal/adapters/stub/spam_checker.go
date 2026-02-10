package stub

import (
	"context"

	"github.com/de4et/office-mail/services/delivery/internal/domain"
)

type StubSpamChecker struct {
}

func NewStubSpamChecker() *StubSpamChecker {
	return &StubSpamChecker{}
}

func (sc *StubSpamChecker) CheckMailForSpam(ctx context.Context, mail domain.Mail) error {
	return nil
}
