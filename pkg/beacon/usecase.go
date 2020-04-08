package beacon

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Usecase interface {
	GetByID(ctx context.Context, id int64) (*models.Beacon, error)
	Delete(ctx context.Context, b *models.Beacon) error
	Store(ctx context.Context, b *models.Beacon) error
	CreateTable(ctx context.Context) error
	RetrieveBeaconBuffer(ctx context.Context, b *models.Beacon) ([]byte, error)
}
