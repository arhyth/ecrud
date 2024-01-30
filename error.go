package ecrud

import (
	"errors"
)

var (
	ErrServerError = errors.New("server error")
)

type ErrBadRequest struct {
	Fields []string `json:"fields"`
}

func (e ErrBadRequest) Error() string {
	return "missing/invalid params"
}

type ErrNotFound struct {
	ID int `json:"id"`
}

func (e ErrNotFound) Error() string {
	return "record not found"
}
