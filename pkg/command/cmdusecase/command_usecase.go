package cmdusecase

import (
	"context"
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

	return err

}
func (cu *commandUsecase) GetByID(ctx context.Context, id int64) (*models.Command, error) {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	cmd, err := cu.commandRepo.GetByID(cctx, id)
	return cmd, err

}
func (cu *commandUsecase) ListCommandsByTargetID(ctx context.Context, tId, amount int64) ([]models.Command, error) {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	cmds, err := cu.commandRepo.GetByTargetID(cctx, amount, tId)
	return cmds, err

}
func (cu *commandUsecase) GetNextCommand(ctx context.Context, targetId int64) (*models.Command, error) {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	cmd, err := cu.commandRepo.GetNextCommand(cctx, targetId)
	return cmd, err

}
func (cu *commandUsecase) Update(ctx context.Context, c *models.Command) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	_, err := cu.GetByID(cctx, c.Id)
	if err != nil {
		return models.ErrItemNotFound
	}
	err = cu.commandRepo.Update(cctx, c)

	return err

}
func (cu *commandUsecase) Store(ctx context.Context, c *models.Command) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	existingCommand, _ := cu.GetByID(cctx, c.Id)
	if existingCommand != nil {
		return models.ErrDuplicate
	}
	err := cu.commandRepo.CreateNewCommand(cctx, c)
	return err

}
func (cu *commandUsecase) Delete(ctx context.Context, c *models.Command) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	_, err := cu.GetByID(cctx, c.Id)
	if err != nil {
		return models.ErrItemNotFound
	}
	err = cu.commandRepo.DeleteByID(cctx, c.Id)

	return err

}
