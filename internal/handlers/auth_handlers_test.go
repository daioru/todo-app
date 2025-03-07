package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daioru/todo-app/internal/handlers"
	"github.com/daioru/todo-app/internal/models"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/daioru/todo-app/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) RegisterUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthService) LoginUser(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func TestRegister(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("Successful registration", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		user := models.User{Username: "testuser", Password: "password123"}
		mockService.On("RegisterUser", &user).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"username": "testuser", "password": "password123"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{"username":}`))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request")
	})

	t.Run("Username already taken", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		user := models.User{Username: "testuser", Password: "password123"}
		mockService.On("RegisterUser", &user).Return(repository.ErrUniqueUser)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"username": "testuser", "password": "password123"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "username already taken")
	})

	t.Run("Internal server error", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		user := models.User{Username: "testuser", Password: "password123"}
		mockService.On("RegisterUser", &user).Return(errors.New("DB error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"username": "testuser", "password": "password123"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "server side error")
	})
}

func TestLogin(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("Successful login", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		mockService.On("LoginUser", "testuser", "password123").Return("valid-token", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"username": "testuser", "password": "password123"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		cookie := w.Header().Get("Set-Cookie")
		assert.Contains(t, cookie, "Authorization=valid-token")
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"username":}`))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request")
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		mockService.On("LoginUser", "testuser", "wrongpassword").Return("", services.ErrInvalidCredentials)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"username": "testuser", "password": "wrongpassword"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "ivalid credentials")
	})

	t.Run("Internal server error", func(t *testing.T) {
		t.Parallel()
		mockService := new(MockAuthService)
		handler := handlers.NewAuthHandler(mockService)

		mockService.On("LoginUser", "testuser", "password123").Return("", errors.New("DB error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := `{"username": "testuser", "password": "password123"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "server side error")
	})
}
