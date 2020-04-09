package client

import (
	"context"
	"time"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Repository interface {
	GetByNameAndPassword(ctx context.Context, name, password string) (*models.Client, error)
	DeleteByID(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*models.Client, error)
	UpdateCSRFToken(ctx context.Context, id int64, csrfToken string, updatedAt time.Time) error
	Update(ctx context.Context, c *models.Client) error
	CreateNewClient(ctx context.Context, c *models.Client) (int64, error)
	CreateTable(ctx context.Context) error
	DropTable(ctx context.Context) error
}
