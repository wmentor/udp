package udp

import (
	"errors"
)

var (
	ErrInvalidAddr           error = errors.New("invalid addr")
	ErrInvalidUdpBufferSize  error = errors.New("invalid udp buffer size")
	ErrInvalidMaxMessageSize error = errors.New("invalid max message size")
	ErrInvaidUdpTimeout      error = errors.New("invalid udp timeout")
)
