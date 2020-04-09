package target

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Usecase interface {
	GetByID(ctx context.Context, id int64) (*models.Target, error)
	Delete(ctx context.Context, t *models.Target) error
	Update(ctx context.Context, t *models.Target) error
	FetchCmdResponse(ctx context.Context, t *models.Target, cmdId int64) ([]byte, error)
	Store(ctx context.Context, t *models.Target) (int64, error)
	CreateTable(ctx context.Context) error
}
