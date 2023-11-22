package node

type Handler struct {
	nodeID [32]byte
	f      func(connection *Connection)
}

func (h *Handler) Close() error {
	return nil
}