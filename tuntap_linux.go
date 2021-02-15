// +build linux

package tuntap

import (
	"os"
	"syscall"
	"unsafe"
)

// Linux supports some standard ioctls to configure network devices.
// They can be used on any socket's file descriptor regardless of
// the family or type.  Most of them pass an ifReq structure.
//
// Based on ifreq struct - https://elixir.bootlin.com/linux/v4.9.164/source/include/uapi/linux/if.h#L226
// Kernel waiting for flags field for tun/tap configuration.
type ifReq struct {
	name  [syscall.IFNAMSIZ]byte
	flags uint16
}

func ioctl(fd, req, arg uintptr) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, req, arg)
	if err != 0 {
		return os.NewSyscallError("ioctl", err)
	}

	return nil
}

type setupOpts struct {
	name        string
	mode        Mode
	packageInfo bool
}

func tunTapSetup(opts setupOpts) (int, error) {
	fd, err := syscall.Open("/dev/net/tun", os.O_RDWR|syscall.O_NONBLOCK, 0);
	if err != nil {
		return 0, err
	}

	var flags uint16
	switch opts.mode {
	case Tun:
		flags |= syscall.IFF_TUN
	case Tap:
		flags |= syscall.IFF_TAP
	}

	if !opts.packageInfo {
		flags |= syscall.IFF_NO_PI
	}

	var r ifReq
	copy(r.name[:], opts.name)
	r.flags = flags

	if err := ioctl(uintptr(fd), syscall.TUNSETIFF, uintptr(unsafe.Pointer(&r))); err != nil {
		return 0, err
	}

	return fd, nil
}
