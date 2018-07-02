// actually send a message on from front to back

func sendOneMessage(theMsg string) {
	conn, hostName, err := getConn()
	if err != nil {
		fmt.Println("No connection dropping msg:", err)
		theContext.writeDrops++
	} else {
		theMsg = tweakMessage(theMsg)
		_, err = conn.Write([]byte(theMsg))
		if err == nil {
			_, err = conn.Write([]byte(theContext.syslogSeparatorString))
		}
		if err != nil {
			theContext.writeErrors++
			fmt.Println("Error on write:", err)
			if conn != nil {
				defer conn.Close()
			}
			return
		}
		theContext.writeMessages++
		theContext.writeBytes += len(theMsg)
		putConn(conn, hostName)
	}
}
