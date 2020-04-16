package abusecase

import (
	"context"
	"fmt"
	"github.com/mojodojo101/c2server/pkg/activebeacon"
	"github.com/mojodojo101/c2server/pkg/models"
	"github.com/mojodojo101/c2server/pkg/target"
	"time"
)

type activebeaconUsecase struct {
	activebeaconRepo activebeacon.Repository
	targetUsecase    target.Usecase
	contextTimeout   time.Duration
}

func NewActiveBeaconUsecase(ar activebeacon.Repository, tu target.Usecase, timeout time.Duration) activebeacon.Usecase {
	return &activebeaconUsecase{
		activebeaconRepo: ar,
		targetUsecase:    tu,
		contextTimeout:   timeout,
	}
}

func (au *activebeaconUsecase) CreateTable(ctx context.Context) error {
	cctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()
	err := au.activebeaconRepo.CreateTable(cctx)
	return err

}

func (au *activebeaconUsecase) GetByID(ctx context.Context, id int64) (*models.ActiveBeacon, error) {

	cctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()

	b, err := au.activebeaconRepo.GetByID(cctx, id)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (au *activebeaconUsecase) GetTargetByIpv4(ctx context.Context, host string) (*models.Target, error) {

	cctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()

	t, err := au.targetUsecase.GetByIpv4(cctx, host)
	return t, err

}
func (au *activebeaconUsecase) SetCmdExecuted(ctx context.Context, a *models.ActiveBeacon, response []byte) error {
	cctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()

	t, err := au.targetUsecase.GetByID(cctx, a.TId)
	fmt.Printf("\na.Cmdid1=%v,\n t1 =%#v\n", a.CmdId, t)

	if err != nil {

		return err
	}
	err = au.targetUsecase.SetCmdExecuted(cctx, t, a.CmdId, response)
	fmt.Printf("\na.Cmdid2=%v,\n t2 =%#v\n", a.CmdId, t)
	return err

}
func (au *activebeaconUsecase) GetNextCmd(ctx context.Context, a *models.ActiveBeacon) error {

	cctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()
	t, err := au.targetUsecase.GetByID(cctx, a.TId)
	if err != nil {
		return err
	}
	cmd, err := au.targetUsecase.GetNextCmd(cctx, t)
	if err != nil {
		return err
	}
	a.Cmd = cmd.Cmd
	a.CmdId = cmd.Id
	a.UpdatedAt = time.Now()
	err = au.Update(cctx, a)
	return err

}
func (au *activebeaconUsecase) Register(ctx context.Context, a *models.ActiveBeacon) error {
	cctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()
	existingBeacon, _ := au.GetByID(cctx, a.Id)
	if existingBeacon != nil {
		return models.ErrDuplicate
	}
	err := au.activebeaconRepo.CreateNewBeacon(cctx, a)
	return err

}

func (au *activebeaconUsecase) Delete(ctx context.Context, a *models.ActiveBeacon) error {
	cctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()
	existingBeacon, _ := au.GetByID(cctx, a.Id)
	if existingBeacon != nil {
		return models.ErrDuplicate
	}
	err := au.activebeaconRepo.DeleteByID(cctx, a.Id)
	return err

}
func (au *activebeaconUsecase) Update(ctx context.Context, a *models.ActiveBeacon) error {
	cctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()
	_, err := au.GetByID(cctx, a.Id)
	if err != nil {
		return models.ErrItemNotFound
	}
	err = au.activebeaconRepo.Update(cctx, a)
	return err

}
