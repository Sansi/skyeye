package client

import (
	"fmt"
	"github.com/edwardtoday/skyeye/utils"
	"net"
	// "time"
)

const strStatus = "brightness_control,brightness_control,automatic,0\nmanual_brightness,manual_brightness,7,0\ncur_photocell,cur_photocell,398,0\nlast_reset_time,last_reset_time,\"2013-12-25 16:50:12\",0\nbox_error,box0,1,1\n"
const strEvent = "cur_photocell,cur_photocell,398,0\nlast_reset_time,last_reset_time,\"2013-12-25 16:50:12\",0\nbox_error,box0,1,1\n"
const strPlaylist = "This is a test playlist.\nNo content is available for preview.\n"

type SkyeyeClient struct {
	servAddr         string
	conn             net.Conn
	sendBuf, recvBuf []byte
	err              error
	quit             bool
}

func (c *SkyeyeClient) Init() {
	// c.servAddr = "202.11.11.162:9800"
	c.servAddr = "202.11.12.186:9912"
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

func (c *SkyeyeClient) Send() {
	_, c.err = c.conn.Write(c.sendBuf[:utils.PacketLen(c.sendBuf)])
	utils.CheckError(c.err)
	utils.PrintSendBuf(c.sendBuf)
}

func (c *SkyeyeClient) DTULogin(data string) {
	copy(c.sendBuf, utils.CreatePacketDTU("00", data))
	c.Send()
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
				c.quit = true // test login
			} else {
				fmt.Printf("failed\n")
			}
		case byte(0x01):
			fmt.Println("Server sends keepalive. Replying...")
			copy(c.sendBuf, utils.CreatePacketDTU("01", ""))
			c.Send()
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

func CreateClient(cid chan int, id string) {
	fmt.Println("CreateClient: ", id)
	c := SkyeyeClient{}
	c.Init()
	c.Connect()
	defer c.Close()
	dtuInfo := utils.CreateDTUInfo(id)
	c.DTULogin(dtuInfo)
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
			copy(c.sendBuf, utils.CreatePacketDTU("ff", utils.CreatePacketDevice("00", "00010000000000", strStatus)))
			c.Send()
		case 2:
			fmt.Printf("event\n")
			copy(c.sendBuf, utils.CreatePacketDTU("ff", utils.CreatePacketDevice("00", "00020000000000", strEvent)))
			c.Send()
		case 3:
			fmt.Printf("playlist\n")
			copy(c.sendBuf, utils.CreatePacketDTU("ff", utils.CreatePacketDevice("00", "00030000000000", strPlaylist)))
			c.Send()
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
