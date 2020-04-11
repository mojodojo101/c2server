package command

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Usecase interface {
	GetByID(ctx context.Context, id int64) (*models.Command, error)
	Delete(ctx context.Context, c *models.Command) error
	Update(ctx context.Context, c *models.Command) error
	Store(ctx context.Context, c *models.Command) error
	GetNextCommand(ctx context.Context, targetId int64) (*models.Command, error)
	ListCommandsByTargetID(ctx context.Context, targetId int64) ([]*models.Command, error)
	CreateTable(ctx context.Context) error
}
