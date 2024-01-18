package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

const DVA_DPORT = 1337

const DVA_MAGIC uint64 = 0x544847494C415644
const DVA_REQUEST_LEN = 8 + 4

type DVARequest struct {
	Magic    uint64
	DeviceID uint32
}

func (req *DVARequest) ParseDVARequest(buf []byte, len int) bool {
	if len < DVA_REQUEST_LEN {
		return false
	}

	req.Magic = binary.LittleEndian.Uint64(buf[0:8])
	req.DeviceID = binary.LittleEndian.Uint32(buf[8:12])

	return req.Magic == DVA_MAGIC
}

func main() {
	anyAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", DVA_DPORT))
	conn, err := net.ListenUDP("udp", anyAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	devices := make(map[uint32]*net.UDPAddr)
	var buf [512]byte

	for {
		len, addr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println(err)
			continue
		}

		var req DVARequest
		if req.ParseDVARequest(buf[0:len], len) {
			devices[req.DeviceID] = addr

			fmt.Printf("device %d request\n", req.DeviceID)
		} else {
			fmt.Println("failed to parse DVA request")
		}
	}
}
