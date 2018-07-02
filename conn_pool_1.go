var (
	connPool      = map[string]chan net.Conn{}
	connPoolMutex = sync.RWMutex{}
	NextConnList  []string
	NextConnIndex int
)
