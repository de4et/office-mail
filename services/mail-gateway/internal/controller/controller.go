package controller

import (
	"context"
	"net/http"
	"time"

	mailhandlers "github.com/de4et/office-mail/services/mail-gateway/internal/controller/handlers/mail"
	"github.com/de4et/office-mail/services/mail-gateway/internal/controller/middleware"
	"github.com/de4et/office-mail/services/mail-gateway/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	g      *gin.Engine
	ctx    context.Context
	server *http.Server
}

func SetupRoutes(ctx context.Context, sendUC *usecase.SendMailUsecase) *Controller {
	controller := &Controller{
		g:   gin.Default(),
		ctx: ctx,
	}

	send_handler := mailhandlers.NewSendHandler(sendUC)

	// TODO: Middlewares(logger, metrics, tracing, auth) (propagate user in context)
	mailGroup := controller.g.Group("/mail/")
	mailGroup.Use(middleware.ErrorHandler())
	mailGroup.POST("/send", send_handler.Handle)

	return controller
}

func (controller *Controller) Listen(addr string) {
	controller.server = &http.Server{
		Addr:         addr,
		Handler:      controller.g,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	controller.server.ListenAndServe()
}

func (contoller *Controller) Shutdown() error {
	return contoller.server.Shutdown(contoller.ctx)
}

// locahost:8080/mail/send body: {"from": "", to: "", "msg": ""}, cookie: {"jwt": ""}
