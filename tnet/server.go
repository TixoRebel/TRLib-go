package tnet

import (
	"net"
	"log"
)

type Server struct {
	network string
	address string
	handle func(p *MultiChannelStream)
}

func NewServer(network string, address string, handle func(p *MultiChannelStream)) (s *Server) {
	s = new(Server)
	s.network = network
	s.address = address
	s.handle = handle
	return
}

func (s *Server) Start() error {
	ln, err := net.Listen(s.network, s.address)
	if err != nil {
		log.Fatal(err)
		return err
	}
	for {
		con, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			mcs := NewMultiChannelStream(con)
			go mcs.Start()
			
			s.handle(mcs)
		}()
	}
	return nil
}