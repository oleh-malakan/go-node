package node

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

type tError struct {
	text string
}

func (e *tError) Error() string {
	return e.text
}

func newError(text string) *tError {
	return &tError{text}
}
