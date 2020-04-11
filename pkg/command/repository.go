package command

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (*models.Command, error)
	DeleteByID(ctx context.Context, id int64) error
	GetByTargetID(ctx context.Context, amount, id int64) (*[]models.Command, error)
	GetNextCommand(ctx context.Context, targetId int64) (*models.Command, error)
	Update(ctx context.Context, c *models.Command) error
	CreateNewCommand(ctx context.Context, c *models.Command) error
	CreateTable(ctx context.Context) error
	DropTable(ctx context.Context) error
}
