package cusecase

import (
	"context"
	"github.com/mojodojo101/c2server/pkg/client"
	"github.com/mojodojo101/c2server/pkg/models"
	"time"
)

type clientUsecase struct {
	clientRepo     client.Repository
	contextTimeout time.Duration
}

func NewClientUsecase(cr client.Repository, timeout time.Duration) client.Usecase {
	return &clientUsecase{
		clientRepo:     cr,
		contextTimeout: timeout,
	}
}

func (cu *clientUsecase) CreateTable(ctx context.Context) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	err := cu.clientRepo.CreateTable(cctx)

	if err != nil {
		return err
	}
	return nil

}

func (cu *clientUsecase) isValidClient(ctx context.Context, cancel context.CancelFunc, c *models.Client) error {
	defer cancel()

	existingClient, err := cu.clientRepo.GetByID(ctx, c.Id)
	if err != nil {
		return err
	}

	//check if given client c has either a valid password or valid tokens
	if existingClient.Name == c.Name {
		if existingClient.Password == c.Password || (existingClient.Token == c.Token && existingClient.CSRFToken == c.Token) {
			return nil
		}
	}

	return models.ErrInvalidClient
}

func (cu *clientUsecase) SignIn(ctx context.Context, name, password string) (*models.Client, error) {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	c, err := cu.clientRepo.GetByNameAndPassword(cctx, name, password)
	if err != nil {
		return nil, err
	}
	return c, nil

}

func (cu *clientUsecase) Store(ctx context.Context, c *models.Client) (int64, error) {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	existingClient, _ := cu.clientRepo.GetByID(cctx, c.Id)
	if existingClient != nil {
		return int64(0), models.ErrDuplicate
	}

	id, err := cu.clientRepo.CreateNewClient(cctx, c)
	if err != nil {
		return int64(0), err
	}
	return id, nil

}

func (cu *clientUsecase) Delete(ctx context.Context, c *models.Client) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	err := cu.isValidClient(cctx, cancel, c)
	if err != nil {
		return err
	}
	return nil

}
func (cu *clientUsecase) Update(ctx context.Context, c *models.Client) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	_, err := cu.clientRepo.GetByID(cctx, c.Id)
	if err != nil {
		return models.ErrItemNotFound
	}
	err = cu.clientRepo.Update(cctx, c)
	if err != nil {
		return err
	}
	return nil

}
