package common

import (
	"log"
	"github.com/guftall/socket-message-boundary-go/packaging/common"
)

type TcpReceiver struct {
	lastFrame *common.Frame
	out chan<- *common.Frame
}

func (tr *TcpReceiver) SetChannel(ch chan<- *common.Frame) {
	tr.out = ch
}

func (tr *TcpReceiver) Received(f *common.Frame) {
	log.Printf("receiver (%p) received frame (%p) by len %d\n", tr, f, f.Data.Len())
	tr.lastFrame = f
	tr.sendToChannel(f)
}

func (tr *TcpReceiver) sendToChannel(f *common.Frame) {
	if tr.out == nil {
		panic("sending channel is null")
	}
	log.Println("sending frame to channel")
	tr.out <- f
}