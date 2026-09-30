//go:build linux
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
//
// The kernel copies the whole sizeof(struct ifreq) (40 bytes) in both
// directions, so the structure is padded to the full size of the union.
type ifReq struct {
	name  [syscall.IFNAMSIZ]byte
	flags uint16
	_     [22]byte
}

func ioctl(fd, req, arg uintptr) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, req, arg)
	if err != 0 {
		return os.NewSyscallError("ioctl", err)
	}

	return nil
}

// tunTapSetup returns file descriptor of the configured device and
// the interface name assigned by the kernel.
func tunTapSetup(opts setupOpts) (int, string, error) {
	fd, err := syscall.Open("/dev/net/tun", os.O_RDWR|syscall.O_NONBLOCK|syscall.O_CLOEXEC, 0)
	if err != nil {
		return 0, "", os.NewSyscallError("open", err)
	}

	var flags uint16
	switch opts.mode {
	case Tun:
		flags |= syscall.IFF_TUN
	case Tap:
		flags |= syscall.IFF_TAP
	}

	if !opts.packetInfo {
		flags |= syscall.IFF_NO_PI
	}

	var r ifReq
	copy(r.name[:syscall.IFNAMSIZ-1], opts.name)
	r.flags = flags

	if err := ioctl(uintptr(fd), syscall.TUNSETIFF, uintptr(unsafe.Pointer(&r))); err != nil {
		syscall.Close(fd)
		return 0, "", err
	}

	name := r.name[:]
	for i, c := range name {
		if c == 0 {
			name = name[:i]
			break
		}
	}

	return fd, string(name), nil
}
