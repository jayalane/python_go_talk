func (m *FrontierClient) mainLoopRead() {
	var err error
	var bytes []byte
	for {
		m.rlock()
		if m.wsConn.ws == nil {
			m.runlock()
			time.Sleep(time.Second * 1)
			continue
		} else {
			m.wsConn.ws.SetReadDeadline(time.Now().Add(5 * time.Second))
			ws := m.wsConn  // keep copy while have lock
			m.runlock()
			bytes, err = ws.readMsg()
		}
		if e, ok := err.(net.Error); ok && e.Timeout() {
			runtime.Gosched()
			continue // ignore read timeouts
		}
		r := readMsg{msg: bytes, err: err}
		m.readMsgs <- r // block? that's ok don't read till can handle
		if err != nil {
			fmt.Println("error for read", m.clientID, err)
			// no restart_internal - will be handled in sync with write loop
			time.Sleep(2 * time.Second) // takes some time to reconnect
		}
	}
}
