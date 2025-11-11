package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ericoliveiras/ecoponto-api/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	AdminEmail = "admin@ecoponto.com"
	AdminPass  = "admin123"
)

func main() {
	log.Println("Iniciando o seeder de admin...")

	// 1. Carrega o .env local
	if err := godotenv.Load(); err != nil {
		log.Printf("Aviso: Erro ao carregar arquivo .env: %v", err)
	}

	// 2. Construir a DATABASE_URL localmente
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	if user == "" || pass == "" || port == "" || dbName == "" {
		log.Fatalf("Erro: Variáveis de banco (POSTGRES_USER, POSTGRES_PASSWORD, DB_PORT, POSTGRES_DB) não encontradas no .env")
	}

	// Esta URL aponta para localhost, que é o correto para o seeder
	localDbURL := fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s?sslmode=disable",
		user, pass, port, dbName,
	)

	// 3. Conectar ao Banco
	db := database.Connect(localDbURL)
	defer db.Close()

	// 4. Gerar o Hash da Senha
	log.Printf("Gerando hash para a senha '%s'...", AdminPass)
	hash, err := bcrypt.GenerateFromPassword([]byte(AdminPass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Erro ao gerar hash: %v", err)
	}
	hashString := string(hash)
	log.Println("Hash gerado com sucesso.")

	// 5. Inserir o usuário no banco
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING
	`

	res, err := db.Exec(query, AdminEmail, hashString)
	if err != nil {
		log.Fatalf("Erro ao inserir admin: %v", err)
	}

	// 6. Checar o resultado
	rows, _ := res.RowsAffected()
	if rows == 0 {
		log.Printf("Usuário '%s' já existia. Nada foi feito.", AdminEmail)
	} else {
		log.Printf("Usuário admin '%s' criado com sucesso!", AdminEmail)
	}
}
