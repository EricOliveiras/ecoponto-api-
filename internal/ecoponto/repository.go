package ecoponto

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Repository gerencia a persistência dos ecopontos
type Repository struct {
	db *sqlx.DB
}

// NewRepository cria uma nova instância do repositório
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Create insere um novo ecoponto no banco de dados
func (r *Repository) Create(ctx context.Context, req CreateEcoPontoRequest) (*EcoPonto, error) {
	query := `
		INSERT INTO ecopontos (nome, tipo_residuo, logradouro, bairro, coordenadas)
		VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($5, $6), 4326))
		RETURNING id, nome, tipo_residuo, logradouro, bairro, created_at
	`
	var novoPonto EcoPonto
	err := r.db.QueryRowxContext(
		ctx,
		query,
		req.Nome,
		req.TipoResiduo,
		req.Logradouro,
		req.Bairro,
		req.Longitude,
		req.Latitude,
	).StructScan(&novoPonto)
	if err != nil {
		return nil, err
	}
	novoPonto.Latitude = req.Latitude
	novoPonto.Longitude = req.Longitude
	return &novoPonto, nil
}

// ListByProximity busca ecopontos dentro de um raio (distância)
func (r *Repository) ListByProximity(ctx context.Context, params ListByProximityParams) ([]EcoPonto, error) {

	query := `
		SELECT 
			id, nome, tipo_residuo, logradouro, bairro, created_at,
			ST_X(coordenadas::geometry) AS longitude,
			ST_Y(coordenadas::geometry) AS latitude
		FROM 
			ecopontos
		WHERE 
			ST_DWithin(
				coordenadas,
				-- AQUI ESTÁ A CORREÇÃO: 4236 -> 4326
				ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
				$3
			)
			AND ($4 = '' OR tipo_residuo = $4)
	`

	var pontos []EcoPonto

	err := r.db.SelectContext(
		ctx,
		&pontos,
		query,
		params.Longitude,   // $1
		params.Latitude,    // $2
		params.Distancia,   // $3
		params.TipoResiduo, // $4
	)

	if err != nil {
		return nil, err
	}

	if pontos == nil {
		pontos = make([]EcoPonto, 0)
	}

	return pontos, nil
}

// GetByID busca um único ecoponto pelo seu ID (UUID)
func (r *Repository) GetByID(ctx context.Context, id string) (*EcoPonto, error) {
	query := `
		SELECT 
			id, nome, tipo_residuo, logradouro, bairro, created_at,
			ST_X(coordenadas::geometry) AS longitude,
			ST_Y(coordenadas::geometry) AS latitude
		FROM 
			ecopontos
		WHERE 
			id = $1
	`

	var ponto EcoPonto

	err := r.db.GetContext(ctx, &ponto, query, id)

	if err != nil {
		return nil, err
	}

	return &ponto, nil
}
