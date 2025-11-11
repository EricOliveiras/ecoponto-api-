package ecoponto

import "time"

// EcoPonto representa a estrutura de dados de um ponto de coleta.
type EcoPonto struct {
	ID          string    `db:"id" json:"id"`
	Logradouro  string    `db:"logradouro" json:"logradouro"`
	Bairro      string    `db:"bairro" json:"bairro"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Nome        string    `db:"nome" json:"nome"`
	TipoResiduo string    `db:"tipo_residuo" json:"tipo_residuo"`
	Latitude    float64   `db:"latitude" json:"latitude"`
	Longitude   float64   `db:"longitude" json:"longitude"`
}

type CreateEcoPontoRequest struct {
	Nome        string  `json:"nome" binding:"required"`
	TipoResiduo string  `json:"tipo_residuo" binding:"required"`
	Logradouro  string  `json:"logradouro"`
	Bairro      string  `json:"bairro"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
}

type ListByProximityParams struct {
	Latitude    float64
	Longitude   float64
	Distancia   int
	TipoResiduo string
}

type UpdateEcoPontoRequest struct {
	Nome        *string  `json:"nome"`
	TipoResiduo *string  `json:"tipo_residuo"`
	Logradouro  *string  `json:"logradouro"`
	Bairro      *string  `json:"bairro"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
}
