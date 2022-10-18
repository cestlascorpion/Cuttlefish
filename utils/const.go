package utils

import (
	"errors"
	"time"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

const (
	defaultPrefix = "online"
	defaultExpire = time.Minute * 3
)
