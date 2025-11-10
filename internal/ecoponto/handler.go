package ecoponto

import (
	"net/http"
	"strconv" // Precisamos dele para converter os parâmetros da URL

	"github.com/gin-gonic/gin"
)

// Handler (sem mudanças)
type Handler struct {
	repo *Repository
}

// NewHandler (sem mudanças)
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// CreateEcoponto (sem mudanças)
func (h *Handler) CreateEcoponto(c *gin.Context) {
	var req CreateEcoPontoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	novoPonto, err := h.repo.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, novoPonto)
}

// --- NOVO CÓDIGO ABAIXO ---

// ListEcopontos é o método para o endpoint GET /api/ecopontos
func (h *Handler) ListEcopontos(c *gin.Context) {
	// 1. Ler os parâmetros da query (URL)
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	distStr := c.Query("dist")

	// --- NOVO CÓDIGO ---
	tipoStr := c.Query("tipo") // Pega o ?tipo=...

	// 2. Validar parâmetros obrigatórios (lat, lon)
	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros 'lat' e 'lon' são obrigatórios"})
		return
	}

	// 3. Converter strings para números (float64 e int)
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'lat' inválido"})
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'lon' inválido"})
		return
	}

	// 4. Definir uma distância padrão
	dist, err := strconv.Atoi(distStr)
	if err != nil || dist <= 0 {
		dist = 5000 // Default: 5 km
	}

	// 5. Montar os parâmetros para o repositório
	params := ListByProximityParams{
		Latitude:  lat,
		Longitude: lon,
		Distancia: dist,
		// --- NOVO CÓDIGO ---
		TipoResiduo: tipoStr, // Passa o filtro para o repositório
	}

	// 6. Chamar o repositório
	pontos, err := h.repo.ListByProximity(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 7. Retornar a lista de pontos
	c.JSON(http.StatusOK, pontos)
}
