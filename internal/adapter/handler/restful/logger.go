package restful

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/labasubagia/realworld-backend/internal/core/port"
)

func (s *Server) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// request id
		reqID := c.GetHeader("x-request-id")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		// make logger and sub-logger
		logger := s.logger.Field("request_id", reqID).Logger()
		c.Set(port.LoggerCtxKey, logger)

		// process request
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		// log
		logEvent := s.logger.Info()
		if c.Writer.Status() >= 500 {
			logEvent = s.logger.Error()
			if c.Request != nil && c.Request.Body != nil {
				if body, err := io.ReadAll(c.Request.Body); err == nil {
					logEvent.Field("body", body)
				}
			}
		}
		logEvent.
			Field("request_id", reqID).
			Field("protocol", "http").
			Field("method", c.Request.Method).
			Field("path", c.Request.URL.Path).
			Field("status_code", c.Writer.Status()).
			Field("status", http.StatusText(c.Writer.Status())).
			Field("duration", duration).
			Msg("received http request")
	}
}
