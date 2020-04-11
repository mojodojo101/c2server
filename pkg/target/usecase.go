package target

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Usecase interface {
	GetByID(ctx context.Context, id int64) (*models.Target, error)
	GetByIpv4(ctx context.Context, ipv4 string) (*models.Target, error)
	Delete(ctx context.Context, t *models.Target) error
	Update(ctx context.Context, t *models.Target) error
	FetchCmdResponse(ctx context.Context, t *models.Target, cmdId int64) ([]byte, error)
	GetNextCmd(ctx context.Context, t *models.Target) (int64, string, error)
	SetCmdExecuted(ctx context.Context, t *models.Target, cmdId int64, response []byte) error
	Store(ctx context.Context, t *models.Target) error
	CreateTable(ctx context.Context) error
}
