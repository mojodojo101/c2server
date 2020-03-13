package models

import "time"

type Client struct {
	Id   int64	`json:"id"`

	Ip  string	`json:"ip" validate:"required"`

	Name string	`json:"name" validate:"required"`

	Password string	`json:"password"`

	Token string	`json:"token"`

	CreatedAt time.Time	`json:"created_at"`

	UpdatedAt time.Time	`json:"updated_at"`
	
}


