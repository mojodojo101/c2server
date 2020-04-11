package models

import "time"

type Target struct {
	Id int64 `json:"id"`

	Ipv4 string `json:"ipv4"`

	Ipv6 string `json:"ipv6"`

	HostName string `json:"host_name"`

	Path string `json:"path"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
