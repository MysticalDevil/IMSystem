package main

import (
	"fmt"
	"io"
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
	user := NewUser(conn, s)
	user.Online()
	// receive messages from users
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// receive users messages
			msg := string(buf[:n-1])
			user.DoMessage(msg)
		}
	}()

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
