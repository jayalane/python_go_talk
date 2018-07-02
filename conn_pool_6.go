// parsedMessageListener is channel listener for parsed syslog messages to send for front end -- there are >1 of these running
func parsedMessageListenerFront() {
	var sleepTime time.Duration
	sleepTime = time.Second * 5
	fmt.Println("PML-front started")
	for i := 0; i < len(nextConnList); i++ {
		for j := 0; j < theConfig["maxConnToUse"].IntVal; j++ {
			// each connection in the pool is a go-routine -- all state is local.
			go func(endPoint string) {
				var conn net.Conn
				var err error
				for {
					if conn == nil {
						conn, err = withCalTxn(calProtocol.TxnTypeConnect, endPoint, func(string) (net.Conn, error) {
							conn, err := getConn(endPoint)
							return conn, err
						})
						if err != nil {
							time.Sleep(sleepTime) // mark down
							if sleepTime < time.Minute*5 {
								sleepTime = sleepTime + sleepTime
							}
							continue
						} else {
							theContext.connectionsReused++
						}
					}
					select {
					case msg := <-theContext.parsedMessages:
						theContext.recvParsedMessagesDequeued++
						if sendOneMessage(msg.msg, conn) {
							conn.Close()
							conn = nil
						}
						theContext.timeSeries.MarkDistribution(time.Now().Sub(msg.when).Seconds()*1000.0,
							"sent_minus_recv")
						sleepTime = time.Second
					case <-time.After(time.Second * 60):
						fmt.Println("60 seconds with no messages")
					}
				}
			}(nextConnList[i])
		}
		fmt.Println("Started", theConfig["maxConnToUse"].IntVal, "go routines for", nextConnList[i])
	}
}
