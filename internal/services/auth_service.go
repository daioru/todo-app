package services

import (
	"errors"
	"os"
	"time"

	"github.com/daioru/todo-app/internal/logger"
	"github.com/daioru/todo-app/internal/models"
	"github.com/daioru/todo-app/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
	log  zerolog.Logger
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
		log:  logger.GetLogger(),
	}
}

func (s *AuthService) RegisterUser(user *models.User) error {
	exists, err := s.repo.UserExists(user)
	if err != nil {
		s.log.Error().
			Object("user", user).
			Err(err).
			Msg("Failed to check UserExists")
		return err
	}
	if exists {
		return repository.ErrUniqueUser
	}

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

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWTSECRET")))
	if err != nil {
		s.log.Error().Err(err).Msg("SignedString error")
		return "", err
	}

	return signedToken, nil
}
