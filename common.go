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
	cRead  chan []byte
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

func compareID(a []byte, b []byte) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3] &&
		a[4] == b[4] && a[5] == b[5] && a[6] == b[6] && a[7] == b[7] &&
		a[8] == b[8] && a[9] == b[9] && a[10] == b[10] && a[11] == b[11] &&
		a[12] == b[12] && a[13] == b[13] && a[14] == b[14] && a[15] == b[15] &&
		a[16] == b[16] && a[17] == b[17] && a[18] == b[18] && a[19] == b[19] &&
		a[20] == b[20] && a[21] == b[21] && a[22] == b[22] && a[23] == b[23] &&
		a[24] == b[24] && a[25] == b[25] && a[26] == b[26] && a[27] == b[27] &&
		a[28] == b[28] && a[29] == b[29] && a[30] == b[30] && a[31] == b[31]
}
