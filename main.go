package main

import (
	"fmt"
	"net"
	"os"
)

const DVA_DPORT = 1337

func main() {
	anyAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", DVA_DPORT))
	conn, err := net.ListenUDP("udp", anyAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		var buf [512]byte
		n, addr, err := conn.ReadFromUDP(buf[0:])
		fmt.Printf("%d %s %s\n", n, addr, err)
	}
}
