package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go-ddd/internal/domain/identity"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserClaims struct {
	UserID string          `json:"user_id"`
	Roles  []identity.Role `json:"roles"`
}

func (c *UserClaims) HasRole(role identity.Role) bool {
	if c == nil {
		return false
	}
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// EnsureAuthenticated valida o token JWT enviado no Header Authorization
func EnsureAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "authorization header missing")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "invalid authorization format")
			return
		}

		tokenString := parts[1]

		// Valida o Token JWT
		claims, err := validateJWT(tokenString)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("invalid token: %v", err))
			return
		}

		// Adiciona o UserClaims no Contexto da Requisição
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateJWT lê e verifica a assinatura e expiração do JWT
func validateJWT(tokenString string) (*UserClaims, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		// Fallback para desenvolvimento caso a variável de ambiente não esteja configurada
		jwtSecret = []byte("dev")
	}

	// Faz o parse e verifica a assinatura do token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Garante que o método de assinatura é HMAC (HS256, HS384, HS512)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token is invalid or expired")
	}

	// Extrai os claims (payload) do token
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot parse token claims")
	}

	userID, _ := mapClaims["sub"].(string)
	if userID == "" {
		// Caso a chave utilizada seja 'user_id' em vez do padrão 'sub'
		userID, _ = mapClaims["user_id"].(string)
	}

	// Processa o array de roles vindos do payload JWT
	var userRoles []identity.Role
	if rawRoles, ok := mapClaims["roles"].([]interface{}); ok {
		for _, r := range rawRoles {
			if roleStr, ok := r.(string); ok {
				role := identity.Role(roleStr)
				if role.IsValid() {
					userRoles = append(userRoles, role)
				}
			}
		}
	}

	return &UserClaims{
		UserID: userID,
		Roles:  userRoles,
	}, nil
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// RequireRole exige que o usuário autenticado possua uma role específica
func RequireRole(requiredRole identity.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserContextKey).(*UserClaims)
			if !ok || !user.HasRole(requiredRole) {
				respondWithError(w, http.StatusForbidden, "insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}