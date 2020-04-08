package target

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Usecase interface {
	GetByID(ctx context.Context, id int64) (*models.Target, error)
	Delete(ctx context.Context, b *models.Target) error
	Update(ctx context.Context, b *models.Target) error
	Store(ctx context.Context, b *models.Target) error
	CreateTable(ctx context.Context) error
}
