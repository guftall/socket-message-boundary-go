package main

import (
	"time"
	"log"
	"flag"
	"strconv"

	"github.com/guftall/socket-message-boundary-go/networking/client"
	"github.com/guftall/socket-message-boundary-go/networking/common"
	_ "github.com/guftall/socket-message-boundary-go/packaging"
)


func main() {

	// b := byte(128)
	// i := uint16(b) << 8
	// i |= uint16(b)
	// println(i)
	// return

	host := flag.String("h", "localhost", "destination host to connect to")
	port := flag.Int("p", 5959, "port of server application")

	flag.Parse()

	client, err := client.Connect("tcp", (*host) + ":" + strconv.Itoa(*port))
	if err != nil {
		log.Fatalln(err)

	}

	log.Printf("connected to %s\n", client.LocalAddr())

	client.Write([]byte{0})
	for i := range [9999999999]int{} {
		i++
	}
	client.Write([]byte{3, 2, 3, 3})
	for i := range [9999999999]int{} {
		i++
	}
	client.Write([]byte{0})
	time.Sleep(10)
	client.Write([]byte{1})

	frame, err := common.ReadNetworkFrame(&client)
	if err != nil {
		log.Println("error: ", err)
	}

	log.Printf("received frame %s\n", frame)

	_ = client
}