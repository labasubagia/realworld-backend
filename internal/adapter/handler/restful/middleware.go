package restful

import (
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationTypeToken = "token"
	authorizationArgKey    = "authorization_arg"
)

func (server *Server) AuthMiddleware(autoDenied bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authArg, err := server.parseToken(c)
		if err != nil {
			if autoDenied {
				errorHandler(c, err)
				return
			}
		}
		c.Set(authorizationArgKey, authArg)
		c.Next()
	}
}
