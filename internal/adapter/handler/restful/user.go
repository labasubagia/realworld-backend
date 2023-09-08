package restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
)

type UserResult struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

type UserRegisterParams struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	User UserRegisterParams `json:"user"`
}

type RegisterResponse struct {
	User UserResult `json:"user"`
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
		User: UserResult{
			Email:    result.User.Email,
			Username: result.User.Username,
			Bio:      result.User.Bio,
			Image:    result.User.Email,
			Token:    result.Token,
		},
	}
	c.JSON(http.StatusCreated, res)
}

type UserLoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	User UserLoginParams `json:"user"`
}

type LoginResponse struct {
	User UserResult `json:"user"`
}

func (server *Server) Login(c *gin.Context) {
	req := LoginRequest{}
	if err := c.BindJSON(&req); err != nil {
		errorHandler(c, err)
		return
	}

	result, err := server.service.User().Login(c, port.LoginUserParams{
		User: domain.User{
			Email:    req.User.Email,
			Password: req.User.Password,
		},
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := LoginResponse{
		User: UserResult{
			Email:    result.User.Email,
			Username: result.User.Username,
			Bio:      result.User.Bio,
			Image:    result.User.Email,
			Token:    result.Token,
		},
	}
	c.JSON(http.StatusCreated, res)
}
