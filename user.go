package main

import (
	"fmt"
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
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

func (u *User) Online() {
	// user is online, add user to OnlineMap
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// broadcast the current user online message
	u.server.Broadcast(u, "Online")
}

func (u *User) Offline() {
	// user is offline, remove user to OnlineMap
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	// broadcast the current user offline message
	u.server.Broadcast(u, "Offline")
}

func (u *User) DoMessage(msg string) {
	u.server.Broadcast(u, msg)
}
