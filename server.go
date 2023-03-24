package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP   string
	Port int

	// Online user list
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// channel for message broadcasting
	Message chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// ListenMessenger is used to listen the goroutine of the Message broadcast channel,
// and once there is a message, it will be sent to all online users
func (s *Server) ListenMessenger() {
	for {
		msg := <-s.Message

		s.mapLock.Lock()

		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}

		s.mapLock.Unlock()
	}
}

func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]:%s:%s", user.Addr, user.Name, msg)

	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	user := NewUser(conn)
	// user is online, add user to OnlineMap
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	// broadcast the current user online message
	s.Broadcast(user, "Online")

	// current handler is blocked
	select {}
}

func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err.Error())
		return
	}
	// close listen socket
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("listener.Close() err:", err.Error())
		}
	}(listener)

	// start the goroutine listening to the messenger
	go s.ListenMessenger()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err.Error())
			continue
		}
		// do handler
		go s.Handler(conn)
	}
}
