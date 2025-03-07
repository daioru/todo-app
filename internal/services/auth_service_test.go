package services_test

import (
	"errors"
	"os"
	"testing"

	"github.com/daioru/todo-app/internal/models"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/daioru/todo-app/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) UserExists(user *models.User) (bool, error) {
	args := m.Called(user)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByID(userID int) (*models.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestRegisterUser(t *testing.T) {
	user := &models.User{Username: "test username", Password: "test password"}

	t.Run("User already exists", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := services.NewAuthService(mockRepo)

		mockRepo.On("UserExists", user).Return(true, nil)

		err := service.RegisterUser(user)
		assert.ErrorIs(t, err, repository.ErrUniqueUser)

		mockRepo.AssertNotCalled(t, "CreateUser")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Successful registration", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := services.NewAuthService(mockRepo)

		mockRepo.On("UserExists", user).Return(false, nil)
		mockRepo.On("CreateUser", user).Return(nil)

		err := service.RegisterUser(user)
		assert.NoError(t, err)

		mockRepo.AssertCalled(t, "CreateUser", user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error checking UserExists", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := services.NewAuthService(mockRepo)

		mockRepo.On("UserExists", user).Return(false, errors.New("some error"))

		err := service.RegisterUser(user)
		assert.Error(t, err)

		mockRepo.AssertNotCalled(t, "CreateUser")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error creating user", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := services.NewAuthService(mockRepo)

		mockRepo.On("UserExists", user).Return(false, nil)
		mockRepo.On("CreateUser", user).Return(errors.New("failed to create user"))

		err := service.RegisterUser(user)
		assert.Error(t, err)

		mockRepo.AssertCalled(t, "CreateUser", user)
		mockRepo.AssertExpectations(t)
	})
}

func TestLoginUser(t *testing.T) {
	t.Run("User not found", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := services.NewAuthService(mockRepo)

		mockRepo.On("GetUserByUsername", "nonexistent").Return((*models.User)(nil), errors.New("user not found"))

		token, err := service.LoginUser("nonexistent", "password")
		assert.Empty(t, token)
		assert.ErrorIs(t, err, repository.ErrUserNotFound)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid password", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := services.NewAuthService(mockRepo)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
		user := &models.User{ID: 1, Username: "testuser", PasswordHash: string(hashedPassword)}

		mockRepo.On("GetUserByUsername", "testuser").Return(user, nil)

		token, err := service.LoginUser("testuser", "wrong_password")
		assert.Empty(t, token)
		assert.ErrorIs(t, err, services.ErrInvalidCredentials)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Successful login", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := services.NewAuthService(mockRepo)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
		user := &models.User{ID: 1, Username: "testuser", PasswordHash: string(hashedPassword)}

		mockRepo.On("GetUserByUsername", "testuser").Return(user, nil)

		os.Setenv("JWTSECRET", "testsecret") // Устанавливаем секретный ключ

		token, err := service.LoginUser("testuser", "correct_password")
		assert.NotEmpty(t, token)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})
}
