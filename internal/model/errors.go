package model

import (
	"errors"
	"strings"
)

var (
	ErrNoRecord           = errors.New("no such item")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type ValidationError []string

func (e ValidationError) Error() string {
	return strings.Join([]string(e), "; ")
}

func (e ValidationError) String() string {
	return strings.Join([]string(e), "\n")
}
