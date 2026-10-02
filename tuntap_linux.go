//go:build linux
// +build linux

package tuntap

import (
	"os"

	"golang.org/x/sys/unix"
)

// tunTapSetup returns file descriptor of the configured device and
// the interface name assigned by the kernel.
func tunTapSetup(opts setupOpts) (int, string, error) {
	ifr, err := unix.NewIfreq(opts.name)
	if err != nil {
		return 0, "", err
	}

	fd, err := unix.Open("/dev/net/tun", unix.O_RDWR|unix.O_NONBLOCK|unix.O_CLOEXEC, 0)
	if err != nil {
		return 0, "", os.NewSyscallError("open", err)
	}

	var flags uint16
	switch opts.mode {
	case Tun:
		flags |= unix.IFF_TUN
	case Tap:
		flags |= unix.IFF_TAP
	}

	if !opts.packetInfo {
		flags |= unix.IFF_NO_PI
	}

	ifr.SetUint16(flags)

	if err := unix.IoctlIfreq(fd, unix.TUNSETIFF, ifr); err != nil {
		unix.Close(fd)
		return 0, "", os.NewSyscallError("ioctl", err)
	}

	return fd, ifr.Name(), nil
}
