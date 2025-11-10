package server

import (
	"fmt"
	"net/http"

	"github.com/ericoliveiras/ecoponto-api/internal/ecoponto"
	"github.com/gin-gonic/gin"
)

// Server é a struct principal do servidor
type Server struct {
	router      *gin.Engine
	ecopontoHdl *ecoponto.Handler
	jwtSecret   string
}

// NewServer cria e configura o servidor com todas as rotas
func NewServer(ecopontoHdl *ecoponto.Handler, jwtSecret string) *Server {
	// Cria o router Gin
	r := gin.Default()

	// Cria a struct do servidor
	s := &Server{
		router:      r,
		ecopontoHdl: ecopontoHdl,
		jwtSecret:   jwtSecret,
	}

	// Registra todas as nossas rotas
	s.registerRoutes()

	return s
}

// registerRoutes é um método privado para organizar o registro
func (s *Server) registerRoutes() {
	// --- Rotas Públicas  ---
	// Health Check
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "pong",
			"db_status": "conectado",
		})
	})

	apiPublic := s.router.Group("/api")
	{
		apiPublic.GET("/ecopontos", s.ecopontoHdl.ListEcopontos)
		apiPublic.GET("/ecopontos/:id", s.ecopontoHdl.GetEcoponto)
	}

	// --- Rotas de Admin (protegidas com JWT) ---
	apiAdmin := s.router.Group("/api")

	// Aplica o middleware de autenticação a este grupo
	apiAdmin.Use(AuthMiddleware(s.jwtSecret))
	{
		apiAdmin.POST("/ecopontos", s.ecopontoHdl.CreateEcoponto)
	}
}

// Run inicia o servidor HTTP
func (s *Server) Run(port string) error {
	addr := fmt.Sprintf(":%s", port)
	if port == "" {
		addr = ":8080" // Fallback
	}

	fmt.Printf("Servidor subindo na porta %s\n", addr)
	return s.router.Run(addr)
}
