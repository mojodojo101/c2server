package cmdusecase

import (
	"context"
	"fmt"
	"github.com/mojodojo101/c2server/pkg/command"
	"github.com/mojodojo101/c2server/pkg/models"
	"time"
)

type commandUsecase struct {
	commandRepo    command.Repository
	contextTimeout time.Duration
}

func NewCommandUsecase(cr command.Repository, timeout time.Duration) command.Usecase {
	return &commandUsecase{
		commandRepo:    cr,
		contextTimeout: timeout,
	}
}

func (cu *commandUsecase) CreateTable(ctx context.Context) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	err := cu.commandRepo.CreateTable(cctx)

	if err != nil {
		return err
	}
	return nil

}
func (cu *commandUsecase) GetByID(ctx context.Context, id int64) (*models.Command, error) {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	b, err := cu.commandRepo.GetByID(cctx, id)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (cu *commandUsecase) Update(ctx context.Context, c *models.Command) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	_, err := cu.GetByID(cctx, c.Id)
	if err != nil {
		return models.ErrItemNotFound
	}
	err = cu.commandRepo.Update(cctx, c)
	if err != nil {
		return err
	}
	return nil

}
func (cu *commandUsecase) Store(ctx context.Context, c *models.Command) (int64, error) {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	existingCommand, _ := cu.GetByID(cctx, c.Id)
	if existingCommand != nil {
		return int64(0), models.ErrDuplicate
	}
	id, err := cu.commandRepo.CreateNewCommand(cctx, c)
	if err != nil {
		return int64(0), err
	}
	return id, nil

}
func (cu *commandUsecase) Delete(ctx context.Context, c *models.Command) error {
	fmt.Printf("c = %#v\n", c)
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	_, err := cu.GetByID(cctx, c.Id)
	if err != nil {
		return models.ErrItemNotFound
	}
	err = cu.commandRepo.DeleteByID(cctx, c.Id)
	if err != nil {
		return err
	}
	return nil

}
