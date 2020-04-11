package activebeacon

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Usecase interface {
	Delete(ctx context.Context, a *models.ActiveBeacon) error
	GetByID(ctx context.Context, id int64) (*models.ActiveBeacon, error)
	GetTargetByIpv4(ctx context.Context, ipv4 string) (*models.Target, error)
	Update(ctx context.Context, a *models.ActiveBeacon) error
	SetCmdExecuted(ctx context.Context, a *models.ActiveBeacon, response []byte) error
	GetNextCommand(ctx context.Context, a *models.ActiveBeacon) error
	Register(ctx context.Context, a *models.ActiveBeacon) error
	CreateTable(ctx context.Context) error
}
