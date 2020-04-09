package activebeacon

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Usecase interface {
	Delete(ctx context.Context, a *models.ActiveBeacon) error
	GetByID(ctx context.Context, id int64) (*models.ActiveBeacon, error)
	Update(ctx context.Context, a *models.ActiveBeacon) error
	Register(ctx context.Context, a *models.ActiveBeacon) (int64, error)
	CreateTable(ctx context.Context) error
}
