package chttp

import (
	"context"
	_ "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mojodojo101/c2server/pkg/client"
	"github.com/mojodojo101/c2server/pkg/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

const (
	LIST_COMMANDS = iota
	LIST_ACTIVE_BEACONS
	LIST_BEACONS
	LIST_TARGETS
	RETRIEVE_COMMAND_RESPONSE
	RETRIEVE_BEACON
	ADD_NEW_COMMAND
	ADD_NEW_TARGET
	REMOVE_COMMAND
	REMOVE_TARGET
	REMOVE_ACTIVE_BEACON
	UPDATE_COMMAND
	UPDATE_TARGET
	UPDATE_ACTIVE_BEACON
)

type ClientRequest struct {
	Id          int64  `json:"id"`
	TId         int64  `json:"tId"`
	TIpv4       string `json:"tIpv4"`
	CmdId       int64  `json:"cmdId"`
	Cmd         string `json:"cmd"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Token       string `json:"token"`
	CSRFToken   string `json:"csrfToken"`
	RequestType int64  `json:"requestType"`
	Amount      int64  `json:"amount"`
}
type ClientResponse struct {
	Token        string `json:"token"`
	CSRFToken    string `json:"csrfToken"`
	ResponseType int64  `json:"responseType"`
	ObjectSize   int    `json:"objectSize"`
	Error        bool   `json:"error"`
	Response     string `json:"response"`
}
type ResponseError struct {
	Message string `json:"message"`
}

type ClientHandler struct {
	CUsecase client.Usecase
}

func NewHandler(cu client.Usecase) ClientHandler {
	return ClientHandler{
		CUsecase: cu,
	}

}
func (ch *ClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	ctx := context.Background()
	cReq, err := ch.decode(ctx, r)
	if err != nil {
		logrus.Error(err)
		return
	}
	var cResp *ClientResponse
	c, err := ch.signIn(ctx, cReq)
	if err != nil {
		logrus.Error(err)
		return
	}
	strings.TrimRight(cReq.TIpv4, " ")
	fmt.Printf("\ncReq := %#v\n client =%#v\n", cReq, c)
	fmt.Printf("req type = %v\n", cReq.RequestType)
	switch cReq.RequestType {

	case LIST_COMMANDS:
		cResp, err = ch.listCommandsOfTarget(ctx, c, cReq.TId, cReq.Amount)
		/*
			case LIST_ACTIVE_BEACONS:
				cResp, err = ch.listActiveBeacons(ctx, c, cReq.TIpv4)

			case LIST_BEACONS:
				cResp, err = ch.listBeacons(ctx, c, cReq.TIpv4)
		*/
	case LIST_TARGETS:
		cResp, err = ch.listTargets(ctx, c, cReq.Amount)
	/*
		case RETRIEVE_COMMAND_RESPONSE:
			cResp, err = ch.retrieveCommandResponse(ctx, c, cReq.TId, cReq.CmdId)

		case RETRIEVE_BEACON:
			cResp, err = ch.retrieveBeacon(ctx, c, cReq.TId, cReq.CmdId)
	*/
	case ADD_NEW_COMMAND:
		cResp, err = ch.addNewCommand(ctx, c, cReq.TId, cReq.Cmd)

	case ADD_NEW_TARGET:
		cResp, err = ch.addNewTarget(ctx, c, cReq.TIpv4)
		/*
			case REMOVE_COMMAND:
				cResp, err = ch.removeCommand(ctx, c, cReq.TId, cReq.Cmd)

			case REMOVE_TARGET:
				cResp, err = ch.removeTarget(ctx, c, cReq.TIpv4)

			case REMOVE_ACTIVE_BEACON:
				cResp, err = ch.removeActiveBeacon(ctx, c, cReq.TIpv4)

			case UPDATE_COMMAND:
				cResp, err = ch.updateCommand(ctx, c, cReq.TId, cReq.Cmd)

			case UPDATE_TARGET:
				cResp, err = ch.updateTarget(ctx, c, cReq.TIpv4)

			case UPDATE_ACTIVE_BEACON:
				cResp, err = ch.updateActiveBeacon(ctx, c, cReq.TIpv4)

		*/
	}

	if err != nil {
		logrus.Error(err)
	}
	jResp, err := ch.encode(ctx, cResp)
	if err != nil {
		logrus.Error(err)
		return
	}
	fmt.Fprintf(w, jResp)

	return
}

func (ah *ClientHandler) encode(ctx context.Context, cResp *ClientResponse) (string, error) {
	jResp, err := json.Marshal(cResp)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return string(jResp), err

}
func (ah *ClientHandler) decode(ctx context.Context, r *http.Request) (*ClientRequest, error) {
	cReq := ClientRequest{}
	err := json.NewDecoder(r.Body).Decode(&cReq)
	if err != nil {
		return nil, errors.New("Could not Unmarshal json, invalid beacon reponse")
	}
	return &cReq, err

}

func (ch *ClientHandler) signIn(ctx context.Context, cReq *ClientRequest) (*models.Client, error) {
	c, err := ch.CUsecase.SignIn(ctx, cReq.Name, cReq.Password)
	if err != nil {
		fmt.Printf("couldnt sign in\n")
	}
	return c, err
}

func (ch *ClientHandler) listCommandsOfTarget(ctx context.Context, c *models.Client, tId, amount int64) (*ClientResponse, error) {

	cmds, err := ch.CUsecase.ListTargetCommands(ctx, c, tId, amount)
	if err != nil {
		return nil, err
	}
	fmt.Printf("cmds := %v", cmds)
	jcmds, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.Token = c.Token
	cResp.CSRFToken = c.CSRFToken
	cResp.ObjectSize = len(cmds)
	cResp.ResponseType = LIST_COMMANDS
	cResp.Response = string(jcmds)

	return &cResp, err
}
func (ch *ClientHandler) listTargets(ctx context.Context, c *models.Client, amount int64) (*ClientResponse, error) {
	fmt.Printf("got to target list with %#v\n", c)
	cmds, err := ch.CUsecase.ListTargets(ctx, c, amount)
	if err != nil {
		return nil, err
	}
	fmt.Printf("cmds := %v", cmds)
	jcmds, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.Token = c.Token
	cResp.CSRFToken = c.CSRFToken
	cResp.ObjectSize = len(cmds)
	cResp.ResponseType = LIST_TARGETS
	cResp.Response = string(jcmds)

	return &cResp, err
}

func (ch *ClientHandler) retrieveCommandResponse(ctx context.Context, c *models.Client, tId, cmdId int64) (*ClientResponse, error) {
	data, err := ch.CUsecase.RetrieveCommandResponse(ctx, c, tId, cmdId)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.CSRFToken = "2831921092109021093"
	cResp.Token = "2389293829128092109"
	cResp.Error = false
	cResp.ResponseType = RETRIEVE_COMMAND_RESPONSE
	cResp.ObjectSize = 1
	cResp.Response = string(data)
	return &cResp, err
}

func (ch *ClientHandler) addNewCommand(ctx context.Context, c *models.Client, tId int64, cmdString string) (*ClientResponse, error) {
	cmd := models.Command{}
	cmd.Cmd = cmdString
	cmd.CreatedAt = time.Now()
	cmd.ExecutedAt = time.Time{}
	cmd.TId = tId
	cmd.Executed = false
	cmd.Executing = false
	err := ch.CUsecase.AddNewCommad(ctx, c, &cmd)
	if err != nil {
		return nil, err
	}
	jCommand, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", jCommand)
	cResp := ClientResponse{}
	cResp.CSRFToken = "2831921092109021093"
	cResp.Token = "2389293829128092109"
	cResp.Error = false
	cResp.ObjectSize = 1
	cResp.Response = string(jCommand)
	cResp.ResponseType = ADD_NEW_COMMAND

	return &cResp, err
}

func (ch *ClientHandler) addNewTarget(ctx context.Context, c *models.Client, ipv4 string) (*ClientResponse, error) {
	t := models.Target{}
	t.Ipv4 = ipv4
	t.CreatedAt = time.Now()
	t.UpdatedAt = t.CreatedAt

	err := ch.CUsecase.AddNewTarget(ctx, c, &t)
	if err != nil {
		return nil, err
	}
	jTarget, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.CSRFToken = "2831921092109021093"
	cResp.Token = "2389293829128092109"
	cResp.Error = false
	cResp.ResponseType = ADD_NEW_TARGET
	cResp.ObjectSize = 1
	cResp.Response = string(jTarget)

	fmt.Printf("CRESP:%v", cResp)
	return &cResp, err
}
