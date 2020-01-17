package common

// Packet is a completed packet received from sender
type Packet struct {
	Size uint16
	Data []byte
	State PacketState
}

const (
	Size = iota
	Data
	Complete
)

// PacketState to indicate state of packet
type PacketState int