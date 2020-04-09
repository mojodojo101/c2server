package models

import "errors"

var (
	ErrInternalServerError = errors.New("Internal server error")
	ErrItemNotFound        = errors.New("Coudld not find requested item")
	ErrDuplicate           = errors.New("Item with this ID already exists")
	ErrBadParam            = errors.New("The supplied parameters arent valid")
	ErrNotExecuted         = errors.New("The command hasnt executed yet")
	ErrInvalidClient       = errors.New("The clients ID/Name and authentication does not match")
)
