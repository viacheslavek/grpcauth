package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
)

func NewToken(owner models.Owner, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = owner.Id()
	claims["email"] = owner.Email()
	claims["login"] = owner.Login()
	claims["exp"] = time.Now().Add(duration).Unix()

	secret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
