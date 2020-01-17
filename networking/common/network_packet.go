package common

import (
	"log"
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/guftall/socket-message-boundary-go/packaging/common"
)

// type NetworkFrame struct {
// 	Size int16
// 	Data []byte
// }

// ReadNetworkFrame reads a frame from underlying connection
func ReadNetworkFrame(conn *net.Conn) (*common.Frame, error) {
	p := common.Frame{}

	buf := make([]byte, 1024)

	n, err := (*conn).Read(buf)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("unable to read connection %v: %v", conn, err)
		}

		return nil, err
	}

	log.Printf("read %d byte from connection", n)

	p.Data = bytes.NewBuffer(buf[0:n])

	return &p, nil
}
