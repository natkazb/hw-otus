package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetConfig struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetConfig{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *TelnetConfig) Connect() (err error) {
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	return
}

func (t *TelnetConfig) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

func (t *TelnetConfig) Send() error {
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *TelnetConfig) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	return err
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
