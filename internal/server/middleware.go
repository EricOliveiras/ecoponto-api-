package server

import (
	"net/http"

	"github.com/ericoliveiras/ecoponto-api/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware cria um middleware Gin para validar o JWT
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Valida o token
		token, err := auth.ValidateToken(c, jwtSecret)
		if err != nil {
			// Se o token for inválido, bloqueia a requisição
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Armazena os claims no contexto para uso posterior
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("claims", claims)
		}

		// Token é válido, continua para o próximo handler
		c.Next()
	}
}
