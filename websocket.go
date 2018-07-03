
func (m *FrontierClient) startConnect() {
	select {
	case m.doConnect <- true:
		// good
	default:
		// bad but ok
	}
}

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
			ws := m.wsConn
			m.runlock()
			bytes, err = ws.readMsg()
		}
		e, ok := err.(net.Error)
		if ok && e.Timeout() {
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

func (m *FrontierClient) mainLoopWrite() {
	go m.mainLoopRead() // also never exits
	m.doConnect <- true
	connectionPending := false
	for {
		var err error
		select {
		case <-m.gotConnect:
			connectionPending = false
		case <-m.doConnect:
			if connectionPending == false {
				go m.connect()
				connectionPending = true
			}
		case <-m.finished:
			fmt.Println("Done - exiting")
			return
		case metric := <-m.metrics:
			if m.checkSession() {
				msg := m.metricToMsg(metric)
				m.lockSession()
				err = m.wsConn.writeMsg(m.clientID, msg)
				m.unlockSession()
				if err != nil {
					fmt.Println("Write msg got error", m.clientID, err)
					m.restart_internal()                                                                   // kill connection
					m.resetCb(msg)                                                                         // clear and ignore CB
					err = m.SendWithCb(metric.dim, metric.data, metric.when, metric.resolution, metric.cb) // reQ
					if err != nil {
						fmt.Println("dropping retried metric - buffer full", m.clientID)
					}
					m.startConnect()
				}
			} else { // no session yet
				go func() {
					time.Sleep(5 * time.Second)
					err = m.SendWithCb(metric.dim, metric.data, metric.when, metric.resolution, metric.cb) // reQ
					if err != nil {
						fmt.Println("dropping retried metric - buffer full", m.clientID)
					}
				}()
			}
		case readMsg := <-m.readMsgs:
			var err error
			if readMsg.err == nil {
				err = m.readOneMsg(readMsg.msg)
				if err != nil {
					if !connectionPending {
						fmt.Println("Error with parse", m.clientID, err)
						m.restart_internal()
						m.startConnect()
					}
				}
			} else {
				if !connectionPending {
					fmt.Println("Error with read", m.clientID, readMsg.err, len(m.doConnect))
					m.restart_internal()
					m.startConnect()
				}
			}
		case <-time.After(time.Second * 300):
			fmt.Println("No activity for 5 minutes", m.clientID)
			continue
		}
	}
}
