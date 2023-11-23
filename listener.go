package node

type Listener struct{}

func (l *Listener) Accept() (*Connection, error) {
	return &Connection{}, nil
}

func (l *Listener) Close() error {
	return nil
}
