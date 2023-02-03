package app

import "github.com/deweppro/go-sdk/errors"

var (
	errDepRunning     = errors.New("dependencies is already running")
	errDepNotRunning  = errors.New("dependencies are not running yet")
	errServiceUnknown = errors.New("unknown service")
	errBadFileFormat  = errors.New("is not a supported file format")
)
