package restful

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// process request
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		// log
		logger := s.logger.Info()
		if c.Writer.Status() >= 500 {
			logger = s.logger.Error()
			if c.Request != nil && c.Request.Body != nil {
				if body, err := io.ReadAll(c.Request.Body); err == nil {
					logger.Field("body", body)
				}
			}
		}
		logger.
			Field("protocol", "http").
			Field("method", c.Request.Method).
			Field("path", c.Request.URL.Path).
			Field("status_code", c.Writer.Status()).
			Field("status", http.StatusText(c.Writer.Status())).
			Field("duration", duration).
			Msg("received http request")
	}
}
