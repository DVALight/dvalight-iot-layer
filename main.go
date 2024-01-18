package main

import (
	"fmt"
	"net"
)

const DVA_DPORT = 1337

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:1337")
	fmt.Printf("%s %s\n", addr, err)
}
