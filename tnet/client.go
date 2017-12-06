package tnet

import (
	"net"
)

func Connect(network string, address string) (net.Conn, error) {
	con, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return con, nil
}