package middleware

import (
	sl "dictionary_app/utils/logger"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		duration := time.Since(start)
		durationMs := duration.Milliseconds()
		sl.GetLogger().Info("Request",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.URL.Path),
			slog.Int("status", ctx.Writer.Status()),
			slog.Int64("duration", durationMs))
	}
}
