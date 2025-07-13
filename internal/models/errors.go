package models

import (
	"errors"
)

var (
	ErrNoRows         = errors.New("models: no matching record found")
	ErrDuplicateEmail = errors.New("models: duplicate user email")
)
