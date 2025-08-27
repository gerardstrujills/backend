package handlers

import (
	"net/http"
	"strconv"

	"github.com/gerardstrujills/backend/internal/usecases"
	"github.com/gin-gonic/gin"
)

type PokemonHandler struct {
	pokemonUseCase *usecases.PokemonUseCase
}

func NewPokemonHandler(pokemonUseCase *usecases.PokemonUseCase) *PokemonHandler {
	return &PokemonHandler{
		pokemonUseCase: pokemonUseCase,
	}
}

// GetPokemonByID maneja GET /pokemon/:id
func (h *PokemonHandler) GetPokemonByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid pokemon ID",
		})
		return
	}

	pokemon, err := h.pokemonUseCase.GetPokemonByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "pokemon not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Pokemon not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get pokemon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   pokemon,
		"cached": false, // Podrías implementar un flag para indicar si vino del caché
	})
}

// GetPokemonByName maneja GET /pokemon/name/:name
func (h *PokemonHandler) GetPokemonByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Pokemon name is required",
		})
		return
	}

	pokemon, err := h.pokemonUseCase.GetPokemonByName(c.Request.Context(), name)
	if err != nil {
		if err.Error() == "pokemon not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Pokemon not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get pokemon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": pokemon,
	})
}

// GetPokemonList maneja GET /pokemon con paginación
func (h *PokemonHandler) GetPokemonList(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	pokemonList, err := h.pokemonUseCase.GetPokemonList(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get pokemon list",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": pokemonList,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
		},
	})
}

// SearchPokemonByTitle maneja GET /pokemon/search con paginación
func (h *PokemonHandler) SearchPokemonByTitle(c *gin.Context) {
	title := c.Query("q")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query 'q' is required",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	pokemonList, err := h.pokemonUseCase.SearchPokemonByTitle(c.Request.Context(), title, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search pokemon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": pokemonList,
		"search": gin.H{
			"query":  title,
			"limit":  limit,
			"offset": offset,
			"count":  len(pokemonList),
		},
	})
}
