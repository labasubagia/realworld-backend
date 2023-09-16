package restful

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
)

func (s *Server) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// request id
		reqID := c.GetHeader("x-request-id")
		if reqID == "" {
			reqID = domain.NewID().String()
		}

		// make logger and sub-logger
		// ? make unique instance for each handler/interactor request
		logger := s.logger.NewInstance().Field("request_id", reqID).Logger()
		// logger := s.logger.Field("request_id", reqID).Logger() // this is use single instance
		c.Set(port.SubLoggerCtxKey, logger)

		// process request
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		// log
		logEvent := logger.Info()
		if c.Writer.Status() >= 500 {
			logEvent = logger.Error()
			if c.Request != nil && c.Request.Body != nil {
				if body, err := io.ReadAll(c.Request.Body); err == nil {
					logEvent.Field("body", body)
				}
			}
		}
		logEvent.
			Field("protocol", "http").
			Field("client_ip", c.ClientIP()).
			Field("user_agent", c.Request.UserAgent()).
			Field("method", c.Request.Method).
			Field("path", c.Request.URL.Path).
			Field("status_code", c.Writer.Status()).
			Field("status", http.StatusText(c.Writer.Status())).
			Field("duration", duration).
			Msg("received http request")
	}
}
