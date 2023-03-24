package main

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:   ip,
		Port: port,
	}
	return server
}

func (s *Server) Handler(conn net.Conn) {
	// Handler
	fmt.Println("Link established successfully")
}

func (s *Server) Start() {
	// Socker listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err.Error())
	}
	// Close listen socket
	defer listener.Close()
	for {
		// Accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err.Error())
		}
		// Do handler
		go s.Handler(conn)
	}
}
