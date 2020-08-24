package udp

import (
	"fmt"
	"net"
	"sync"

	"github.com/wmentor/dsn"
	"github.com/wmentor/log"
)

type ServerHandler func(msg []byte) []byte

func Server(opts string, fn ServerHandler) error {

	params, err := dsn.New(opts)
	if err != nil {
		return err
	}

	addr := params.GetString("addr", "")
	if addr == "" {
		return ErrInvalidAddr
	}

	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	udpBufferSize := params.GetInt("udp_buffer_size", 1024*1024)
	if udpBufferSize < 1024 {
		return ErrInvalidUdpBufferSize
	}

	maxMessageSize := params.GetInt("max_nessage_size", 128*1024)
	if maxMessageSize < 1024 {
		return ErrInvalidMaxMessageSize
	}

	ln, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}

	ln.SetReadBuffer(udpBufferSize)
	ln.SetWriteBuffer(udpBufferSize)

	pool := sync.Pool{New: func() interface{} { return make([]byte, maxMessageSize) }}

	for {

		buffer := pool.Get().([]byte)

		n, addr, err := ln.ReadFromUDP(buffer)
		if err != nil {
			pool.Put(buffer)
			continue
		}

		go func() {

			defer func() {

				if r := recover(); r != nil {
					log.Error(fmt.Sprint(r))
				}

			}()

			defer pool.Put(buffer)

			msg := buffer[:n]

			msg = fn(msg)

			if len(msg) > 0 {
				ln.WriteToUDP(msg, addr)
			}

		}()

	}

	return nil

}
