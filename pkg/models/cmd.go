package models

import "time"

type Command struct {

	//unique id
	Id int64 `json:"id"`

	TId int64 `json:"target_id"`

	//path to beacon template
	Cmd string `json:"cmd"`

	Executed bool `json:"executed"`

	CreatedAt time.Time `json:"created_at"`

	ExecutedAt time.Time `json:"executed_at"`
}
