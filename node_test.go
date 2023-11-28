package node

var (
	cont *serverController
)

func init() {
	cont = &serverController{
		config: &Config{
			ClientsLimit: 524288,
			HeapCap:      512,
		},
		in:       make(chan *incomingPackage),
		nextDrop: make(chan *serverContainer),
	}
}
