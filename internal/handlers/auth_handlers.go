package handlers

import (
	"errors"
	"net/http"

	"github.com/daioru/todo-app/internal/models"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/daioru/todo-app/internal/services"
	"github.com/gin-gonic/gin"

	_ "github.com/daioru/todo-app/docs"
)

type IAuthService interface {
	RegisterUser(user *models.User) error
	LoginUser(username, password string) (string, error)
}

type AuthHandler struct {
	service IAuthService
}

func NewAuthHandler(service IAuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// @Summary Register
// @Description create account
// @Accept  json
// @Produce  json
// @Tags auth
// @Param input body UserData true "user info"
// @Success 201
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.service.RegisterUser(&req)
	if err != nil {
		if err == repository.ErrUniqueUser {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already taken"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "server side error"})
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}

// @Summary Login
// @Description user login to set auth cookie
// @Accept  json
// @Produce  json
// @Tags auth
// @Param input body UserData true "user info"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.service.LoginUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) || errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ivalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server side error"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600*24*3, "", "", false, true)

	c.AbortWithStatus(http.StatusOK)
}
