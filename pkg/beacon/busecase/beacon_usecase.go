package busecase

import (
	"context"
	"github.com/mojodojo101/c2server/pkg/beacon"
	"github.com/mojodojo101/c2server/pkg/models"
	"io/ioutil"
	"time"
)

type beaconUsecase struct {
	beaconRepo     beacon.Repository
	contextTimeout time.Duration
}

func NewBeaconUsecase(br beacon.Repository, timeout time.Duration) beacon.Usecase {
	return &beaconUsecase{
		beaconRepo:     br,
		contextTimeout: timeout,
	}
}

func (bu *beaconUsecase) CreateTable(ctx context.Context) error {
	cctx, cancel := context.WithTimeout(ctx, bu.contextTimeout)
	defer cancel()
	err := bu.beaconRepo.CreateTable(cctx)

	if err != nil {
		return err
	}
	return nil

}
func (bu *beaconUsecase) RetrieveBeaconBuffer(ctx context.Context, b *models.Beacon) ([]byte, error) {
	//probably wanna change this to io.Pipe
	data, err := ioutil.ReadFile(b.Path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (bu *beaconUsecase) GetByID(ctx context.Context, id int64) (*models.Beacon, error) {

	cctx, cancel := context.WithTimeout(ctx, bu.contextTimeout)
	defer cancel()

	b, err := bu.beaconRepo.GetByID(cctx, id)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (bu *beaconUsecase) Store(ctx context.Context, b *models.Beacon) error {
	cctx, cancel := context.WithTimeout(ctx, bu.contextTimeout)
	defer cancel()
	existingBeacon, _ := bu.GetByID(cctx, b.Id)
	if existingBeacon != nil {
		return models.ErrDuplicate
	}
	err := bu.beaconRepo.CreateNewBeacon(cctx, b)

	return err

}

func (bu *beaconUsecase) Delete(ctx context.Context, b *models.Beacon) error {
	cctx, cancel := context.WithTimeout(ctx, bu.contextTimeout)
	defer cancel()
	existingBeacon, _ := bu.GetByID(cctx, b.Id)
	if existingBeacon != nil {
		return models.ErrDuplicate
	}
	err := bu.beaconRepo.DeleteByID(cctx, b.Id)
	if err != nil {
		return err
	}
	return nil

}
