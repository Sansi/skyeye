package main

import (
	// "bytes"
	"fmt"
	"github.com/edwardtoday/skyeye/utils"
	"net"
)

const status = "\x02\xFF\x1B\xE7\x00\x00\x01\x00\x00\x00\x00\x00\x62\x72\x69\x67\x68\x74\x6E\x65\x73\x73\x5F\x63\x6F\x6E\x74\x72\x6F\x6C\x2C\x62\x72\x69\x67\x68\x74\x6E\x65\x73\x73\x5F\x63\x6F\x6E\x74\x72\x6F\x6C\x2C\x61\x75\x74\x6F\x6D\x61\x74\x69\x63\x2C\x30\x0A\x6D\x61\x6E\x75\x61\x6C\x5F\x62\x72\x69\x67\x68\x74\x6E\x65\x73\x73\x2C\x6D\x61\x6E\x75\x61\x6C\x5F\x62\x72\x69\x67\x68\x74\x6E\x65\x73\x73\x2C\x37\x2C\x30\x0A\x63\x75\x72\x5F\x70\x68\x6F\x74\x6F\x63\x65\x6C\x6C\x2C\x63\x75\x72\x5F\x70\x68\x6F\x74\x6F\x63\x65\x6C\x6C\x2C\x33\x39\x38\x2C\x30\x0A\x63\x75\x72\x5F\x62\x72\x69\x67\x68\x74\x6E\x65\x73\x73\x2C\x63\x75\x72\x5F\x62\x72\x69\x67\x68\x74\x6E\x65\x73\x73\x2C\x31\x38\x2C\x30\x0A\x73\x69\x67\x6E\x5F\x6F\x66\x66\x2C\x73\x69\x67\x6E\x5F\x6F\x66\x66\x2C\x30\x2C\x30\x0A\x6C\x61\x73\x74\x5F\x72\x65\x73\x65\x74\x5F\x74\x69\x6D\x65\x2C\x6C\x61\x73\x74\x5F\x72\x65\x73\x65\x74\x5F\x74\x69\x6D\x65\x2C\x22\x32\x30\x31\x33\x2D\x31\x32\x2D\x32\x35\x20\x31\x36\x3A\x35\x30\x3A\x31\x32\x22\x2C\x30\x0A\x62\x6F\x78\x5F\x65\x72\x72\x6F\x72\x2C\x62\x6F\x78\x30\x2C\x31\x2C\x31\x0A\x62\x6F\x78\x5F\x65\x72\x72\x6F\x72\x2C\x62\x6F\x78\x31\x2C\x31\x2C\x31\x0A\x62\x6F\x78\x5F\x65\x72\x72\x6F\x72\x2C\x62\x6F\x78\x32\x2C\x31\x2C\x31\x0A\x7B\x1B\xE8\xC9\x09\x03"

type SkyeyeClient struct {
	servAddr         string
	conn             net.Conn
	sendBuf, recvBuf []byte
	err              error
	quit             bool
}

func (c *SkyeyeClient) Init() {
	c.servAddr = "202.11.11.162:9801"
	c.sendBuf = make([]byte, 2048)
	c.recvBuf = make([]byte, 2048)
	c.quit = false
}

func (c *SkyeyeClient) Connect() {
	c.conn, c.err = net.Dial("tcp", c.servAddr)
	utils.CheckError(c.err)
}

func (c *SkyeyeClient) Close() {
	c.conn.Close()
}

func (c *SkyeyeClient) DTULogin(data string) {
	copy(c.sendBuf, utils.CreatePacketDTU("00", data))
	_, c.err = c.conn.Write(c.sendBuf[:utils.PacketLen(c.sendBuf)])
	utils.CheckError(c.err)
	utils.PrintSendBuf(c.sendBuf)
}

func (c *SkyeyeClient) Loop() {
	for c.quit == false {
		_, c.err = c.conn.Read(c.recvBuf)
		utils.CheckError(c.err)
		utils.PrintRecvBuf(c.recvBuf)
		// decode
		packet := utils.Unescape(c.recvBuf[:utils.PacketLen(c.recvBuf)])
		utils.PrintPacket(packet)
		switch packet[1] {
		case byte(0x00):
			fmt.Printf("DTU login ")
			if packet[2] == byte(0x00) {
				fmt.Printf("successful\n")
			} else {
				fmt.Printf("failed\n")
			}
		case byte(0x01):
			fmt.Println("Server sends keepalive. Replying...")
		case byte(0x02):
			fmt.Println("Server queries DTU info.")
		case byte(0xff):
			fmt.Println("Server talks to device. Forwarding...")
			c.ProcDevicePacket(packet[2:utils.PacketLen(packet)])
		default:
			fmt.Printf("Invalid DTU packet type from server: %x", packet[1])
		}
	}

}

var clientNum int

func createClient(cid chan int) {
	c := SkyeyeClient{}
	c.Init()
	c.Connect()
	defer c.Close()
	c.DTULogin("34363030323038323234383534323000187141C3")
	c.Loop()

	clientNum++
	fmt.Println("Client", clientNum, "finished.")
	cid <- clientNum
}

func (c *SkyeyeClient) ProcDevicePacket(src []byte) {
	packet := utils.Unescape(src)
	t := packet[1]
	switch t {
	case byte(0x00):
		fmt.Printf("Server get ")
		obj := utils.BytesToUint16(packet[2:4])
		switch obj {
		case 1:
			fmt.Printf("status\n")
			_, c.err = c.conn.Write([]byte(status))
			utils.CheckError(c.err)
		case 2:
			fmt.Printf("event\n")
			c.quit = true // game over
		case 3:
			fmt.Printf("playlist\n")
		default:
			fmt.Printf("unknown object: %x\n", obj)
		}
	case byte(0x01):
		fmt.Println("Server set function is not implemented yet.")
	case byte(0x02):
		fmt.Println("Do-not-reply function is not implemented yet.")
	default:
		fmt.Printf("Invalid device packet type from server: %x\n", t)
	}
}

func main() {
	cid := make(chan int)
	numBenchmark := 1
	for i := 0; i < numBenchmark; i++ {
		go createClient(cid)
	}
	for i := 0; i < numBenchmark; i++ {
		<-cid
	}

}
