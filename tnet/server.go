package tnet

import (
	"net"
	"log"
)

type Server struct {
	net string
	laddr string
	handle func(p *MultiChannelStream)
}

func NewServer(net, laddr string, handle func(p *MultiChannelStream)) (s *Server) {
	s = new(Server)
	s.net = net
	s.laddr = laddr
	s.handle = handle
	return
}

func (s *Server) Start() error {
	ln, err := net.Listen(s.net, s.laddr)
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