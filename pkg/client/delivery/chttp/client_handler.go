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

// this enum tells the server and the client what kind of request/response to expect
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

//subject to change
type ClientRequest struct {
	//Client ID
	Id int64 `json:"id"`
	//Target Id
	TId int64 `json:"tId"`
	//Target Ipv4
	TIpv4 string `json:"tIpv4"`
	//Target Hostname
	THostName string `json:"tHostName"`
	//Command id
	CmdId int64 `json:"cmdId"`
	//Command
	Cmd string `json:"cmd"`
	//Active Beacon ID
	AbId int64 `json:"abId"`
	//ACtive Beacon Ping duration
	AbPing int64 `json:"abPing"`
	//Beacon Id
	BId int64 `json:"bId"`
	//Beacon Operating System
	BOs string `json:"bOs"`
	//Beacon Language
	BLang string `json:"bLang"`
	//Name of the client
	Name string `json:"name"`
	//Password of the client
	Password string `json:"password"`
	//Token of the client
	Token string `json:"token"`
	//CSRF Token of the client
	CSRFToken string `json:"csrfToken"`
	//Request type check out the enum above for the request types
	RequestType int64 `json:"requestType"`
	//Amount is only for how much information the client asks for (for example 20 targets max)
	Amount int64 `json:"amount"`
}
type ClientResponse struct {
	//Server Token
	Token string `json:"token"`
	//Next CSRF
	CSRFToken string `json:"csrfToken"`
	//Again the same
	ResponseType int64 `json:"responseType"`
	//The amount of objects the server found for the request type
	ObjectSize int `json:"objectSize"`
	//If there was a Server error
	Error bool `json:"error"`
	//Json encoded Server Response
	Response string `json:"response"`
}

type ClientHandler struct {
	CUsecase client.Usecase
}

func NewHandler(cu client.Usecase) ClientHandler {
	return ClientHandler{
		CUsecase: cu,
	}

}

//Handles the Client Requests
// decode -> signin -> parse request for requesttype and change execution flow depending on it -> write response
//subject to change
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
	switch cReq.RequestType {

	case LIST_COMMANDS:
		cResp, err = ch.listCommandsOfTarget(ctx, c, cReq.TId, cReq.Amount)

	case LIST_ACTIVE_BEACONS:
		cResp, err = ch.listActiveBeacons(ctx, c, cReq.Amount)

		/*
			case LIST_BEACONS:
				cResp, err = ch.listBeacons(ctx, c, cReq.TIpv4)
		*/
	case LIST_TARGETS:
		cResp, err = ch.listTargets(ctx, c, cReq.Amount)
	case RETRIEVE_COMMAND_RESPONSE:
		cResp, err = ch.retrieveCommandResponse(ctx, c, cReq.CmdId)

	/*
		case RETRIEVE_BEACON:
			cResp, err = ch.retrieveBeacon(ctx, c, cReq.TId, cReq.CmdId)
	*/
	case ADD_NEW_COMMAND:
		cResp, err = ch.addNewCommand(ctx, c, cReq.TId, cReq.Cmd)

	case ADD_NEW_TARGET:
		cResp, err = ch.addNewTarget(ctx, c, cReq.TIpv4)

	case REMOVE_COMMAND:
		cResp, err = ch.removeCommand(ctx, c, cReq.TId, cReq.CmdId)

	case REMOVE_TARGET:
		cResp, err = ch.removeTarget(ctx, c, cReq.TId)

	case REMOVE_ACTIVE_BEACON:
		cResp, err = ch.removeActiveBeacon(ctx, c, cReq.AbId)

	case UPDATE_COMMAND:
		cResp, err = ch.updateCommand(ctx, c, cReq.CmdId, cReq.Cmd)

	case UPDATE_TARGET:
		cResp, err = ch.updateTarget(ctx, c, cReq.TId, cReq.TIpv4, cReq.THostName)

	case UPDATE_ACTIVE_BEACON:
		cResp, err = ch.updateActiveBeacon(ctx, c, cReq.AbId, cReq.AbPing)

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

//marshals the server response
func (ah *ClientHandler) encode(ctx context.Context, cResp *ClientResponse) (string, error) {
	jResp, err := json.Marshal(cResp)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return string(jResp), err

}

//unmarshals the clientrequest
func (ah *ClientHandler) decode(ctx context.Context, r *http.Request) (*ClientRequest, error) {
	cReq := ClientRequest{}
	err := json.NewDecoder(r.Body).Decode(&cReq)
	if err != nil {
		return nil, errors.New("Could not Unmarshal json, invalid beacon reponse")
	}
	return &cReq, err

}

//currently just checks if password and name are valid
//subject to change
func (ch *ClientHandler) signIn(ctx context.Context, cReq *ClientRequest) (*models.Client, error) {
	c, err := ch.CUsecase.SignIn(ctx, cReq.Name, cReq.Password)
	if err != nil {
		fmt.Printf("couldnt sign in\n")
	}
	return c, err
}

func (ch *ClientHandler) listCommandsOfTarget(ctx context.Context, c *models.Client, tId, amount int64) (*ClientResponse, error) {

	fmt.Printf("\nlist commands with c := %v m tId=%v, amount %v\n", c, tId, amount)
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
	cResp.ObjectSize = len(cmds)
	cResp.ResponseType = LIST_COMMANDS
	cResp.Response = string(jcmds)

	return &cResp, err
}
func (ch *ClientHandler) listTargets(ctx context.Context, c *models.Client, amount int64) (*ClientResponse, error) {
	targets, err := ch.CUsecase.ListTargets(ctx, c, amount)
	if err != nil {
		return nil, err
	}
	jtargets, err := json.Marshal(targets)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.ObjectSize = len(targets)
	cResp.ResponseType = LIST_TARGETS
	cResp.Response = string(jtargets)

	return &cResp, err
}

func (ch *ClientHandler) listActiveBeacons(ctx context.Context, c *models.Client, amount int64) (*ClientResponse, error) {
	abs, err := ch.CUsecase.ListActiveBeacons(ctx, c, amount)
	if err != nil {
		return nil, err
	}
	jabs, err := json.Marshal(abs)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.ObjectSize = len(abs)
	cResp.ResponseType = LIST_ACTIVE_BEACONS
	cResp.Response = string(jabs)

	return &cResp, err

}

//gets the commandresponse from the filesystem internal_resources/target/<tid>/<cmdid>
func (ch *ClientHandler) retrieveCommandResponse(ctx context.Context, c *models.Client, cmdId int64) (*ClientResponse, error) {
	data, err := ch.CUsecase.RetrieveCommandResponse(ctx, c, cmdId)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.Error = false
	cResp.ResponseType = RETRIEVE_COMMAND_RESPONSE
	cResp.ObjectSize = 1
	cResp.Response = string(data)
	return &cResp, err
}

//adds a new command with a targetid and commandstring
func (ch *ClientHandler) addNewCommand(ctx context.Context, c *models.Client, tId int64, cmdString string) (*ClientResponse, error) {
	cmd := models.Command{}
	cmd.Cmd = cmdString
	cmd.CreatedAt = time.Now()
	cmd.ExecutedAt = time.Time{}
	cmd.TId = tId
	cmd.Executed = false
	cmd.Executing = false
	err := ch.CUsecase.AddNewCommand(ctx, c, &cmd)
	if err != nil {
		return nil, err
	}
	jCommand, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.Error = false
	cResp.ObjectSize = 1
	cResp.Response = string(jCommand)
	cResp.ResponseType = ADD_NEW_COMMAND

	return &cResp, err
}

//adds a new target with an ipv4
//subject to change
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
	cResp.Error = false
	cResp.ResponseType = ADD_NEW_TARGET
	cResp.ObjectSize = 1
	cResp.Response = string(jTarget)

	return &cResp, err
}

