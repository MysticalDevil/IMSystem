package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	go user.ListenMessage()

	return user
}

// ListenMessage  is used to monitor the current User channel,
// and once there is a message, it will be sent to the peer client
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Listen Message err:", err.Error())
			return
		}
	}
}
