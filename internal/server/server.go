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
}

// NewServer cria e configura o servidor com todas as rotas
func NewServer(ecopontoHdl *ecoponto.Handler) *Server {
	// Cria o router Gin
	r := gin.Default()

	// Cria a struct do servidor
	s := &Server{
		router:      r,
		ecopontoHdl: ecopontoHdl,
	}

	// Registra todas as nossas rotas
	s.registerRoutes()

	return s
}

// registerRoutes é um método privado para organizar o registro
func (s *Server) registerRoutes() {
	// Rota de health-check
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "pong",
			"db_status": "conectado",
		})
	})

	// Agrupa todas as rotas da API sob o prefixo /api
	api := s.router.Group("/api")
	{
		// Rotas de Ecoponto
		api.POST("/ecopontos", s.ecopontoHdl.CreateEcoponto)
		api.GET("/ecopontos", s.ecopontoHdl.ListEcopontos)
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
