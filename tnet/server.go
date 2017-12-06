package tnet

import (
	"net"
)

type Server struct {
	network string
	address string
	handle func(*Server, net.Conn)
	listener net.Listener
}

func NewServer(network string, address string, handle func(*Server, net.Conn)) (s *Server) {
	s = new(Server)
	s.network = network
	s.address = address
	s.handle = handle
	return
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	for {
		con, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			s.handle(s, con)
		}()
	}
	return nil
}

func (s *Server) Close() (e1 error, e2 error) {
	e2 = s.listener.Close()
	return
}