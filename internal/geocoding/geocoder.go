package geocoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// NominatimResult define a estrutura da resposta da API Nominatim
// A API retorna uma lista, mesmo que seja só um resultado
type NominatimResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// GetCoordsFromAddress converte um endereço em coordenadas
func GetCoordsFromAddress(address string) (float64, float64, error) {
	// A API Nominatim requer um User-Agent.
	// Substitua "EcopontoApp" pelo nome do seu projeto.
	const userAgent = "EcopontoApp/1.0 (oiericoi@hotmail.com)"

	// Monta a URL
	baseURL := "https://nominatim.openstreetmap.org/search"
	query := fmt.Sprintf("?q=%s&format=json&limit=1", url.QueryEscape(address))

	// Cria a requisição
	req, err := http.NewRequest("GET", baseURL+query, nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("User-Agent", userAgent)

	// Faz a chamada HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, 0, fmt.Errorf("API Nominatim retornou status: %s", resp.Status)
	}

	// Lê o JSON da resposta
	var results []NominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return 0, 0, err
	}

	// Verifica se encontrou algum resultado
	if len(results) == 0 {
		return 0, 0, errors.New("endereço não encontrado")
	}

	// Converte as strings Lat/Lon para float64
	lat, errLat := strconv.ParseFloat(results[0].Lat, 64)
	lon, errLon := strconv.ParseFloat(results[0].Lon, 64)
	if errLat != nil || errLon != nil {
		return 0, 0, errors.New("erro ao converter coordenadas")
	}

	return lat, lon, nil
}
