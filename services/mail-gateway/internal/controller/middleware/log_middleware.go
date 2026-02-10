package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func LogHandler() gin.HandlerFunc {
	slog.Error("here11")
	return sloggin.NewWithConfig(slog.Default(), sloggin.Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,
		WithRequestID:    true,
		WithSpanID:       true,
		WithTraceID:      true,
		WithClientIP:     true,
		HandleGinDebug:   true,
	})
}
