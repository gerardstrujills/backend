package routes

import (
	"github.com/gerardstrujills/backend/internal/interfaces/http/handlers"
	"github.com/gerardstrujills/backend/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, pokemonHandler *handlers.PokemonHandler) {
	// Middleware global
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "pokemon-backend",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Pokemon routes
		pokemon := v1.Group("/pokemon")
		{
			pokemon.GET("", pokemonHandler.GetPokemonList)              // GET /api/v1/pokemon?limit=20&offset=0
			pokemon.GET("/search", pokemonHandler.SearchPokemonByTitle) // GET /api/v1/pokemon/search?q=pika&limit=10&offset=0
			pokemon.GET("/:id", pokemonHandler.GetPokemonByID)          // GET /api/v1/pokemon/25
			pokemon.GET("/name/:name", pokemonHandler.GetPokemonByName) // GET /api/v1/pokemon/name/pikachu
		}
	}
}
