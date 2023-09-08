package restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
)

func errorHandler(c *gin.Context, err error) {
	fail, ok := err.(*exception.Exception)
	if !ok {
		c.JSON(http.StatusInternalServerError, fail)
		return
	}
	var statusCode int
	switch fail.Type {
	case exception.TypeNotFound:
		statusCode = http.StatusNotFound
	case exception.TypeTokenExpired, exception.TypeTokenInvalid:
		statusCode = http.StatusUnauthorized
	case exception.TypeValidation:
		if !fail.HasError() {
			fail.AddError("exception", fail.Cause.Error())
		}
		statusCode = http.StatusUnprocessableEntity
	default:
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, fail)
}
