package model

import "errors"

var ErrNoRecord = errors.New("no such item")

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserAlreadyExists  = errors.New("user already exists")
)
