package node

import (
	"net"
	"time"
)

type dataport struct {
	read  chan []byte
	write chan []byte
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
