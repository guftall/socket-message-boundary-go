package packaging

import (
	"fmt"
	"io"
	"log"

	"github.com/guftall/socket-message-boundary-go/packaging/common"
)

// Ethernet MTU = 1500

// MaxPacketLength is maximum size that a packet can be
const MaxPacketLength = 1500 + 1024

// Packager is responsible to extract packages from frames
// out channel should have capacity more than 1, so frames can be readed from underlying networking channel
// after completing read of one package.
type Packager struct {
	lastFrame     *common.Frame
	currentPacket *common.Packet
	in            <-chan *common.Frame
	out           chan<- *common.Packet
	reading       bool
}

// Init initializes Packager
func (p *Packager) Init(in <-chan *common.Frame, out chan<- *common.Packet) {
	p.in = in
	p.out = out
}

// Start receiving frames from underlying channel
func (p *Packager) Start() {
	p.reading = true
	go p.startReceiving()
}

// Stop reading from undertlying channel
func (p *Packager) Stop() {
	log.Printf("stopping packager\n")
	p.stop()
}

func (p *Packager) startReceiving() {

	for p.reading {

		var frame *common.Frame
		frame, ok := <-p.in
		if !ok {
			log.Printf("packager input channel closed\n")
			p.close()
		}

		log.Printf("read frame (%p) by Len %d from channel\n", frame, frame.Data.Len())

		if frame.Data.Len() < 0 {
			log.Printf("frame with negative size!\n")
			p.close()
		}

		if frame.Data.Len() == 0 {
			log.Printf("empty frame received, is there anything wrong??\n")
			continue
		}

		for {
			if frame.Data.Len() <= 0 {
				break
			}

			if p.currentPacket == nil {
				log.Println("creating new packet")
				p.currentPacket = &common.Packet{}
				p.currentPacket.State = common.Size
			}

			if p.currentPacket.State == common.Size {
				err := p.fillCurrentPacketSize(frame)
				if err != nil {
					log.Printf("unable to fill packet size: %s", err)
					p.close()
					return
				}

				if p.currentPacket.State == common.Data {
					log.Printf("packet (%p) size calculated: %d\n", p.currentPacket, p.currentPacket.Size)

					// if length of packet exceeds MAX value, close connection
					if p.currentPacket.Size > MaxPacketLength {
						log.Printf("packet size(%d) exceeds MAX size(%d)\n", p.currentPacket.Size, MaxPacketLength)
						p.close()
						return
					}
				}
			}

			if p.currentPacket.State == common.Data {
				p.fillCurrentPacketData(frame)
			}

			if p.currentPacket.State == common.Complete {
				p.packetCompleted()
			}
		}
	}
}

func (p *Packager) fillCurrentPacketSize(f *common.Frame) error {

	if f.Data.Len() == 0 {
		return nil
	}
	// if lastFrame is not nil, so previous frame had size one,
	// because our size is 2byte, so after receiving this new frame packet size completed

	if p.lastFrame != nil {
		b, err := f.Data.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return fmt.Errorf("can't read from frame Data: %s", err)
		}

		p.currentPacket.Size |= uint16(b)
		p.currentPacket.State = common.Data
		return nil
	}

	if f.Data.Len() == 1 {
		log.Printf("packet(%p) size not can't be calculated yet\n", p.currentPacket)
		p.lastFrame = f
		b, _ := f.Data.ReadByte()

		// b is most singnificant byte
		p.currentPacket.Size = uint16(b) << 8
		return nil
	}

	b0, _ := f.Data.ReadByte()
	b1, _ := f.Data.ReadByte()
	p.currentPacket.Size = (uint16(b0) << 8) | uint16(b1)

	p.currentPacket.State = common.Data
	return nil
}

func (p *Packager) fillCurrentPacketData(f *common.Frame) {
	if len(p.currentPacket.Data) == int(p.currentPacket.Size) || f.Data.Len() == 0 {
		return
	}

	remaining := int(p.currentPacket.Size) - len(p.currentPacket.Data)

	count := remaining - f.Data.Len()

	if count <= 0 {
		// so all remaining data is in current frame

		log.Printf("read remaining packet data(%d byte) from frame\n", remaining)
		slice := make([]byte, remaining, remaining)

		f.Data.Read(slice)

		p.currentPacket.Data = append(p.currentPacket.Data, slice...)

		p.currentPacket.State = common.Complete
	} else {
		log.Printf("packet(%p) not completed by frame(%p), wait for next frame", p.currentPacket, f)
		p.lastFrame = f

		// means packet remaining bytes are more than this frame, so read all frame data and append them to
		// current packet
		slice := f.Data.Bytes()

		// we read remaining bytes to slice variable, so discard all remaining bytes
		f.Data.Truncate(0)
		p.currentPacket.Data = append(p.currentPacket.Data, slice...)
	}
}

// packetCompleted when all packet data readed and should start of reading new packet from frame
func (p *Packager) packetCompleted() {
	log.Printf("packet completed by size %d", p.currentPacket.Size)
	p.out <- p.currentPacket

	// current packet indicates an uncompleted packet, so we should set it to nil when packet reassembled completly
	p.currentPacket = nil

	// we should set lastFrame to nil, because this field only used when we have an uncompleted packet that some of
	// its data received in previous frame.
	p.lastFrame = nil
}

func (p *Packager) close() {
	log.Printf("closing packager receive channel\n")
	p.stop()
}

func (p *Packager) stop() {
	p.reading = false
}
