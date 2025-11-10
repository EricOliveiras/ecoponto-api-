package config

import (
	"errors"
	"os"
)

// Config armazena todas as configurações da aplicação
type Config struct {
	APIPort     string
	DatabaseURL string
}

// LoadConfig lê a configuração das variáveis de ambiente
func LoadConfig() (Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return Config{}, errors.New("DATABASE_URL não está definida")
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8080"
	}

	config := Config{
		APIPort:     apiPort,
		DatabaseURL: dbURL,
	}

	return config, nil
}
