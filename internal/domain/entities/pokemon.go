package entities

import "time"

// Entidad principal del dominio
type Pokemon struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Height  int     `json:"height"`
	Weight  int     `json:"weight"`
	Types   []Type  `json:"types"`
	Sprites Sprites `json:"sprites"`
	BaseExp int     `json:"base_experience"`
}

type Type struct {
	Slot int      `json:"slot"`
	Type TypeInfo `json:"type"`
}

type TypeInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Sprites struct {
	FrontDefault string `json:"front_default"`
	BackDefault  string `json:"back_default"`
}

// Lista paginada de Pokemon
type PokemonList struct {
	Count    int             `json:"count"`
	Next     *string         `json:"next"`
	Previous *string         `json:"previous"`
	Results  []PokemonResult `json:"results"`
}

type PokemonResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Elemento en caché con TTL
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// Verifica si el item del caché ha expirado
func (c *CacheItem) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}
