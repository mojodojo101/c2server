package models

import "time"

type Target struct {
	Id int64 `json:"id"`

	Ipv4 string `json:"ipv4"`

	Ipv6 string `json:"ipv6"`

	HostName string `json:"hostName"`

	Path string `json:"path"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}
