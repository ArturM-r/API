package jwt

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)
import "github.com/golang-jwt/jwt/v5"

func Signing(UserID uuid.UUID, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": UserID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("process tokenizing is going wrong: %w", err)
	}
	return tokenString, nil
}
