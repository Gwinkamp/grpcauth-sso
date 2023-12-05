package jwt

import (
	"time"

	"github.com/Gwinkamp/grpcauth-sso/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

// NewToken генерирует новый доступа к конкретному сервису токен для пользователя
func NewToken(user models.User, service models.Service, duraton time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duraton).Unix()
	claims["service_id"] = service.ID

	tokenString, err := token.SignedString([]byte(service.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
