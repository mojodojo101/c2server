package handlers

import (
	"github.com/mojodojo101/c2server/pkg/activebeacon"
	"github.com/mojodojo101/c2server/pkg/beacon"
	"github.com/mojodojo101/c2server/pkg/client"
	"github.com/mojodojo101/c2server/pkg/command"
	_ "github.com/mojodojo101/c2server/pkg/models"
	"github.com/mojodojo101/c2server/pkg/target"
	"io"
	"encoding/json"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ClientHandler struct {
	ActiveBeaconUsecase activebeacon.Usecase
	TargetUsecase       target.Usecase
	CommandUsecase      command.Usecase
	ClientUsecase       client.Usecase
	BeaconUsecase       beacon.Usecase
}

func NewClientHandler(ab activebeacon.Usecase, t target.Usecase, cmd command.Usecase, c client.Usecase, b beacon.Usecase) *ClientHandler {
	return &ClientHandler{
		ActiveBeaconUsecase: ab,
		TargetUsecase:       t,
		CommandUsecase:      cmd,
		ClientUsecase:       c,
		BeaconUsecase:       b,
	}
}

func (ch *ClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
}

func (ch *ClientHandler) addTarget() error {
	ch.
}
