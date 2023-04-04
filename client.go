package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		flag:       99,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err.Error())
		return nil
	}

	client.conn = conn

	return client
}

func (c *Client) menu() bool {
	var _flag int

	fmt.Println("1. Public chat")
	fmt.Println("2. Private chat")
	fmt.Println("3. Rename")
	fmt.Println("0. Exit")

	_, err := fmt.Scanln(&_flag)
	if err != nil {
		fmt.Println("fmt.Scanln err:", err.Error())
	}

	if _flag >= 0 && _flag <= 3 {
		c.flag = _flag
		return true
	}
	fmt.Println(">>>>>>>> Please enter legal range. <<<<<<<<")
	return false
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {

		}

		switch c.flag {
		case 1:
			fmt.Println("Public chat...")
			c.PublicChat()
			break
		case 2:
			fmt.Println("Private chat...")
			c.PrivateChat()
			break
		case 3:
			c.UpdateName()
			break
		}
	}
}

func (c *Client) DealResponse() {
	_, err := io.Copy(os.Stdout, c.conn)
	if err != nil {
		fmt.Println("io.Copy err:", err.Error())
	}
}

func (c *Client) UpdateName() bool {
	fmt.Println(">>>>>>>> Please input username:")
	_, err := fmt.Scanln(&c.Name)
	if err != nil {
		fmt.Println("fmt.Scanln err:", err.Error())
	}

	sendMsg := fmt.Sprintf("rename|%s\n", c.Name)
	_, err = c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err.Error())
		return false
	}
	return true
}

func (c *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>>>>>>> Please input message:")
	_, err := fmt.Scanln(&chatMsg)
	if err != nil {
		fmt.Println("fmt.Scanln err:", err.Error())
	}

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err = c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err.Error())
				break
			}
		}

		chatMsg = ""

		fmt.Println(">>>>>>>> Please input message:")
		_, err := fmt.Scanln(&chatMsg)
		if err != nil {
			fmt.Println("fmt.Scanln err:", err.Error())
		}
	}

}

func (c *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

func (c *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	c.SelectUsers()
	_, err := fmt.Scanln(&remoteName)
	if err != nil {
		fmt.Println("fmt.Scanln err:", err.Error())
	}

	for remoteName != "exit" {
		fmt.Println(">>>>>>>> Please input message")
		_, err := fmt.Scanln(&chatMsg)
		if err != nil {
			fmt.Println("fmt.Scanln err:", err.Error())
		}
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := fmt.Sprintf("to|%s|%s\n\n", remoteName, chatMsg)
				_, err = c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write err:", err.Error())
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>>>>> Please input message:")
			_, err := fmt.Scanln(&chatMsg)
			if err != nil {
				fmt.Println("fmt.Scanln err:", err.Error())
			}
		}

		remoteName = ""
		fmt.Println(">>>>>>>> Please enter username you want to chat with:")
		_, err = fmt.Scanln(&remoteName)
		if err != nil {
			fmt.Println("fmt.Scanln err:", err.Error())
		}
	}
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

	go client.DealResponse()

	fmt.Println(">>>>>>>> Connection success...")
	client.Run()
}
