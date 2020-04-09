package client

import (
	"context"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Usecase interface {
	SignIn(ctx context.Context, name, password string) (*models.Client, error)
	Delete(ctx context.Context, c *models.Client) error
	Update(ctx context.Context, c *models.Client) error
	Store(ctx context.Context, c *models.Client) (int64, error)
	CreateTable(ctx context.Context) error
}