func (ch *ClientHandler) removeCommand(ctx context.Context, c *models.Client, tId, cmdId int64) (*ClientResponse, error) {
	err := ch.CUsecase.RemoveCommand(ctx, c, tId, cmdId)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.Error = false
	cResp.ResponseType = REMOVE_COMMAND
	return &cResp, err
}

func (ch *ClientHandler) removeTarget(ctx context.Context, c *models.Client, tId int64) (*ClientResponse, error) {
	err := ch.CUsecase.RemoveTarget(ctx, c, tId)
	if err != nil {
		return nil, err
	}

	cResp := ClientResponse{}
	cResp.Error = false
	cResp.ResponseType = REMOVE_TARGET
	return &cResp, err
}
func (ch *ClientHandler) removeActiveBeacon(ctx context.Context, c *models.Client, abId int64) (*ClientResponse, error) {
	err := ch.CUsecase.RemoveActiveBeacon(ctx, c, abId)
	if err != nil {
		return nil, err
	}

	cResp := ClientResponse{}
	cResp.Error = false
	cResp.ResponseType = REMOVE_ACTIVE_BEACON
	return &cResp, err
}

//updates a command with what should be executed
func (ch *ClientHandler) updateCommand(ctx context.Context, c *models.Client, cmdId int64, cmdString string) (*ClientResponse, error) {
	cmd := models.Command{}
	cmd.Cmd = cmdString
	cmd.Id = cmdId
	err := ch.CUsecase.UpdateCommand(ctx, c, &cmd)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.Error = false
	cResp.ResponseType = UPDATE_COMMAND

	return &cResp, err
}

//updates the target with new hostname and ipv4
func (ch *ClientHandler) updateTarget(ctx context.Context, c *models.Client, tId int64, ipv4, hostName string) (*ClientResponse, error) {
	target := models.Target{}
	target.Ipv4 = ipv4
	target.HostName = hostName
	target.Id = tId
	err := ch.CUsecase.UpdateTarget(ctx, c, &target)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.Error = false
	cResp.ResponseType = UPDATE_TARGET

	return &cResp, err
}

func (ch *ClientHandler) updateActiveBeacon(ctx context.Context, c *models.Client, abId, abPing int64) (*ClientResponse, error) {
	ab := models.ActiveBeacon{}
	ab.Id = abId
	ab.Ping = abPing
	err := ch.CUsecase.UpdateActiveBeacon(ctx, c, &ab)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.Error = false
	cResp.ResponseType = UPDATE_ACTIVE_BEACON

	return &cResp, err
}
