package udp

import (
	"net"
	"time"

	"github.com/wmentor/dsn"
)

type Client struct {
	con     *net.UDPConn
	addr    *net.UDPAddr
	timeout time.Duration
	buffer  []byte
}

func New(opts string) (*Client, error) {

	params, err := dsn.New(opts)
	if err != nil {
		return nil, err
	}

	addr := params.GetString("addr", "")
	if addr == "" {
		return nil, ErrInvalidAddr
	}

	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	udpBufferSize := params.GetInt("udp_buffer_size", 1024*1024)
	if udpBufferSize < 1024 {
		return nil, ErrInvalidUdpBufferSize
	}

	maxMessageSize := params.GetInt("max_nessage_size", 64*1024)
	if maxMessageSize < 1024 {
		return nil, ErrInvalidMaxMessageSize
	}

	timeout := params.GetInt("udp_timeout", 5)
	if timeout < 1 {
		return nil, ErrInvaidUdpTimeout
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}

	conn.SetReadBuffer(udpBufferSize)
	conn.SetWriteBuffer(udpBufferSize)

	return &Client{
		con: conn, addr: raddr,
		buffer:  make([]byte, maxMessageSize),
		timeout: time.Second * time.Duration(timeout),
	}, nil
}

func (c *Client) Write(msg []byte) (int, error) {
	c.con.SetWriteDeadline(time.Now().Add(c.timeout))
	return c.con.Write(msg)
}

func (c *Client) Receive() ([]byte, error) {
	if n, err := c.Read(c.buffer); err != nil {
		return nil, err
	} else {
		return c.buffer[:n], nil
	}
}

func (c *Client) Read(buf []byte) (int, error) {
	c.con.SetReadDeadline(time.Now().Add(c.timeout))
	n, _, err := c.con.ReadFromUDP(buf)
	return n, err
}

func (c *Client) Close() {
	if c.con != nil {
		c.con.Close()
		c.con = nil
	}
}
