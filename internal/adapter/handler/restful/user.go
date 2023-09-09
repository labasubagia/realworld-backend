package restful

import (
	"context"
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
	c.JSON(http.StatusOK, res)
}

type CurrentUserResponse struct {
	User UserResult `json:"user"`
}

func (server *Server) CurrentUser(c *gin.Context) {
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}
	result, err := server.service.User().Current(context.Background(), authArg)
	if err != nil {
		errorHandler(c, err)
		return
	}
	res := CurrentUserResponse{
		User: UserResult{
			Email:    result.User.Email,
			Username: result.User.Username,
			Bio:      result.User.Bio,
			Image:    result.User.Image,
			Token:    result.Token,
		},
	}
	c.JSON(http.StatusOK, res)
}

type UserUpdateParams struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Image    string `json:"image,omitempty"`
}

type UpdateUserRequest struct {
	User UserUpdateParams `json:"user"`
}

type UpdateUserResult struct {
	User UserResult `json:"user"`
}

func (server *Server) UpdateUser(c *gin.Context) {
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}
	req := UpdateUserRequest{}
	if err := c.BindJSON(&req); err != nil {
		errorHandler(c, err)
		return
	}
	result, err := server.service.User().Update(context.Background(), port.UpdateUserParams{
		AuthArg: authArg,
		User: domain.User{
			ID:       authArg.Payload.UserID,
			Email:    req.User.Email,
			Username: req.User.Username,
			Password: req.User.Password,
			Image:    req.User.Image,
			Bio:      req.User.Bio,
		},
	})
	if err != nil {
		errorHandler(c, err)
		return
	}
	res := UpdateUserResult{
		User: UserResult{
			Email:    result.User.Email,
			Username: result.User.Username,
			Bio:      result.User.Bio,
			Image:    result.User.Image,
			Token:    result.Token,
		},
	}
	c.JSON(http.StatusOK, res)
}

type UserProfile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type ProfileUserResult struct {
	Profile UserProfile `json:"profile"`
}

func (server *Server) Profile(c *gin.Context) {
	username := c.Param("username")
	authArg, _ := getAuthArg(c)
	result, err := server.service.User().Profile(context.Background(), port.ProfileParams{
		Username: username,
		AuthArg:  authArg,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}
	res := ProfileUserResult{
		Profile: UserProfile{
			Username:  result.User.Username,
			Bio:       result.User.Bio,
			Image:     result.User.Image,
			Following: result.IsFollow,
		},
	}
	c.JSON(http.StatusOK, res)
}
