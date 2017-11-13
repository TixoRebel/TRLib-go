package tnet

import (
	"net"
	"log"
)

type Client struct {
	network string
	address string
	handle func(p *MultiChannelStream)
}

func NewClient(network string, address string, handle func(p *MultiChannelStream)) (c *Client) {
	c = new(Client)
	c.network = network
	c.address = address
	c.handle = handle
	return
}

func (c *Client) Connect() error {
	con, err := net.Dial(c.network, c.address)
	if err != nil {
		log.Fatal(err)
		return err
	}
	
	mcs := NewMultiChannelStream(con)
	go mcs.Start()
	c.handle(mcs)
	return nil
}