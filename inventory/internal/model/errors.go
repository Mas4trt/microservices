package model

import "errors"

var (
	ErrPartNotFound       = errors.New("part not found")
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrAlreadyInitialized = errors.New("repository is already initialized")
	ErrEmptyData          = errors.New("cannot initialize repository with empty data")
)
