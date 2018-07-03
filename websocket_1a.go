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
		}
	}
}
y
