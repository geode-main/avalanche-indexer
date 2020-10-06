package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func requestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)
		msg := "request"

		field := logger.
			WithField("method", c.Request.Method).
			WithField("client", c.ClientIP()).
			WithField("status", status).
			WithField("duration", duration.Milliseconds()).
			WithField("path", c.Request.URL.Path)

		if err := c.Errors.Last(); err != nil {
			msg = err.Error()
		}

		switch {
		case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
			field.Warn(msg)
		case status >= http.StatusInternalServerError:
			field.Error(msg)
		default:
			field.Info(msg)
		}
	}
}
