package activebeacon

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (*models.ActiveBeacon, error)
	GetByParentID(ctx context.Context, pId int64) (*models.ActiveBeacon, error)
	DeleteByID(ctx context.Context, id int64) error
	CreateNewBeacon(ctx context.Context, m *models.ActiveBeacon) error
	Update(ctx context.Context, m *models.ActiveBeacon) error
	GetAllActiveBeacons(ctx context.Context, amount int64) ([]models.ActiveBeacon, error)
	CreateTable(ctx context.Context) error
	DropTable(ctx context.Context) error
}
