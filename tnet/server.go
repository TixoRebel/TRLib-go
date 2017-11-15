package tnet

import (
	"net"
	"log"
)

type Server struct {
	network  string
	address  string
	handle   func(*Server, *MultiChannelStream)
	listener net.Listener
	mcs *MultiChannelStream
}

func NewServer(network string, address string, handle func(*Server, *MultiChannelStream)) (s *Server) {
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
		log.Fatal(err)
		return err
	}
	for {
		con, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			mcs := NewMultiChannelStream(con)
			go mcs.Start()
			s.mcs = mcs

			s.handle(s, mcs)
		}()
	}
	return nil
}

func (s *Server) Close() (e1 error, e2 error) {
	e1 = s.mcs.Close()
	e2 = s.listener.Close()
	return
}