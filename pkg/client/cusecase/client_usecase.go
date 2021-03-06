package cusecase

import (
	"context"
	"fmt"
	"github.com/mojodojo101/c2server/config"
	"github.com/mojodojo101/c2server/pkg/activebeacon"
	"github.com/mojodojo101/c2server/pkg/client"
	"github.com/mojodojo101/c2server/pkg/models"
	"github.com/mojodojo101/c2server/pkg/target"
	"os"
	"time"
)

//DISCLAIMER
//this might be a very bad way to do the core domain logic, this is my first real go project and i am still trying to figure stuff out
// it might be better to just include the repositorys instead of going through multiple layers of usecases
//pros of doing it the way i implemented it:
//you can resuse some function calls and dont have to rewrite a bunch of code
//cons:
//complicated and no seperation of core logic for this module
type clientUsecase struct {
	clientRepo          client.Repository
	targetUsecase       target.Usecase
	activeBeaconUsecase activebeacon.Usecase
	contextTimeout      time.Duration
}

func NewClientUsecase(cr client.Repository, tu target.Usecase, au activebeacon.Usecase, timeout time.Duration) client.Usecase {
	return &clientUsecase{
		clientRepo:          cr,
		targetUsecase:       tu,
		activeBeaconUsecase: au,
		contextTimeout:      timeout,
	}
}

func (cu *clientUsecase) CreateTable(ctx context.Context) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	err := cu.clientRepo.CreateTable(cctx)

	return err

}
func (cu *clientUsecase) AddNewCommand(ctx context.Context, c *models.Client, cmd *models.Command) error {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	t, err := cu.targetUsecase.GetByID(cctx, cmd.TId)
	if err != nil {
		return err
	}
	err = cu.targetUsecase.StoreCmd(cctx, t, cmd)
	return err

}
func (cu *clientUsecase) RetrieveCommandResponse(ctx context.Context, c *models.Client, cmdId int64) (string, error) {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	data, err := cu.targetUsecase.FetchCmdResponse(cctx, cmdId)

	return string(data), err

}
func (cu *clientUsecase) ListActiveBeacons(ctx context.Context, c *models.Client, amount int64) ([]models.ActiveBeacon, error) {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	err := cu.isValidClient(cctx, c)
	if err != nil {
		return nil, err
	}
	ab, err := cu.activeBeaconUsecase.ListActiveBeacons(cctx, amount)
	return ab, err

}
func (cu *clientUsecase) ListTargets(ctx context.Context, c *models.Client, amount int64) ([]models.Target, error) {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	err := cu.isValidClient(cctx, c)
	if err != nil {
		return nil, err
	}
	targets, err := cu.targetUsecase.ListTargets(cctx, amount)
	return targets, err

}
func (cu *clientUsecase) ListTargetCommands(ctx context.Context, c *models.Client, tId, amount int64) ([]models.Command, error) {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	err := cu.isValidClient(cctx, c)
	if err != nil {
		return nil, err
	}

	t, err := cu.targetUsecase.GetByID(cctx, tId)

	if err != nil {
		return nil, err
	}
	cmds, err := cu.targetUsecase.ListCommands(cctx, t, amount)
	return cmds, err

}
func (cu *clientUsecase) UpdateCommand(ctx context.Context, c *models.Client, cmd *models.Command) error {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	err := cu.isValidClient(cctx, c)
	if err != nil {
		return err
	}
	t, err := cu.targetUsecase.GetByID(cctx, cmd.TId)
	if err != nil {
		return err
	}
	command, err := cu.targetUsecase.GetCommandById(cctx, t, cmd.Id)

	if err != nil {
		return err
	}
	command.Cmd = cmd.Cmd
	err = cu.targetUsecase.UpdateCommand(cctx, t, command)
	return err

}
func (cu *clientUsecase) UpdateActiveBeacon(ctx context.Context, c *models.Client, ab *models.ActiveBeacon) error {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	err := cu.isValidClient(cctx, c)
	if err != nil {
		return err
	}

	activeb, err := cu.activeBeaconUsecase.GetByID(cctx, ab.Id)

	if err != nil {
		return err
	}
	activeb.Ping = ab.Ping
	err = cu.activeBeaconUsecase.Update(cctx, activeb)
	return err

}
func (cu *clientUsecase) UpdateTarget(ctx context.Context, c *models.Client, t *models.Target) error {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	err := cu.isValidClient(cctx, c)
	if err != nil {
		return err
	}

	target, err := cu.targetUsecase.GetByID(cctx, t.Id)

	if err != nil {
		return err
	}
	target.HostName = t.HostName
	target.Ipv4 = t.Ipv4
	target.Ipv6 = t.Ipv6
	target.UpdatedAt = time.Now()
	err = cu.targetUsecase.Update(cctx, target)
	return err

}
func (cu *clientUsecase) RemoveCommand(ctx context.Context, c *models.Client, tId, cmdId int64) error {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	err := cu.isValidClient(cctx, c)
	if err != nil {
		return err
	}

	t, err := cu.targetUsecase.GetByID(cctx, tId)

	if err != nil {
		return err
	}
	err = cu.targetUsecase.RemoveCommand(cctx, t, cmdId)
	return err

}

func (cu *clientUsecase) RemoveTarget(ctx context.Context, c *models.Client, tId int64) error {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	err := cu.isValidClient(cctx, c)
	if err != nil {
		return err
	}

	t, err := cu.targetUsecase.GetByID(cctx, tId)

	if err != nil {
		return err
	}
	err = cu.targetUsecase.Delete(cctx, t)
	return err

}

func (cu *clientUsecase) RemoveActiveBeacon(ctx context.Context, c *models.Client, abId int64) error {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	err := cu.isValidClient(cctx, c)
	if err != nil {
		return err
	}

	ab, err := cu.activeBeaconUsecase.GetByID(cctx, abId)

	if err != nil {
		return err
	}
	err = cu.activeBeaconUsecase.Delete(cctx, ab)
	return err

}
func (cu *clientUsecase) AddNewTarget(ctx context.Context, c *models.Client, t *models.Target) error {

	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	err := cu.isValidClient(cctx, c)
	if err != nil {
		return err
	}

	err = cu.targetUsecase.Store(cctx, t)
	if err != nil {
		return err
	}

	//I am adding the targetpath afterwards because it depends on the id of the target we get form the db
	//+ i dont want the user to supply this
	TARGETPATH := config.GetResourceConfig().TargetsPath
	t.Path = fmt.Sprintf("%v%v", TARGETPATH, t.Id)

	err = os.Mkdir(t.Path, 0600)
	if err != nil {
		return err
	}
	err = cu.targetUsecase.Update(cctx, t)
	return err

}
func (cu *clientUsecase) isValidClient(ctx context.Context, c *models.Client) error {

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
	err := cu.isValidClient(cctx, c)
	if err != nil {
		return err
	}
	err = cu.clientRepo.DeleteByID(cctx, c.Id)

	return err

}
func (cu *clientUsecase) Update(ctx context.Context, c *models.Client) error {
	cctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()
	err := cu.isValidClient(cctx, c)
	if err != nil {
		return err
	}

	err = cu.clientRepo.Update(cctx, c)

	return err

}
