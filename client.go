package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err.Error())
		return nil
	}

	client.conn = conn

	return client
}

var serverIP string
var serverPort int

// client -ip 127.0.0.1
func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "Set server ip address")
	flag.IntVar(&serverPort, "port", 8888, "Set server port")
}

func main() {
	flag.Parse()

	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>> Connection failure...")
		return
	}

	fmt.Println(">>>>>>>> Connection success...")
	select {}
}
