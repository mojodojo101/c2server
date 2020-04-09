package abusecase

import (
	"context"
	"github.com/mojodojo101/c2server/pkg/activebeacon"
	"github.com/mojodojo101/c2server/pkg/models"
	"time"
)

type activebeaconUsecase struct {
	activebeaconRepo activebeacon.Repository
	contextTimeout   time.Duration
}

func NewActiveBeaconUsecase(ar activebeacon.Repository, timeout time.Duration) activebeacon.Usecase {
	return &activebeaconUsecase{
		activebeaconRepo: ar,
		contextTimeout:   timeout,
	}
}

func (au *activebeaconUsecase) CreateTable(ctx context.Context) error {
	cctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()
	err := au.activebeaconRepo.CreateTable(cctx)

	if err != nil {
		return err
	}
	return nil

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

func (au *activebeaconUsecase) Register(ctx context.Context, a *models.ActiveBeacon) (int64, error) {
	cctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()
	existingBeacon, _ := au.GetByID(cctx, a.Id)
	if existingBeacon != nil {
		return int64(0), models.ErrDuplicate
	}
	id, err := au.activebeaconRepo.CreateNewBeacon(cctx, a)
	if err != nil {
		return id, err
	}
	return id, nil

}

func (au *activebeaconUsecase) Delete(ctx context.Context, a *models.ActiveBeacon) error {
	cctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()
	existingBeacon, _ := au.GetByID(cctx, a.Id)
	if existingBeacon != nil {
		return models.ErrDuplicate
	}
	err := au.activebeaconRepo.DeleteByID(cctx, a.Id)
	if err != nil {
		return err
	}
	return nil

}
func (au *activebeaconUsecase) Update(ctx context.Context, a *models.ActiveBeacon) error {
	cctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()
	_, err := au.GetByID(cctx, a.Id)
	if err != nil {
		return models.ErrItemNotFound
	}
	err = au.activebeaconRepo.Update(cctx, a)
	if err != nil {
		return err
	}
	return nil

}
