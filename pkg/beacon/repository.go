package beacon

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (*models.Beacon, error)
	DeleteByID(ctx context.Context, id int64) error
	CreateNewBeacon(ctx context.Context, b *models.Beacon) (int64, error)
	CreateTable(ctx context.Context) error
	DropTable(ctx context.Context) error
}
