package client

import (
	"fmt"
	"net"
)

func Connect(network string, destination string) (net.Conn, error) {
	conn, err := net.Dial(network, destination)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to [%s, %s] : %s", network, destination, err)
	}

	return conn, nil
}