package tusecase

import (
	"context"
	"fmt"
	"github.com/mojodojo101/c2server/pkg/command"
	"github.com/mojodojo101/c2server/pkg/models"
	"github.com/mojodojo101/c2server/pkg/target"
	"io/ioutil"
	"time"
)

type targetUsecase struct {
	targetRepo     target.Repository
	cmdRepo        command.Repository
	contextTimeout time.Duration
}

func NewTargetUsecase(tr target.Repository, cr command.Repository, timeout time.Duration) target.Usecase {
	return &targetUsecase{
		targetRepo:     tr,
		cmdRepo:        cr,
		contextTimeout: timeout,
	}
}

func (tu *targetUsecase) CreateTable(ctx context.Context) error {
	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()
	err := tu.targetRepo.CreateTable(cctx)

	if err != nil {
		return err
	}
	return nil

}
func (tu *targetUsecase) FetchCmdResponse(ctx context.Context, t *models.Target, cmdId int64) ([]byte, error) {
	//probably wanna change this to io.Pipe
	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	c, err := tu.cmdRepo.GetByID(cctx, cmdId)
	if err != nil {
		return nil, err
	}
	if c.Executed == false {
		return nil, models.ErrNotExecuted
	}
	cmdPath := fmt.Sprintf("%v/%v", t.Path, cmdId)
	fmt.Printf("command path = %v", cmdPath)
	data, err := ioutil.ReadFile(cmdPath)

	return data, nil
}

func (tu *targetUsecase) GetByID(ctx context.Context, id int64) (*models.Target, error) {

	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	b, err := tu.targetRepo.GetByID(cctx, id)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (tu *targetUsecase) Store(ctx context.Context, t *models.Target) (int64, error) {
	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()
	existingTarget, _ := tu.GetByID(cctx, t.Id)
	if existingTarget != nil {
		return int64(0), models.ErrDuplicate
	}
	id, err := tu.targetRepo.CreateNewTarget(cctx, t)
	if err != nil {
		return int64(0), err
	}
	return id, nil

}

func (tu *targetUsecase) Delete(ctx context.Context, t *models.Target) error {
	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()
	existingTarget, _ := tu.GetByID(cctx, t.Id)
	if existingTarget != nil {
		return models.ErrDuplicate
	}
	err := tu.targetRepo.DeleteByID(cctx, t.Id)
	if err != nil {
		return err
	}
	return nil

}
func (tu *targetUsecase) Update(ctx context.Context, t *models.Target) error {
	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()
	_, err := tu.GetByID(cctx, t.Id)
	if err != nil {
		return models.ErrItemNotFound
	}
	err = tu.targetRepo.Update(cctx, t)
	if err != nil {
		return err
	}
	return nil

}
