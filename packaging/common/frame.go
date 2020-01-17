package common

import "bytes"

// Frame is wrapper for received bytes from underlying socket
type Frame struct {
	Data *bytes.Buffer
}