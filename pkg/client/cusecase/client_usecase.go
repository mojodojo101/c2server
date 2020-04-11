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

	return err

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

	return c, err

}

func (cu *clientUsecase) Store(ctx context.Context, c *models.Client) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	existingClient, _ := cu.clientRepo.GetByID(cctx, c.Id)
	if existingClient != nil {
		return models.ErrDuplicate
	}

	err := cu.clientRepo.CreateNewClient(cctx, c)

	return err

}

func (cu *clientUsecase) Delete(ctx context.Context, c *models.Client) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	err := cu.isValidClient(cctx, cancel, c)
	if err != nil {
		return err
	}
	err = cu.clientRepo.DeleteByID(cctx, c.Id)

	return err

}
func (cu *clientUsecase) Update(ctx context.Context, c *models.Client) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	err := cu.isValidClient(cctx, cancel, c)
	if err != nil {
		return err
	}

	err = cu.clientRepo.Update(cctx, c)

	return err

}
