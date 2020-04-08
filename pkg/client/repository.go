package client

import (
	"context"
	"time"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Repository interface {
	GetByNameAndPassword(ctx context.Context, name, password string) (*models.Client, error)
	UpdateCSRFToken(ctx context.Context, id int64, csrfToken string, updatedAt time.Time) error
	CreateNewClient(ctx context.Context, c *models.Client) error
	CreateTable(ctx context.Context) error
	DropTable(ctx context.Context) error
}
