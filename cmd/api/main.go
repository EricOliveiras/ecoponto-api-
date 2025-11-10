package main

import (
	"log"

	"github.com/ericoliveiras/ecoponto-api/internal/config"
	"github.com/ericoliveiras/ecoponto-api/internal/database"
	"github.com/ericoliveiras/ecoponto-api/internal/ecoponto"
	"github.com/ericoliveiras/ecoponto-api/internal/server"
)

func main() {
	// 1. Carrega as configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	// 2. Conecta ao banco
	db := database.Connect(cfg.DatabaseURL)
	defer db.Close()

	// 3. Injeção de Dependência
	ecopontoRepo := ecoponto.NewRepository(db)
	ecopontoHandler := ecoponto.NewHandler(ecopontoRepo)

	// 4. Configura o Servidor
	srv := server.NewServer(ecopontoHandler, cfg.JWTSecret)

	// 5. Sobe o servidor
	if err := srv.Run(cfg.APIPort); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
