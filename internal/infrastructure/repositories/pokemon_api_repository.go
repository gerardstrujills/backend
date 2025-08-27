package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"time"

	"github.com/gerardstrujills/backend/internal/domain/entities"
)

type PokemonAPIRepository struct {
	baseURL    string
	httpClient *http.Client
}

func NewPokemonAPIRepository(baseURL string) *PokemonAPIRepository {
	return &PokemonAPIRepository{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *PokemonAPIRepository) GetByID(ctx context.Context, id int) (*entities.Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%d", r.baseURL, id)
	return r.fetchPokemon(ctx, url)
}

func (r *PokemonAPIRepository) GetByName(ctx context.Context, name string) (*entities.Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", r.baseURL, strings.ToLower(name))
	return r.fetchPokemon(ctx, url)
}

func (r *PokemonAPIRepository) GetList(ctx context.Context, limit, offset int) (*entities.PokemonList, error) {
	url := fmt.Sprintf("%s/pokemon?limit=%d&offset=%d", r.baseURL, limit, offset)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear la solicitud: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la lista de pokemon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API return status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el cuerpo de la respuesta: %w", err)
	}

	var pokemonList entities.PokemonList
	if err := json.Unmarshal(body, &pokemonList); err != nil {
		return nil, fmt.Errorf("no se pudo serializar la lista de pokemon: %w", err)
	}

	return &pokemonList, nil
}

func (r *PokemonAPIRepository) SearchByTitle(ctx context.Context, title string, limit, offset int) ([]*entities.Pokemon, error) {
	// Estrategia optimizada: usar caché inteligente + búsqueda incremental
	searchTerm := strings.ToLower(title)

	// 1. Intentar con lista pequeña primero (más común)
	initialLimit := 100
	pokemonList, err := r.GetList(ctx, initialLimit, 0)
	if err != nil {
		return nil, err
	}

	var results []*entities.Pokemon
	var candidates []string

	// 2. Filtrar candidatos por nombre
	for _, pokemonResult := range pokemonList.Results {
		if strings.Contains(strings.ToLower(pokemonResult.Name), searchTerm) {
			candidates = append(candidates, pokemonResult.Name)
		}
	}

	// 3. Si no hay suficientes resultados, expandir búsqueda
	if len(candidates) < limit+offset && pokemonList.Count > initialLimit {
		// Expandir a 500 máximo (balance entre eficiencia y cobertura)
		expandedList, err := r.GetList(ctx, 500, 0)
		if err == nil {
			candidates = candidates[:0] // Reset
			for _, pokemonResult := range expandedList.Results {
				if strings.Contains(strings.ToLower(pokemonResult.Name), searchTerm) {
					candidates = append(candidates, pokemonResult.Name)
				}
			}
		}
	}

	// 4. Aplicar paginación a candidatos
	start := offset
	end := offset + limit
	if start >= len(candidates) {
		return results, nil
	}
	if end > len(candidates) {
		end = len(candidates)
	}

	// 5. Obtener detalles solo de los Pokemon necesarios
	for i := start; i < end; i++ {
		pokemon, err := r.GetByName(ctx, candidates[i])
		if err != nil {
			continue // Skip errores individuales
		}
		results = append(results, pokemon)
	}

	return results, nil
}

func (r *PokemonAPIRepository) fetchPokemon(ctx context.Context, url string) (*entities.Pokemon, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear la solicitud: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener el pokemon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("pokemon no encontrado")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API return status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el cuerpo de la respuesta: %w", err)
	}

	var pokemon entities.Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return nil, fmt.Errorf("no se pudo serializar a los pokemon: %w", err)
	}

	return &pokemon, nil
}
