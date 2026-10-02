# Go TUN/TAP support

Inspired by https://docs.rs/tun-tap/

## Examples

```go
package main

import (
	"fmt"
	"log"

	"github.com/rokkerruslan/tuntap"
)

func main() {
	nic, err := tuntap.New(tuntap.Opts{
		Name: "tun%d", // the kernel picks the first free number
		Mode: tuntap.Tun,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := nic.Close(); err != nil {
			log.Println(err)
		}
	}()

	fmt.Println("Interface:", nic.Name())

	// MTU (1500); add 4 bytes if PacketInfo is enabled.
	packet := make([]byte, 1500)

	n, err := nic.Read(packet)
	if err != nil {
		log.Fatalf("failed to read packet: %v", err)
	}

	fmt.Println("Packet:", packet[:n])
}
```

```shell
$ go build -o bin
$ sudo setcap cap_net_admin=eip ./bin
$ ./bin
```

## Usages

github.com/rokkerruslan/netstack

## Known issues

- It is tested only on Linux and probably doesn't work anywhere else, even
  though other systems have some TUN/TAP support. Reports that it works (or not)
  and pull request to add other system's support are welcome.

Creating the devices requires CAP_NET_ADMIN privileges (most commonly done by
running as root).

## Alternatives

1. water. Has incorrect flags for created interface and we can't get EtherType
   from packet. Maybe now supported.
