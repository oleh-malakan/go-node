package node

import (
	"net"
	"time"
)

type Connection struct{}

func (c *Connection) Send(b []byte) error {
	return nil
}

func (c *Connection) Receive() ([]byte, error) {
	return nil, nil
}

func (c *Connection) Close() error {
	return nil
}

type dataport struct {
	cRead chan []byte
	cWrite chan []byte
}

func (d *dataport) Read(b []byte) (n int, err error) {
	return
}

func (d *dataport) Write(b []byte) (n int, err error) {
	return
}

func (d *dataport) Close() error {
	return nil
}

func (d *dataport) LocalAddr() net.Addr {
	return nil
}

func (d *dataport) RemoteAddr() net.Addr {
	return nil
}

func (d *dataport) SetDeadline(t time.Time) error {
	return nil
}

func (d *dataport) SetReadDeadline(t time.Time) error {
	return nil
}

func (d *dataport) SetWriteDeadline(t time.Time) error {
	return nil
}
