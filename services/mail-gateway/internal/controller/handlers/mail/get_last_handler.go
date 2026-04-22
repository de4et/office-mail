package mail

import (
	"github.com/de4et/office-mail/pkg/logger"
	"github.com/de4et/office-mail/services/mail-gateway/internal/controller/dto"
	"github.com/de4et/office-mail/services/mail-gateway/internal/domain"
	"github.com/de4et/office-mail/services/mail-gateway/internal/usecase"
	"github.com/gin-gonic/gin"
)

type GetLastHandler struct {
	uc *usecase.MailUsecase
}

func NewGetLastHandler(uc *usecase.MailUsecase) *GetLastHandler {
	return &GetLastHandler{
		uc: uc,
	}
}

type getLastRequest struct {
	From  string `form:"from"`
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
}

type getLastResponseSuccess struct {
	Mails []dto.Mail `json:"mails"`
}

func (h *GetLastHandler) Handle(c *gin.Context) {
	var req getLastRequest
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithError(400, err)
		return
	}

	ctx := logger.WithContext(c.Request.Context(), "limit", req.Limit)
	ctx = logger.WithContext(ctx, "page", req.Page)
	ctx = logger.WithContext(ctx, "from", req.From)

	mails, err := h.uc.GetLast(
		ctx,
		req.From,
		req.Limit,
		req.Page,
	)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.JSON(200, &getLastResponseSuccess{
		Mails: toDtoMails(mails),
	})
}

func toDtoMails(ms []domain.Mail) []dto.Mail {
	mails := make([]dto.Mail, len(ms))
	for i, m := range ms {
		mails[i] = dto.Mail{
			ID:   m.ID,
			To:   string(m.To),
			From: string(m.From),
			Body: m.Body,
		}
	}
	return mails
}
