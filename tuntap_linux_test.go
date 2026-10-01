//go:build linux
// +build linux

package tuntap

import (
	"testing"
	"unsafe"
)

// The kernel copies sizeof(struct ifreq) bytes in both directions,
// which is at most 40 bytes (on 64-bit architectures).
func TestIfReqSize(t *testing.T) {
	if got := unsafe.Sizeof(ifReq{}); got < 40 {
		t.Fatalf("sizeof(ifReq) = %d, want at least 40", got)
	}
}
