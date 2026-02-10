package middleware

import (
	"github.com/de4et/office-mail/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func LogTraceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()

		lctx := logger.WithContext(ctx, "trace_id", traceID)
		c.Request = c.Request.WithContext(lctx)
		c.Next()
	}
}
