package target

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (*models.Target, error)
	GetByIpv4(ctx context.Context, ipv4 string) (*models.Target, error)
	DeleteByID(ctx context.Context, id int64) error
	CreateNewTarget(ctx context.Context, t *models.Target) error
	// i changed this from updatecmdid and update hostname just to keep it a bit simpler for now
	// will be worth reimplementing if i actually extend this simple c2
	Update(ctx context.Context, t *models.Target) error
	CreateTable(ctx context.Context) error
	DropTable(ctx context.Context) error
}
