package main

import (
	//"bytes"
	//"strings"
	"encoding/hex"
	"fmt"
	"strconv"
)

const PACKET_STX = "02"
const PACKET_ETX = "03"
const PACKET_ESC = "1b"

var PACKET_STX_Uint, _ = strconv.ParseUint(PACKET_STX, 16, 8)
var PACKET_ETX_Uint, _ = strconv.ParseUint(PACKET_ETX, 16, 8)
var PACKET_ESC_Uint, _ = strconv.ParseUint(PACKET_ESC, 16, 8)
var PACKET_STX_B = byte(PACKET_STX_Uint)
var PACKET_ETX_B = byte(PACKET_ETX_Uint)
var PACKET_ESC_B = byte(PACKET_ESC_Uint)

func CreatePacketDTU(frameType string, frameData string) []byte {
	var packet []byte
	// combine frameType and frameData for crc16
	s := frameType + frameData
	// calc crc
	var checksum uint16 = 24964
	s += strconv.FormatUint(uint64(checksum), 16)
	// escape
	s = escape(s)
	// add prefix and postfix
	s = PACKET_STX + s + PACKET_ETX
	fmt.Println(s)
	return packet
}

func escape(src string) string {
	slice, _ := hex.DecodeString(src)
	newSlice := []byte{}
	for i := 0; i < len(slice); i++ {
		switch slice[i] {
		case PACKET_STX_B:
			fallthrough
		case PACKET_ETX_B:
			fallthrough
		case PACKET_ESC_B:
			fmt.Printf("replace %x\n", slice[i])
			newSlice = append(newSlice, PACKET_ESC_B, slice[i]-PACKET_ESC_B)
		default:
			newSlice = append(newSlice, slice[i])
		}
	}
	return hex.EncodeToString(newSlice)
}

func main() {
	p := CreatePacketDTU("02", "34363030323038323234383534323000187141C3")
	fmt.Printf("%X\n", p)
}
