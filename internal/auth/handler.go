package auth

import (
	"database/sql"
	"net/http"
	"time"


	"github.com/ericoliveiras/ecoponto-api/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Handler é a camada que lida com requisições de autenticação
type Handler struct {
	userRepo  *user.Repository 
	jwtSecret string
}

// NewHandler cria um novo handler de autenticação
func NewHandler(userRepo *user.Repository, jwtSecret string) *Handler {
	return &Handler{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// LoginRequest é a struct para o JSON de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login é o método para o endpoint POST /api/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	// 1. Valida o JSON de entrada
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Busca o usuário pelo email
	u, err := h.userRepo.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		// Se o usuário não existe (ErrNoRows) ou outro erro,
		// damos uma resposta genérica de "não autorizado"
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Compara a senha enviada com o hash salvo no banco
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))
	if err != nil {
		// Senha não bate (o erro é 'mismatched hash and password')
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"})
		return
	}

	// 4. Senha correta! Gerar o token JWT
	tokenString, err := h.generateToken(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	// 5. Retorna o token para o cliente
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

// generateToken é um helper privado para criar o JWT
func (h *Handler) generateToken(u *user.User) (string, error) {
	// 1. Define os "Claims" (dados) do token
	claims := jwt.MapClaims{
		"sub":   u.ID, // "Subject" (assunto) - o ID do usuário
		"email": u.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Expira em 24h
		"iat":   time.Now().Unix(),                     // "Issued At" (criado em)
	}

	// 2. Cria o token com o método de assinatura HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Assina o token com o nosso segredo
	return token.SignedString([]byte(h.jwtSecret))
}
