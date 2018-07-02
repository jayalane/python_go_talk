
func getConnPool(hostName string) (conn net.Conn) {
	connPoolMutex.Lock() // hash not thread safe
	defer connPoolMutex.Unlock()

	poolChan, ok := connPool[hostName] // check for pool

	if ok { // found pool check for conn
		conn = nil
		select {
		case conn = <-poolChan:
			theContext.connectionsReused++
		case <-time.After(time.Millisecond * 10):
			conn = nil
		}
		return conn
	}
	connPool[hostName] = make(chan net.Conn, theConfig["maxConnToUse"].IntVal) // for return
	return nil
}
