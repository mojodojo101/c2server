package tusecase

import (
	"context"
	"fmt"
	"github.com/mojodojo101/c2server/pkg/command"
	"github.com/mojodojo101/c2server/pkg/models"
	"github.com/mojodojo101/c2server/pkg/target"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type targetUsecase struct {
	targetRepo     target.Repository
	cmdUsecase     command.Usecase
	contextTimeout time.Duration
}

func NewTargetUsecase(tr target.Repository, cu command.Usecase, timeout time.Duration) target.Usecase {
	return &targetUsecase{
		targetRepo:     tr,
		cmdUsecase:     cu,
		contextTimeout: timeout,
	}
}

func (tu *targetUsecase) CreateTable(ctx context.Context) error {
	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()
	err := tu.targetRepo.CreateTable(cctx)
	return err

}
func (tu *targetUsecase) FetchCmdResponse(ctx context.Context, t *models.Target, cmdId int64) ([]byte, error) {
	//probably wanna change this to io.Pipe
	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	c, err := tu.cmdUsecase.GetByID(cctx, cmdId)
	if err != nil {
		return nil, err
	}
	if c.Executed == false {
		return nil, models.ErrNotExecuted
	}
	cmdPath := fmt.Sprintf("%v/%v", t.Path, cmdId)
	data, err := ioutil.ReadFile(cmdPath)

	return data, err
}

//GetNextCmd
//fetch cmd to execute from command table and set status to executing
func (tu *targetUsecase) GetNextCmd(ctx context.Context, t *models.Target) (int64, string, error) {
	//probably wanna change this to io.Pipe
	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	c, err := tu.cmdUsecase.GetNextCommand(cctx, t.Id)
	if err != nil {
		return 0, "", err
	}
	c.Executing = true
	err = tu.cmdUsecase.Update(cctx, c)
	return c.Id, c.Cmd, err
}

func (tu *targetUsecase) SetCmdExecuted(ctx context.Context, t *models.Target, cmdId int64, response []byte) error {
	//probably wanna change this to io.Pipe
	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	c, err := tu.cmdUsecase.GetByID(cctx, t.Id)
	if err != nil {
		return models.ErrItemNotFound
	}
	c.Executed = true
	c.ExecutedAt = time.Now()
	cmdPath := fmt.Sprintf("%v/%v", strings.TrimSpace(t.Path), cmdId)
	fmt.Printf("path =%v", cmdPath)
	err = ioutil.WriteFile(cmdPath, response, os.FileMode(0600))

	if err != nil {
		return err
	}

	err = tu.cmdUsecase.Update(cctx, c)
	return err
}

func (tu *targetUsecase) GetByID(ctx context.Context, id int64) (*models.Target, error) {

	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	b, err := tu.targetRepo.GetByID(cctx, id)

	return b, err

}
func (tu *targetUsecase) GetByIpv4(ctx context.Context, ipv4 string) (*models.Target, error) {

	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	t, err := tu.targetRepo.GetByIpv4(cctx, ipv4)

	return t, err

}
func (tu *targetUsecase) Store(ctx context.Context, t *models.Target) error {
	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()
	existingTarget, _ := tu.GetByID(cctx, t.Id)
	if existingTarget != nil {
		return models.ErrDuplicate
	}
	err := tu.targetRepo.CreateNewTarget(cctx, t)
	return err

}

func (tu *targetUsecase) Delete(ctx context.Context, t *models.Target) error {
	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()
	existingTarget, _ := tu.GetByID(cctx, t.Id)
	if existingTarget != nil {
		return models.ErrDuplicate
	}
	err := tu.targetRepo.DeleteByID(cctx, t.Id)

	return err

}
func (tu *targetUsecase) Update(ctx context.Context, t *models.Target) error {
	cctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()
	_, err := tu.GetByID(cctx, t.Id)
	if err != nil {
		return models.ErrItemNotFound
	}
	err = tu.targetRepo.Update(cctx, t)

	return err

}
