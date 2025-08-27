package main

import (
	"log"

	"time"

	"github.com/gerardstrujills/backend/internal/infrastructure/cache"
	"github.com/gerardstrujills/backend/internal/infrastructure/repositories"
	"github.com/gerardstrujills/backend/internal/interfaces/http/handlers"
	"github.com/gerardstrujills/backend/internal/interfaces/http/routes"
	"github.com/gerardstrujills/backend/internal/usecases"
	"github.com/gin-gonic/gin"
)

func main() {
	// Configuración
	const (
		serverPort = ":8080"
		pokeAPIURL = "https://pokeapi.co/api/v2"
		cacheSize  = 1000
		cacheTTL   = 15 * time.Minute
	)

	// Inicializar caché LRU
	cacheService, err := cache.NewLRUCache(cacheSize, cacheTTL)
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	// Inicializar repositorio
	pokemonRepo := repositories.NewPokemonAPIRepository(pokeAPIURL)

	// Inicializar casos de uso
	pokemonUseCase := usecases.NewPokemonUseCase(pokemonRepo, cacheService)

	// Inicializar handlers
	pokemonHandler := handlers.NewPokemonHandler(pokemonUseCase)

	// Configurar Gin
	if gin.Mode() == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Configurar rutas
	routes.SetupRoutes(r, pokemonHandler)

	log.Printf("🚀 Pokemon Backend Server starting on port %s", serverPort)
	log.Printf("📋 Available endpoints:")
	log.Printf("   GET /health")
	log.Printf("   GET /api/v1/pokemon?limit=20&offset=0")
	log.Printf("   GET /api/v1/pokemon/:id")
	log.Printf("   GET /api/v1/pokemon/name/:name")
	log.Printf("   GET /api/v1/pokemon/search?q=pika&limit=10&offset=0")

	if err := r.Run(serverPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
