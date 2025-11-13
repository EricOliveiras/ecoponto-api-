package ecoponto

import "time"

type EcoPonto struct {
	ID                   string    `db:"id" json:"id"`
	Logradouro           string    `db:"logradouro" json:"logradouro"`
	Bairro               string    `db:"bairro" json:"bairro"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	Nome                 string    `db:"nome" json:"nome"`
	TipoResiduo          string    `db:"tipo_residuo" json:"tipo_residuo"`
	Latitude             float64   `db:"latitude" json:"latitude"`
	Longitude            float64   `db:"longitude" json:"longitude"`
	HorarioFuncionamento *string   `db:"horario_funcionamento" json:"horario_funcionamento,omitempty"`
	FotoURL              *string   `db:"foto_url" json:"foto_url,omitempty"`
}

type CreateEcoPontoRequest struct {
	Nome                 string   `json:"nome" binding:"required"`
	TipoResiduo          string   `json:"tipo_residuo" binding:"required"`
	Logradouro           string   `json:"logradouro" binding:"required"`
	Bairro               string   `json:"bairro" binding:"required"`
	Cidade               string   `json:"cidade" binding:"required"`
	Estado               string   `json:"estado" binding:"required"`
	Latitude             *float64 `json:"latitude"`
	Longitude            *float64 `json:"longitude"`
	HorarioFuncionamento *string  `json:"horario_funcionamento"`
	FotoURL              *string  `json:"foto_url"`
}

type UpdateEcoPontoRequest struct {
	Nome                 *string  `json:"nome"`
	TipoResiduo          *string  `json:"tipo_residuo"`
	Logradouro           *string  `json:"logradouro"`
	Bairro               *string  `json:"bairro"`
	Latitude             *float64 `json:"latitude"`
	Longitude            *float64 `json:"longitude"`
	HorarioFuncionamento *string  `json:"horario_funcionamento"`
	FotoURL              *string  `json:"foto_url"`
}

type ListByProximityParams struct {
	Latitude    float64
	Longitude   float64
	Distancia   int
	TipoResiduo string
}
