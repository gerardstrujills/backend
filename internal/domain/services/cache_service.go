package services

import "context"

// Operaciones de caché
type CacheService interface {
	Get(ctx context.Context, key string) (interface{}, bool)
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}
