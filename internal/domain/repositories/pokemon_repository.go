package repositories

import (
	"context"

	"github.com/gerardstrujills/backend/internal/domain/entities"
)

// Operaciones de acceso a datos
type PokemonRepository interface {
	GetByID(ctx context.Context, id int) (*entities.Pokemon, error)
	GetByName(ctx context.Context, name string) (*entities.Pokemon, error)
	GetList(ctx context.Context, limit, offset int) (*entities.PokemonList, error)
	SearchByTitle(ctx context.Context, title string, limit, offset int) ([]*entities.Pokemon, error)
}
