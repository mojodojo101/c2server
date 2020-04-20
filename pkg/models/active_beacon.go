package models

import "time"

const (
	HTTPS uint32 = 1 << iota
	HTTP
	SMB
	DNS
	DNSHTTPS
	TCP
	NBNS
)

type ActiveBeacon struct {
	//Id of the active beacon
	Id int64 `json:"id"`

	//base beacon id points to an entry in the beacon table
	BId int64 `json:"bId"`

	//------gonna add this later-------
	//points to the parent beacon id in the active_beacon table
	//this id can be used to trace the edges/paths from a child/node beacon to the c2
	PId int64 `json:"pId"`

	//TargetId defines which target the beacon goes to for commands and where to store response information
	TId int64 `json:"tId"`

	//current command to execute
	CmdId int64 `json:"cmdId"`

	//Token should change between each interaction (Token generator algorithm)
	Token string `json:"token"`

	//Command to execute
	Cmd string `json:"command"`

	//C2-mode defines which protocol the beacon should use for communication to the server (HTTPS,HTTP,SMB...)
	C2m uint32 `json:"c2Mode"`

	//Peer mode defines which protocol the beacon should use for communication to a peer (HTTPS,HTTP,SMB...)
	Pm uint32 `json:"peerMode"`

	//Counter how many pings the beacon missed to send
	MissedPings int64 `json:"missedPings"`

	// duration of time in which a beacon should ping the server for new commands
	Ping int64 `json:"ping"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}
