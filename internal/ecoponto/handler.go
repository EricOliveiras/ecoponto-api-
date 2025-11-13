package ecoponto

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ericoliveiras/ecoponto-api/internal/geocoding"
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
	// 1. Faz o "bind" do JSON (que agora pode ter lat/lon)
	var req CreateEcoPontoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var lat, lon float64
	var err error

	// Verificamos se os ponteiros NÃO SÃO nulos
	if req.Latitude != nil && req.Longitude != nil {

		log.Println("Recebidas coordenadas manuais do frontend.")
		lat = *req.Latitude
		lon = *req.Longitude

	} else {
		// CENÁRIO 2: O frontend NÃO enviou. Usar Geocoding (como antes).
		log.Println("Coordenadas não fornecidas. A acionar geocoding...")

		// Monta a string de endereço
		addressString := fmt.Sprintf("%s, %s, %s, %s",
			req.Logradouro,
			req.Bairro,
			req.Cidade,
			req.Estado,
		)

		// Chama o Geocoder
		lat, lon, err = geocoding.GetCoordsFromAddress(addressString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Endereço não encontrado ou inválido"})
			return
		}
	}

	// 3. Chama o repositório
	novoPonto, err := h.repo.Create(c.Request.Context(), req, lat, lon)
	if err != nil {
		// Se der erro no banco
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Retorna o objeto criado
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

// UpdateEcoponto é o método para o endpoint PUT /api/ecopontos/:id
func (h *Handler) UpdateEcoponto(c *gin.Context) {
	// 1. Ler o ID do parâmetro da URL
	id := c.Param("id")

	// 2. Define uma variável para o JSON de entrada (UpdateEcoPontoRequest)
	var req UpdateEcoPontoRequest

	// 3. Faz o "bind" do JSON do body para a struct 'req'
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Lógica Condicional de Coordenadas
	// Verificamos se o frontend enviou novas coordenadas (pin arrastado)
	if req.Latitude != nil && req.Longitude != nil {

		log.Println("UPDATE: Recebidas coordenadas manuais do frontend.")
		// O repositório (com squirrel) saberá lidar com estes campos
		// não precisamos de 'lat' e 'lon' separadas aqui.

	} else if req.Logradouro != nil || req.Bairro != nil {
		// CENÁRIO 2: O frontend NÃO enviou coords, mas enviou um NOVO endereço
		// (Precisamos de fazer geocoding neste novo endereço)

		log.Println("UPDATE: A acionar geocoding para novo endereço.")
	}

	// 5. Chama o repositório para atualizar o ecoponto
	// O 'repo.Update' (com Squirrel) é inteligente e irá
	// atualizar apenas os campos que não são nulos no 'req'.
	pontoAtualizado, err := h.repo.Update(c.Request.Context(), id, req)
	if err != nil {
		// 6. Checar os tipos de erro
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ecoponto não encontrado"})
			return
		}
		if err == sql.ErrTxDone {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Latitude e Longitude devem ser enviadas juntas"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 7. Retorna o objeto atualizado
	c.JSON(http.StatusOK, pontoAtualizado)
}

// DeleteEcoponto é o método para o endpoint DELETE /api/ecopontos/:id
func (h *Handler) DeleteEcoponto(c *gin.Context) {
	// 1. Ler o ID do parâmetro da URL
	id := c.Param("id")

	// 2. Chamar o repositório para apagar
	err := h.repo.Delete(c.Request.Context(), id)

	if err != nil {
		// 3. Checar se o ID não existia
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ecoponto não encontrado"})
			return
		}

		// Outro erro de banco
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Retornar 204 No Content (sucesso, sem corpo de resposta)
	c.Status(http.StatusNoContent)
}

// ListAllEcopontos é o método para o endpoint de admin GET /api/ecopontos/all
func (h *Handler) ListAllEcopontos(c *gin.Context) {
	// 1. Chama o repositório
	pontos, err := h.repo.ListAll(c.Request.Context())
	if err != nil {
		// 2. Se der erro, retorna 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Retorna a lista de pontos (mesmo que esteja vazia)
	c.JSON(http.StatusOK, pontos)
}
