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

type TelnetClientImplement struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	connection net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImplement{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *TelnetClientImplement) Connect() (err error) {
	t.connection, err = net.DialTimeout("tcp", t.address, t.timeout)
	return err
}

func (t *TelnetClientImplement) Close() error {
	return t.connection.Close()
}

func (t *TelnetClientImplement) Send() error {
	_, err := io.Copy(t.connection, t.in)
	return err
}

func (t *TelnetClientImplement) Receive() error {
	_, err := io.Copy(t.out, t.connection)
	return err
}
