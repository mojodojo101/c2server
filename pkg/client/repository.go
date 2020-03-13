package client

import(
	"context"
	"time"

	"github.com/mojodojo101/c2server/pkg/models"
)

type Repository interface {
	GetByID(ctx context.Context,id int64)( *models.Client,error)
	CreateNewClient(ctx context.Context,ip,name,password,token string,createdAt,updatedAt time.Time)(*models.Client,error)
	CreateClientTable(ctx context.Context)(error)
}