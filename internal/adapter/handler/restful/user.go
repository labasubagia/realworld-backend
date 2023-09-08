package restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type RegisterRequest struct {
	User User `json:"user"`
}

type RegisterResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

func (server *Server) Register(c *gin.Context) {
	req := RegisterRequest{}
	if err := c.BindJSON(&req); err != nil {
		errorHandler(c, err)
		return
	}

	result, err := server.service.User().Register(c, port.RegisterUserParams{
		User: domain.User{
			Email:    req.User.Email,
			Username: req.User.Username,
			Password: req.User.Password,
		},
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := RegisterResponse{
		Email:    result.User.Email,
		Username: result.User.Username,
		Bio:      result.User.Bio,
		Image:    result.User.Image,
		Token:    result.Token,
	}
	c.JSON(http.StatusCreated, res)
}
