package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ericoliveiras/ecoponto-api/internal/auth"
	"github.com/ericoliveiras/ecoponto-api/internal/ecoponto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server é a struct principal do servidor
type Server struct {
	router      *gin.Engine
	ecopontoHdl *ecoponto.Handler
	authHdl     *auth.Handler
	jwtSecret   string
}

// NewServer cria e configura o servidor com todas as rotas
func NewServer(ecopontoHdl *ecoponto.Handler, authHdl *auth.Handler, jwtSecret string) *Server {
	// Cria o router Gin
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		// Permite que origens específicas façam requisições
		// (Mude para a porta do seu frontend, ex: 5173 para Vite)
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"},

		// Quais métodos HTTP são permitidos
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},

		// Quais headers o frontend pode enviar
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},

		// Permite que o navegador exponha o header "Authorization"
		ExposeHeaders: []string{"Content-Length"},

		// Permite que cookies/credenciais sejam enviados
		AllowCredentials: true,

		// Duração máxima que o "preflight" (OPTIONS) é cacheado
		MaxAge: 12 * time.Hour,
	}))

	// Cria a struct do servidor
	s := &Server{
		router:      r,
		ecopontoHdl: ecopontoHdl,
		authHdl:     authHdl,
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
		apiPublic.POST("/auth/login", s.authHdl.Login)
	}

	// --- Rotas de Admin (protegidas com JWT) ---
	apiAdmin := s.router.Group("/api")

	// Aplica o middleware de autenticação a este grupo
	apiAdmin.Use(AuthMiddleware(s.jwtSecret))
	{
		apiAdmin.POST("/ecopontos", s.ecopontoHdl.CreateEcoponto)
		apiAdmin.PUT("/ecopontos/:id", s.ecopontoHdl.UpdateEcoponto)
		apiAdmin.DELETE("/ecopontos/:id", s.ecopontoHdl.DeleteEcoponto)
		apiAdmin.GET("/ecopontos/all", s.ecopontoHdl.ListAllEcopontos)
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
