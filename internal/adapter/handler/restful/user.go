package restful

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
)

type RegisterRequestUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	User RegisterRequestUser `json:"user"`
}

func (server *Server) Register(c *gin.Context) {
	req := RegisterRequest{}
	if err := c.BindJSON(&req); err != nil {
		errorHandler(c, err)
		return
	}

	user, err := server.service.User().Register(c, port.RegisterParams{
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

	res := UserResponse{serializeUser(user)}
	c.JSON(http.StatusCreated, res)
}

type LoginParamUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	User LoginParamUser `json:"user"`
}

func (server *Server) Login(c *gin.Context) {
	req := LoginRequest{}
	if err := c.BindJSON(&req); err != nil {
		errorHandler(c, err)
		return
	}

	user, err := server.service.User().Login(c, port.LoginParams{
		User: domain.User{
			Email:    req.User.Email,
			Password: req.User.Password,
		},
	})
	if err != nil {
		errorHandler(c, err)
		return
	}

	res := UserResponse{serializeUser(user)}
	c.JSON(http.StatusOK, res)
}

func (server *Server) CurrentUser(c *gin.Context) {
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}
	user, err := server.service.User().Current(context.Background(), authArg)
	if err != nil {
		errorHandler(c, err)
		return
	}
	res := UserResponse{serializeUser(user)}
	c.JSON(http.StatusOK, res)
}

type UpdateUser struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Image    string `json:"image,omitempty"`
}

type UpdateUserRequest struct {
	User UpdateUser `json:"user"`
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
	user, err := server.service.User().Update(context.Background(), port.UpdateUserParams{
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
	res := UserResponse{serializeUser(user)}
	c.JSON(http.StatusOK, res)
}

func (server *Server) Profile(c *gin.Context) {
	username := c.Param("username")
	authArg, _ := getAuthArg(c)
	user, err := server.service.User().Profile(context.Background(), port.ProfileParams{
		Username: username,
		AuthArg:  authArg,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}
	res := ProfileResponse{serializeProfile(user)}
	c.JSON(http.StatusOK, res)
}

func (server *Server) FollowUser(c *gin.Context) {
	username := c.Param("username")
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}
	user, err := server.service.User().Follow(context.Background(), port.ProfileParams{
		Username: username,
		AuthArg:  authArg,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}
	res := ProfileResponse{serializeProfile(user)}
	c.JSON(http.StatusOK, res)
}

func (server *Server) UnFollowUser(c *gin.Context) {
	username := c.Param("username")
	authArg, err := getAuthArg(c)
	if err != nil {
		errorHandler(c, err)
		return
	}
	user, err := server.service.User().UnFollow(context.Background(), port.ProfileParams{
		Username: username,
		AuthArg:  authArg,
	})
	if err != nil {
		errorHandler(c, err)
		return
	}
	res := ProfileResponse{serializeProfile(user)}
	c.JSON(http.StatusOK, res)
}
