package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Connect abre a conexão com o banco de dados e a retorna
func Connect(databaseURL string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Testa a conexão
	if err = db.Ping(); err != nil {
		log.Fatalf("Erro ao 'pingar' o banco de dados: %v", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso.")
	return db
}
