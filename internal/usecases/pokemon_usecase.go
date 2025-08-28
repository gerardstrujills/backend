package usecases

import (
	"context"
	"fmt"

	"strings"

	"github.com/gerardstrujills/backend/internal/domain/entities"
	"github.com/gerardstrujills/backend/internal/domain/repositories"
	"github.com/gerardstrujills/backend/internal/domain/services"
)

type PokemonUseCase struct {
	pokemonRepo  repositories.PokemonRepository
	cacheService services.CacheService
}

func NewPokemonUseCase(pokemonRepo repositories.PokemonRepository, cacheService services.CacheService) *PokemonUseCase {
	return &PokemonUseCase{
		pokemonRepo:  pokemonRepo,
		cacheService: cacheService,
	}
}

// GetPokemonByID obtiene un Pokemon por ID con cache
func (uc *PokemonUseCase) GetPokemonByID(ctx context.Context, id int) (*entities.Pokemon, error) {
	cacheKey := fmt.Sprintf("pokemon:id:%d", id)

	// Intentar obtener del cache primero
	if cachedData, found := uc.cacheService.Get(ctx, cacheKey); found {
		if pokemon, ok := cachedData.(*entities.Pokemon); ok {
			return pokemon, nil
		}
	}

	// Si no esta en cache, obtener del repositorio
	pokemon, err := uc.pokemonRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener pokemon por ID: %d: %w", id, err)
	}

	// Guardar en cache
	if err := uc.cacheService.Set(ctx, cacheKey, pokemon); err != nil {
		// Registrar error pero no fallar en la solicitud
		fmt.Printf("no se pudo almacenar en cache el pokemon: %d: %v\n", id, err)
	}

	return pokemon, nil
}

// GetPokemonByName obtiene un Pokemon por nombre con cache
func (uc *PokemonUseCase) GetPokemonByName(ctx context.Context, name string) (*entities.Pokemon, error) {
	cacheKey := fmt.Sprintf("pokemon:name:%s", strings.ToLower(name))

	// Intentar obtener del cache primero
	if cachedData, found := uc.cacheService.Get(ctx, cacheKey); found {
		if pokemon, ok := cachedData.(*entities.Pokemon); ok {
			return pokemon, nil
		}
	}

	// Si no esta en cache, obtener del repositorio
	pokemon, err := uc.pokemonRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("no pude obtener el pokemon por su nombre: %s: %w", name, err)
	}

	// Guardar en cache
	if err := uc.cacheService.Set(ctx, cacheKey, pokemon); err != nil {
		fmt.Printf("no se pudo almacenar en cache el pokemon: %s: %v\n", name, err)
	}

	return pokemon, nil
}

// GetPokemonList obtiene una lista paginada de Pokemon con cache
func (uc *PokemonUseCase) GetPokemonList(ctx context.Context, limit, offset int) (*entities.PokemonList, error) {
	cacheKey := fmt.Sprintf("pokemon:list:%d:%d", limit, offset)

	// Intentar obtener del cache primero
	if cachedData, found := uc.cacheService.Get(ctx, cacheKey); found {
		if pokemonList, ok := cachedData.(*entities.PokemonList); ok {
			return pokemonList, nil
		}
	}

	// Si no esta en cache, obtener del repositorio
	pokemonList, err := uc.pokemonRepo.GetList(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la lista de pokemon: %w", err)
	}

	// Guardar en cache
	if err := uc.cacheService.Set(ctx, cacheKey, pokemonList); err != nil {
		fmt.Printf("no se pudo almacenar en cache el pokemon list: %v\n", err)
	}

	return pokemonList, nil
}

// SearchPokemonByTitle busca Pokemon por titulo/nombre con cache
func (uc *PokemonUseCase) SearchPokemonByTitle(ctx context.Context, title string, limit, offset int) ([]*entities.Pokemon, error) {
	searchTerm := strings.ToLower(title)
	cacheKey := fmt.Sprintf("pokemon:search:%s:%d:%d", searchTerm, limit, offset)

	// Intentar obtener del cache primero
	if cachedData, found := uc.cacheService.Get(ctx, cacheKey); found {
		if pokemonList, ok := cachedData.([]*entities.Pokemon); ok {
			return pokemonList, nil
		}
	}

	// Cache adicional para candidatos de busqueda (evita re-filtrar)
	candidatesCacheKey := fmt.Sprintf("pokemon:search_candidates:%s", searchTerm)
	var candidates []string

	if cachedCandidates, found := uc.cacheService.Get(ctx, candidatesCacheKey); found {
		if candidatesList, ok := cachedCandidates.([]string); ok {
			candidates = candidatesList
		}
	}

	// Si no hay candidatos en cache, buscar
	if len(candidates) == 0 {
		pokemonList, err := uc.pokemonRepo.SearchByTitle(ctx, title, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("no se pudo buscar Pokémon por titulo: %s: %w", title, err)
		}

		// Cachear resultado final
		if err := uc.cacheService.Set(ctx, cacheKey, pokemonList); err != nil {
			fmt.Printf("no se pudo almacenar en cache el pokemon search: %v\n", err)
		}

		return pokemonList, nil
	}

	// Usar candidatos cacheados para paginación eficiente
	var results []*entities.Pokemon
	start := offset
	end := offset + limit
	if start >= len(candidates) {
		return results, nil
	}
	if end > len(candidates) {
		end = len(candidates)
	}

	for i := start; i < end; i++ {
		pokemon, err := uc.GetPokemonByName(ctx, candidates[i]) // Usa cache individual
		if err != nil {
			continue
		}
		results = append(results, pokemon)
	}

	// Cachear resultado paginado
	if err := uc.cacheService.Set(ctx, cacheKey, results); err != nil {
		fmt.Printf("no se pudo almacenar en cache el pokemon search: %v\n", err)
	}

	return results, nil
}
