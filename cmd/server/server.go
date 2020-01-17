package main

import (
	"github.com/guftall/socket-message-boundary-go/networking/server"
)

func main() {
	server.StartListening(5959)
}