package ginmiddleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GinLogFunc - Logger forces gin to use our logger
// Adapted from gin.Logger
func GinLogFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		logger := log.Debug()

		switch {
		case statusCode < 300:
			logger = log.Debug()
		case statusCode < 400:
			logger = log.Info()
		default:
			logger = log.Error()
		}

		logger.Fields(map[string]interface{}{
			"clientIP": clientIP,
			"path":     path,
			"method":   method,
			"status":   statusCode,
		}).Msgf("[GIN] %13v | %v",
			latency,
			comment,
		)
	}
}
