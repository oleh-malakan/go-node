package node

var (
	server *Server
)

func init() {
	server = &Server{
		memoryLock: make(chan *struct{}, 1),
		checkLock:  make(chan *struct{}, 1),
	}
}

