package main

import (
	"log"

	"github.com/ericoliveiras/ecoponto-api/internal/auth"
	"github.com/ericoliveiras/ecoponto-api/internal/config"
	"github.com/ericoliveiras/ecoponto-api/internal/database"
	"github.com/ericoliveiras/ecoponto-api/internal/ecoponto"
	"github.com/ericoliveiras/ecoponto-api/internal/server"
	"github.com/ericoliveiras/ecoponto-api/internal/user"
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

	// 3. Cria os repositórios
	ecopontoRepo := ecoponto.NewRepository(db)
	userRepo := user.NewRepository(db) //

	// Agora criamos os handlers, injetando os repositórios
	ecopontoHandler := ecoponto.NewHandler(ecopontoRepo)
	authHandler := auth.NewHandler(userRepo, cfg.JWTSecret)

	// Passamos o novo authHandler para o servidor
	srv := server.NewServer(ecopontoHandler, authHandler, cfg.JWTSecret)

	// 5. Sobe o servidor
	if err := srv.Run(cfg.APIPort); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
