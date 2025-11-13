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
func (r *Repository) Create(ctx context.Context, req CreateEcoPontoRequest, lat, lon float64) (*EcoPonto, error) {
	query := `
		INSERT INTO ecopontos (
			nome, tipo_residuo, logradouro, bairro, coordenadas, 
			horario_funcionamento, foto_url
		)
		VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($5, $6), 4326), $7, $8)
		RETURNING id, nome, tipo_residuo, logradouro, bairro, created_at, 
				  horario_funcionamento, foto_url
	`

	var novoPonto EcoPonto
	err := r.db.QueryRowxContext(
		ctx,
		query,
		req.Nome,
		req.TipoResiduo,
		req.Logradouro,
		req.Bairro,
		lon,
		lat,
		req.HorarioFuncionamento,
		req.FotoURL,              
	).StructScan(&novoPonto)

	if err != nil {
		return nil, err
	}
	novoPonto.Latitude = lat
	novoPonto.Longitude = lon
	return &novoPonto, nil
}

// Adiciona os novos campos ao SELECT
func (r *Repository) ListByProximity(ctx context.Context, params ListByProximityParams) ([]EcoPonto, error) {
	query := `
		SELECT 
			id, nome, tipo_residuo, logradouro, bairro, created_at,
			horario_funcionamento, foto_url, -- <--- ADICIONADO
			ST_X(coordenadas::geometry) AS longitude,
			ST_Y(coordenadas::geometry) AS latitude
		FROM 
			ecopontos
		WHERE 
			ST_DWithin(
				coordenadas,
				ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
				$3
			)
			AND ($4 = '' OR tipo_residuo = $4)
	`
	var pontos []EcoPonto
	err := r.db.SelectContext(ctx, &pontos, query, params.Longitude, params.Latitude, params.Distancia, params.TipoResiduo)
	if err != nil {
		return nil, err
	}
	if pontos == nil {
		pontos = make([]EcoPonto, 0)
	}
	return pontos, nil
}

// Adiciona os novos campos ao SELECT
func (r *Repository) GetByID(ctx context.Context, id string) (*EcoPonto, error) {
	query := `
		SELECT 
			id, nome, tipo_residuo, logradouro, bairro, created_at,
			horario_funcionamento, foto_url, -- <--- ADICIONADO
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

// Adiciona os novos campos ao builder do squirrel 
func (r *Repository) Update(ctx context.Context, id string, req UpdateEcoPontoRequest) (*EcoPonto, error) {
	// 1. Inicia o construtor de SQL
	qb := sq.Update("ecopontos").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, nome, tipo_residuo, logradouro, bairro, created_at, horario_funcionamento, foto_url, ST_X(coordenadas::geometry) AS longitude, ST_Y(coordenadas::geometry) AS latitude") 

	// 2. Adiciona campos de texto
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

	if req.HorarioFuncionamento != nil {
		qb = qb.Set("horario_funcionamento", *req.HorarioFuncionamento)
	}
	if req.FotoURL != nil {
		qb = qb.Set("foto_url", *req.FotoURL)
	}

	// 3. Lógica especial para coordenadas 
	if req.Latitude != nil && req.Longitude != nil {
		qb = qb.Set("coordenadas", sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)", *req.Longitude, *req.Latitude))
	} else if req.Latitude != nil || req.Longitude != nil {
		// Se o admin só enviar Lat ou Lon, é um erro
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

func (r *Repository) ListAll(ctx context.Context) ([]EcoPonto, error) {
	query := `
		SELECT 
			id, nome, tipo_residuo, logradouro, bairro, created_at,
			horario_funcionamento, foto_url, -- <--- ADICIONADO
			ST_X(coordenadas::geometry) AS longitude,
			ST_Y(coordenadas::geometry) AS latitude
		FROM 
			ecopontos
		ORDER BY
			created_at DESC
	`
	var pontos []EcoPonto
	err := r.db.SelectContext(ctx, &pontos, query)
	if err != nil {
		return nil, err
	}
	if pontos == nil {
		pontos = make([]EcoPonto, 0)
	}
	return pontos, nil
}
