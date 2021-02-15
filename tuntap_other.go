// +build !linux

package tuntap

import (
	"errors"
	"os"
	"runtime"
)

func tunTapSetup() (*os.File, error) {
	return nil, errors.New("not supported on" + runtime.GOOS)
}

