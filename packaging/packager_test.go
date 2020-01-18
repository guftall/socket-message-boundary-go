package packaging

import (
	"bytes"
	"log"
	"io/ioutil"
	"github.com/guftall/socket-message-boundary-go/packaging/common"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestFillCurrentPacketSizeWithFrameLen1(t *testing.T) {
	p := Packager{}
	pkt := common.Packet{}

	p.currentPacket = &pkt

	f := common.Frame{}
	f.Data = bytes.NewBuffer([]byte{})

	b := []byte{1}
	f.Data.Write(b)

	p.fillCurrentPacketSize(&f)

	if p.currentPacket.Size >> 8 != uint16(b[0]) {
		t.Errorf("packager should set most significant byte of package size\n")
	}
}

func TestFillCurrentPacketDataAllDataInFrame(t *testing.T) {
	pkt := common.Packet {}
	pkt.Size = 64

	f := common.Frame{}
	f.Data = bytes.NewBuffer([]byte{})
	f.Data.Write(make([]byte, 64))


	// Act
	fillCurrentPacketData(&pkt, &f)

	// Assert
	if pkt.State != common.Complete {
		t.Errorf("packet state not changed after receiving all data in frame\n")
	}
}

func BenchmarkFillCurrentPacketDataAllDataInFrame(b *testing.B) {

	pkt := common.Packet {}
	pkt.Size = 64

	f := common.Frame{}
	f.Data = bytes.NewBuffer([]byte{})

	d := make([]byte, 64)
	dd := []byte{}

	// b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Data.Reset()
		f.Data.Write(d)
		pkt.Data = dd
		fillCurrentPacketData(&pkt, &f)
	}
}

func BenchmarkFillCurrentPacketDataNotCompleteFrame(b *testing.B) {

	pkt := common.Packet {}
	pkt.Size = 64

	f := common.Frame{}
	f.Data = bytes.NewBuffer([]byte{})

	d := make([]byte, 32)
	dd := []byte{}

	// b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Data.Reset()
		f.Data.Write(d)
		pkt.Data = dd
		fillCurrentPacketData(&pkt, &f)
	}
}

func BenchmarkFillCurrentPacketSizeWithFrameLen1(b *testing.B) {
	p := Packager{}
	pkt := common.Packet{}

	p.currentPacket = &pkt

	f := common.Frame{}
	f.Data = bytes.NewBuffer([]byte{})

	f.Data.Write([]byte{1})

	for i := 0; i < b.N; i++ {

		p.fillCurrentPacketSize(&f)
	}
}
