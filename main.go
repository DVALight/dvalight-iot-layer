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
	magic    uint64
	deviceId uint32
}

func (req *DVARequest) parseDVARequest(buf []byte, len int) bool {
	if len < DVA_REQUEST_LEN {
		return false
	}

	req.magic = binary.LittleEndian.Uint64(buf[0:8])
	req.deviceId = binary.LittleEndian.Uint32(buf[8:12])

	return req.magic == DVA_MAGIC
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
		if req.parseDVARequest(buf[0:len], len) {
			devices[req.deviceId] = addr

			fmt.Printf("device %d request\n", req.deviceId)
		} else {
			fmt.Println("failed to parse DVA request")
		}
	}
}
