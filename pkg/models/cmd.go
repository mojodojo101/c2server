package models

import "time"

type Command struct {

	//unique id
	Id int64 `json:"id"`

	TId int64 `json:"target_id"`

	//path to beacon template
	Cmd string `json:"cmd"`

	Executed bool `json:"executed"`

	//make sure there are no 2 beacons excuting the same command
	Executing bool `json:"executing"`

	CreatedAt time.Time `json:"created_at"`

	ExecutedAt time.Time `json:"executed_at"`
}
