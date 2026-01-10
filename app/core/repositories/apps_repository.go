package repositories

import (
	"context"

	"github.com/brunojet/my-store-go/app/core/models"
	"github.com/brunojet/my-store-go/app/core/repositories/ports"
	"gorm.io/gorm"
)

// Repository is the application repository for Apps.
//
// Conventions:
//   - Get/Patch return (nil, nil) when the record does not exist.
//   - Delete is idempotent.
type Repository interface {
	ports.CRUDRepository[models.App]
}

type appsRepository struct {
	crud *ports.GormCRUDRepository[models.App]
}

func New(db *gorm.DB) Repository {
	return NewAppsRepository(db)
}

func NewAppsRepository(db *gorm.DB) Repository {
	return &appsRepository{
		crud: ports.NewGormCRUDRepository[models.App](db),
	}
}

func (r *appsRepository) List(ctx context.Context) ([]models.App, error) {
	return r.crud.List(ctx)
}

func (r *appsRepository) Get(ctx context.Context, id uint) (*models.App, error) {
	return r.crud.Get(ctx, id)
}

func (r *appsRepository) Create(ctx context.Context, v *models.App) (*models.App, error) {
	return r.crud.Create(ctx, v)
}

func (r *appsRepository) Delete(ctx context.Context, id uint) error {
	return r.crud.Delete(ctx, id)
}

func (r *appsRepository) Patch(ctx context.Context, id uint, updates map[string]any) (*models.App, error) {
	return r.crud.Patch(ctx, id, updates)
}
