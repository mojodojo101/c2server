package models

import "time"

type Client struct {
	Id int64 `json:"id"`

	Ip string `json:"ip" validate:"required"`

	Name string `json:"name" validate:"required"`

	Password string `json:"password"`

	Token string `json:"token"`

	CSRFToken string `json:"csrfToken"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}
