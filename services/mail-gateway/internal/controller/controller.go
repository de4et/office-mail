package controller

import (
	"context"
	"net/http"
	"time"

	mailhandlers "github.com/de4et/office-mail/services/mail-gateway/internal/controller/handlers/mail"
	"github.com/de4et/office-mail/services/mail-gateway/internal/controller/middleware"
	"github.com/de4et/office-mail/services/mail-gateway/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Controller struct {
	g      *gin.Engine
	ctx    context.Context
	server *http.Server
}

func SetupRoutes(ctx context.Context, sendUC *usecase.SendMailUsecase) *Controller {
	controller := &Controller{
		g:   gin.New(),
		ctx: ctx,
	}

	controller.g.Use(gin.Recovery())
	controller.g.Use(otelgin.Middleware("mail-gateway"))
	controller.g.Use(middleware.LogHandler())
	controller.g.Use(middleware.LogTraceHandler())
	controller.g.Use(middleware.ErrorHandler())
	controller.g.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	send_handler := mailhandlers.NewSendHandler(sendUC)

	// TODO: Middlewares(logger, metrics, tracing, auth) (propagate user in context)
	mailGroup := controller.g.Group("/mail/")
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
