//go:build !linux
// +build !linux

package tuntap

import (
	"errors"
	"runtime"
)

func tunTapSetup(opts setupOpts) (int, string, error) {
	return 0, "", errors.New("tuntap: not supported on " + runtime.GOOS)
}
