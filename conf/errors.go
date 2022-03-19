package conf

import "errors"

var (
	ErrInvalidPort = errors.New("invalid port: must be in range [0-65535]")
)
