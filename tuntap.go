package tuntap

import (
	"fmt"
	"os"
)

// Opts represents configuration for interface. The "name"
// can be empty.
type Opts struct {
	Name        string
	Mode        Mode
	PackageInfo bool
}

// Interface represents virtual network interface.
type Interface struct {
	name string
	mode Mode

	f *os.File
}

// Read receives a packet from the interface.
//
// Blocks until a packet is sent into the virtual
// interface. At that point, the content of the packet
// is copied into the provided buffer.
//
// Make sure the buffer is large enough. It is MTU of the
// interface (usually 1500, unless reconfigured) + 4 for
// the header in case that packet info is prepended, MTU + size
// of ethernet frame (38 bytes, unless VLan tags are enabled). If
// the buffer isn't large enough, the packet gets truncated.
//
func (i *Interface) Read(b []byte) (int, error) {
	return i.f.Read(b)
}

// Write sends a packet into the interface.
//
// Sends a packet through the interface. The buffer
// must be valid representation of a packet (with
// appropriate headers).
//
// It is up to the caller to provide only packets
// that fit MTU.
//
func (i *Interface) Write(b []byte) (int, error) {
	return i.f.Write(b)
}

// Close closes underlying file descriptor. When the program
// closes the file descriptor, the network device and all
// corresponding routes will disappear.
func (i *Interface) Close() error {
	return i.f.Close()
}

// Mode returns the mode of the adapter. It is
// always the same as the one passed to New.
func (i *Interface) Mode() Mode {
	return i.mode
}

// The mode in which open the virtual network adapter.
type Mode int

const (
	// Depending on the type of device chosen the userspace program has to read/write
	// IP packets (with Tun) or ethernet frames (with Tap). Which one is being used
	// depends on the flags given with the ioctl().
	_ Mode = iota
	Tun
	Tap
)

// New creates TUN/TAP interface.
func New(opts Opts) (*Interface, error) {
	switch opts.Mode {
	case Tun:
	case Tap:
	default:
		return nil, fmt.Errorf("invalid interface mode: %v", opts.Mode)
	}

	fd, err := tunTapSetup(setupOpts{
		name:        opts.Name,
		mode:        opts.Mode,
		packageInfo: opts.PackageInfo,
	})
	if err != nil {
		return nil, err
	}

	return &Interface{
		name: opts.Name,
		mode: opts.Mode,
		f:    os.NewFile(uintptr(fd), opts.Name),
	}, nil
}
