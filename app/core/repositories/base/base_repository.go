package base

import "context"

type Repository[T any] interface {
	List(ctx context.Context) ([]T, error)
	Get(ctx context.Context, id string) (*T, error)
	Create(ctx context.Context, v *T) (*T, error)
	Update(ctx context.Context, v *T) (*T, error)
	Delete(ctx context.Context, id string) error
}
