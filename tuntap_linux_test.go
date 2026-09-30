//go:build linux
// +build linux

package tuntap

import (
	"testing"
	"unsafe"
)

// The kernel copies sizeof(struct ifreq) bytes in both directions.
func TestIfReqSize(t *testing.T) {
	if got := unsafe.Sizeof(ifReq{}); got != 40 {
		t.Fatalf("sizeof(ifReq) = %d, want 40", got)
	}
}
