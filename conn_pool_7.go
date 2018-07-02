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
