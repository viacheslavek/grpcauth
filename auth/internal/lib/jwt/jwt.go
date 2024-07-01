package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
)

// NewToken creates new JWT token for given owner and app
func NewToken(owner models.Owner, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = owner.Id
	claims["email"] = owner.Email
	claims["login"] = owner.Login
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.Id

	// TODO: откуда лучше брать секрет приложения? Лучше загружать в env и брать из env -> улучшить потом - техдолг
	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
