package tuntap

import "testing"

func TestNewInvalidName(t *testing.T) {
	for _, name := range []string{"0123456789abcdef", "tun\x000"} {
		if _, err := New(Opts{Name: name, Mode: Tun}); err == nil {
			t.Errorf("New(%q): expected error", name)
		}
	}
}
