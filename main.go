package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

const DVA_DPORT = 1337

const DVA_MAGIC uint64 = 0x544847494C415644
const DVA_REQUEST_LEN = 8 + 4
const DVA_RESPONSE_LEN = 8 + 4 + 1

const DVA_ENDPOINT = "http://localhost:4444/device"

var LE = binary.LittleEndian

type DVARequest struct {
	Magic    uint64
	DeviceID uint32
}

func (req *DVARequest) ParseDVARequest(buf []byte, len int) bool {
	if len < DVA_REQUEST_LEN {
		return false
	}

	req.Magic = LE.Uint64(buf[0:8])
	req.DeviceID = LE.Uint32(buf[8:12])

	return req.Magic == DVA_MAGIC
}

type DVAResponse struct {
	Magic    uint64
	DeviceID uint32 `json:"id"`
	State    bool   `json:"state"`
}

func BoolToByte(val bool) byte {
	if val {
		return 1
	} else {
		return 0
	}
}

func (res *DVAResponse) MakeDVAResponse() []byte {
	res.Magic = DVA_MAGIC

	buf := make([]byte, DVA_RESPONSE_LEN)
	LE.PutUint64(buf[0:8], res.Magic)
	LE.PutUint32(buf[8:12], res.DeviceID)
	buf[12] = BoolToByte(res.State)

	return buf
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

			httpRes, err := http.Get(fmt.Sprintf("%s/%d", DVA_ENDPOINT, req.DeviceID))
			if err != nil {
				fmt.Printf("Error %s when requesting back-end", err)
				continue
			}

			httpData, err := io.ReadAll(httpRes.Body)
			httpRes.Body.Close()

			if err != nil {
				fmt.Printf("http body error: %s\n", err)
				continue
			}

			var res DVAResponse
			err = json.Unmarshal(httpData, &res)
			if err != nil {
				fmt.Printf("json error: %s\n", err)
			}

			_, err = conn.WriteToUDP(res.MakeDVAResponse(), devices[res.DeviceID])
			if err == nil {
				fmt.Printf("sent response to %d device\n", res.DeviceID)
			} else {
				fmt.Printf("WriteToUDP error: %s\n", err)
			}
		} else {
			fmt.Println("failed to parse DVA request")
		}
	}
}
