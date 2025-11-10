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

	// 2. Conecta ao banco de dados 
	db := database.Connect(cfg.DatabaseURL)
	defer db.Close()

	// 3. Injeção de Dependência (a "amarração")
	// Criamos as instâncias que nosso servidor precisa
	ecopontoRepo := ecoponto.NewRepository(db)
	ecopontoHandler := ecoponto.NewHandler(ecopontoRepo)

	// 4. Configura o Servidor
	// Passamos o handler para o novo servidor
	srv := server.NewServer(ecopontoHandler)

	// 5. Sobe o servidor
	// Chamamos o método Run do nosso novo servidor
	if err := srv.Run(cfg.APIPort); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
