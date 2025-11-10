package ecoponto

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler gerencia as requisições HTTP relacionadas aos ecopontos
type Handler struct {
	repo *Repository
}

// NewHandler cria uma nova instância do handler
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// CreateEcoponto é o método para o endpoint POST /api/ecopontos
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

// ListEcopontos é o método para o endpoint GET /api/ecopontos
func (h *Handler) ListEcopontos(c *gin.Context) {
	// 1. Ler os parâmetros da query (URL)
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	distStr := c.Query("dist")

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
		Latitude:    lat,
		Longitude:   lon,
		Distancia:   dist,
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

func (h *Handler) GetEcoponto(c *gin.Context) {
	// 1. Ler o ID do parâmetro da URL
	id := c.Param("id")

	// 2. Chamar o repositório
	ponto, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		// 3. Checar o tipo do erro
		if err == sql.ErrNoRows {
			// Se o banco não achou nada, retornamos 404
			c.JSON(http.StatusNotFound, gin.H{"error": "Ecoponto não encontrado"})
			return
		}

		// Para qualquer outro erro do banco
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Retornar o ponto encontrado
	c.JSON(http.StatusOK, ponto)
}
