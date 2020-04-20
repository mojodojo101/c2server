package client

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Usecase interface {
	ListTargetCommands(ctx context.Context, c *models.Client, targetId, amount int64) ([]models.Command, error)
	ListTargets(ctx context.Context, c *models.Client, amount int64) ([]models.Target, error)
	ListActiveBeacons(ctx context.Context, c *models.Client, amount int64) ([]models.ActiveBeacon, error)
	AddNewCommand(ctx context.Context, c *models.Client, cmd *models.Command) error
	RemoveCommand(ctx context.Context, c *models.Client, tId, cmdId int64) error
	RemoveTarget(ctx context.Context, c *models.Client, tId int64) error
	RemoveActiveBeacon(ctx context.Context, c *models.Client, tId int64) error
	UpdateCommand(ctx context.Context, c *models.Client, cmd *models.Command) error
	UpdateTarget(ctx context.Context, c *models.Client, t *models.Target) error
	UpdateActiveBeacon(ctx context.Context, c *models.Client, ab *models.ActiveBeacon) error
	AddNewTarget(ctx context.Context, c *models.Client, target *models.Target) error
	RetrieveCommandResponse(ctx context.Context, c *models.Client, cmdId int64) (string, error)
	SignIn(ctx context.Context, name, password string) (*models.Client, error)
	Delete(ctx context.Context, c *models.Client) error
	Update(ctx context.Context, c *models.Client) error
	Store(ctx context.Context, c *models.Client) error
	CreateTable(ctx context.Context) error
}
