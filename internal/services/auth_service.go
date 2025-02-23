package services

import (
	"errors"

	"github.com/daioru/todo-app/internal/models"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/daioru/todo-app/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       *repository.UserRepository
	jwtService *utils.JWTService
}

func NewAuthService(repo *repository.UserRepository, jwtService *utils.JWTService) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *AuthService) RegisterUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	return s.repo.CreateUser(user)
}

func (s *AuthService) LoginUser(username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credential")
	}

	return s.jwtService.GenerateToken(user.ID)
}