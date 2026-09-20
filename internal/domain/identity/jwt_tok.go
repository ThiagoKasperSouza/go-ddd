package identity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, roles []Role, secret string) (string, error) {
	var rolesStr []string
	for _, r := range roles {
		rolesStr = append(rolesStr, string(r))
	}

	claims := jwt.MapClaims{
		"sub":   userID,
		"roles": rolesStr,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Expira em 24 horas
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}