# Go TUN/TAP support

Designed by https://docs.rs/tun-tap/

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
		Name: "tun0",
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

	packet := make([]byte, 1504)

	n, err := nic.Read(packet)
	if err != nil {
		log.Fatalf("failed to read packet: %v", err)
	}

	fmt.Println("Packet:", packet[:n])
}
```

```shell
$ go build
$ sudo setcap cap_net_admin=eip ./bin
$ ./bin
```

## Usages

github.com/rokkerruslan/tcp

## Known issues

- It is tested only on Linux and probably doesn't work anywhere else, even
  though other systems have some TUN/TAP support. Reports that it works (or not)
  and pull request to add other sustem's support are welcome.

Creating the devices requires CAP_NET_ADMIN privileges (most commonly done by
running as root).

```
$ cat main.go
// file...

$ go build
$ sudo setcap cap_net_admin=eip ./bin
$ ./bin
```

## Alternatives

1. water. Has incorrect flags for created interface and we can't get EtherType
   from packet. Maybe now suppored.
