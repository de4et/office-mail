package mail

import (
	"encoding/json"

	"github.com/de4et/office-mail/services/mail-gateway/internal/domain"
	"github.com/de4et/office-mail/services/mail-gateway/internal/usecase"
	"github.com/gin-gonic/gin"
)

type SendHandler struct {
	// TODO: logger
	sendUC *usecase.SendMailUsecase
}

func NewSendHandler(sendUC *usecase.SendMailUsecase) *SendHandler {
	return &SendHandler{
		sendUC: sendUC,
	}
}

type sendMailRequest struct {
	To   string `json:"to"`
	From string `json:"from"`
	Body string `json:"body"`
}

type sendMailResponseSuccess struct {
	Message string `json:"messages"`
}

func (h *SendHandler) Handle(c *gin.Context) {
	var req sendMailRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if err := h.sendUC.Send(
		c,
		domain.Mail{
			To:   domain.Address(req.To),
			From: domain.Address(req.From),
			Body: req.Body,
		},
	); err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.JSON(200, &sendMailResponseSuccess{
		Message: "successfully sent message",
	})
}
