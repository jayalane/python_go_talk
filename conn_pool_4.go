
func putConn(conn net.Conn, hostName string) {
	connPoolMutex.Lock()
	defer connPoolMutex.Unlock()
	poolChan, ok := connPool[hostName]
	if ok {
		select {
		case poolChan <- conn:
			theContext.connectionsReuseQueued++
		case <-time.After(time.Millisecond * 1):
			conn.Close()
			theContext.connectionsReuseClosed++
		}
	}
}
