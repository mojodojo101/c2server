package client

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Usecase interface {
	ListTargetCommands(ctx context.Context, c *models.Client, targetId int64) ([]*models.Command, error)
	AddNewCommad(ctx context.Context, c *models.Client, tId int64, cmd string) error
	AddNewTarget(ctx context.Context, c *models.Client, ipv4 string) error
	RetrieveCommandResponse(ctx context.Context, c *models.Client, tId, cmdId int64) (string, error)
	SignIn(ctx context.Context, name, password string) (*models.Client, error)
	Delete(ctx context.Context, c *models.Client) error
	Update(ctx context.Context, c *models.Client) error
	Store(ctx context.Context, c *models.Client) error
	CreateTable(ctx context.Context) error
}
