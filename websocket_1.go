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
