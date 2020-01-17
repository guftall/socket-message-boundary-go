package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/guftall/socket-message-boundary-go/networking/common"
	pack "github.com/guftall/socket-message-boundary-go/packaging"
	packcommon "github.com/guftall/socket-message-boundary-go/packaging/common"
)

// StartListening in 0.0.0.0:port address
func StartListening(port int) {

	addr := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("unable to listen on address %s", addr)
	}

	log.Printf("now listening on %s\n", listener.Addr())

	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Printf("can't accept connection: %s\n", err)
			continue
		}

		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	log.Printf("Serving %s\n", c.RemoteAddr())
	r := common.TcpReceiver{}
	p := pack.Packager{}

	defer closeConnection(c)
	defer p.Stop()
	
	chFrame := make(chan *packcommon.Frame)
	chPacket := make(chan *packcommon.Packet, 2)

	r.SetChannel(chFrame)
	p.Init(chFrame, chPacket)

	p.Start()

	for {
		f, err := common.ReadNetworkFrame(&c)
		if err != nil {
			if err == io.EOF {
				log.Printf("connection(%p) EOF\n", c)
				return
			}
			log.Printf("unable to read network frame: %s\n", err)
			return
		}

		r.Received(f)
	}
}

func closeConnection(conn net.Conn) {
	log.Printf("closing connection(%p)\n", conn)
	conn.Close()
}
