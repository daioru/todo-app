package utils

import (
	"errors"
	"time"

	"github.com/daioru/todo-app/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
)

type JWTService struct {
	jwtSecret []byte
	log       zerolog.Logger
}

func NewJwtService(jwtSecret []byte) *JWTService {
	return &JWTService{
		jwtSecret: jwtSecret,
		log:       logger.GetLogger(),
	}
}

func (s *JWTService) GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		s.log.Error().Err(err).Msg("SignedString error")
		return "", err
	}

	return signedToken, nil
}

func (s *JWTService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	//Expiration check
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return 0, errors.New("token expired")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id in token")
	}

	return int(userID), nil
}
