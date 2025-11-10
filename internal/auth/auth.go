package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ValidateToken extrai e valida o token JWT do header
func ValidateToken(c *gin.Context, jwtSecret string) (*jwt.Token, error) {
	// 1. Pegar o header "Authorization"
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("header Authorization não encontrado")
	}

	// 2. Checar se o formato é "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("formato do header Authorization inválido")
	}
	tokenString := parts[1]

	// 3. Parse e validação do token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Valida o método de assinatura (esperamos HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		// Retorna o segredo
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token inválido: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	return token, nil
}
