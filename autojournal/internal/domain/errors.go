package domain

import "errors"

var (
	ErrValidation     = errors.New("validation error")
	ErrNotFound       = errors.New("not found")
	ErrInfrastructure = errors.New("infrastructure error")
)
