func getConn() (net.Conn, string, error) {
	var conn net.Conn
	var err error
	err = nil
	conn = nil
	host := NextConnList[NextConnIndex%len(NextConnList)]
	NextConnIndex++
	if NextConnIndex == 1024*1024 {
		NextConnIndex = 0
	}
	conn = getConnPool(host)
	if conn == nil {
		theContext.connectionsMade++
		port := theConfig["relayPort"].StrVal
		addr := strings.Join([]string{host, port}, ":")
		fmt.Println("About to dial", addr)
		conn, err = net.Dial("tcp", addr)
	}
	if err != nil {
		theContext.connectionsError++
		if conn != nil {
			defer conn.Close()
		}
		return nil, "", err
	}
	return conn, host, nil
}
