package ecoponto

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
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

func (r *Repository) Update(ctx context.Context, id string, req UpdateEcoPontoRequest) (*EcoPonto, error) {
	// 1. Inicia o construtor de SQL (squirrel)
	// Usamos PlaceholderFormat(sq.Dollar) para usar $1, $2...
	qb := sq.Update("ecopontos").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, nome, tipo_residuo, logradouro, bairro, created_at, ST_X(coordenadas::geometry) AS longitude, ST_Y(coordenadas::geometry) AS latitude")

	// 2. Adiciona campos à query APENAS se eles não forem nil
	if req.Nome != nil {
		qb = qb.Set("nome", *req.Nome)
	}
	if req.TipoResiduo != nil {
		qb = qb.Set("tipo_residuo", *req.TipoResiduo)
	}
	if req.Logradouro != nil {
		qb = qb.Set("logradouro", *req.Logradouro)
	}
	if req.Bairro != nil {
		qb = qb.Set("bairro", *req.Bairro)
	}

	// 3. Lógica especial para coordenadas (devem ser atualizadas juntas)
	if req.Latitude != nil && req.Longitude != nil {
		// Criamos uma expressão SQL customizada
		qb = qb.Set("coordenadas", sq.Expr("ST_SetSRID(ST_MakePoint($?), $?)", *req.Longitude, *req.Latitude, 4326))
	} else if req.Latitude != nil || req.Longitude != nil {
		// Se o admin só enviar Lat ou Lon, é um erro de lógica
		return nil, sql.ErrTxDone
	}

	// 4. Constrói e executa a query
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var pontoAtualizado EcoPonto
	err = r.db.QueryRowxContext(ctx, query, args...).StructScan(&pontoAtualizado)
	if err != nil {
		// Pode ser sql.ErrNoRows se o ID não existir
		return nil, err
	}

	return &pontoAtualizado, nil
}

// Delete remove um ecoponto do banco de dados pelo seu ID
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM ecopontos WHERE id = $1"

	// Executa a query
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Verifica se alguma linha foi realmente deletada
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Se nenhuma linha foi afetada, significa que o ID não existia
	if rows == 0 {
		return sql.ErrNoRows 
	}

	return nil
}
