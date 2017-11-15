package tnet

import (
	"net"
	"log"
)

func Connect(network string, address string) (*MultiChannelStream, error) {
	con, err := net.Dial(network, address)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	
	mcs := NewMultiChannelStream(con)
	go mcs.Start()
	return mcs, nil
}