package node

func Handler(nodeID string, f func(queryData []byte, connection *Connection)) {}

func ListenAndServe() error {
	return nil
}

type Client struct{}

func (c *Client) Query(nodeID string, queryData []byte) (*Connection, error) {
	return &Connection{}, nil
}

func Dial() (*Client, error) {
	return &Client{}, nil
}

type Connection struct{}

func (c *Connection) Send(data []byte) error {
	return nil
}

func (c *Connection) Receive() ([]byte, error) {
	return nil, nil
}

func (c *Connection) Close() error {
	return nil
}
