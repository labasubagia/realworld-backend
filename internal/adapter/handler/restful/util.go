package restful

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getPagination(c *gin.Context) (offset, limit int) {
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	limit, err = strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 20
	}
	return offset, limit
}
