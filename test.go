func doTestFrontier(t *testing.T, done2Chan chan bool, evil bool) {
	a, b := getProfile()
	f, err := NewFrontierClient("sherlock-frontier-vip.qa.paypal.com",
		80,
		"pp", "qa",
		a, b)

	if err != nil {
		t.Log("Can't connect to frontier", err)
		t.Fail()
		return
	}
	var dones = []chan bool{}
	for i := 0; i <= numTrials; i++ {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		doneChan := make(chan bool, 1)
		go sendSomeStuff(f, doneChan)
		if evil {
			rn := rand.Intn(1000)
			if rn > 995 {
				f.lock()
				f.wsConn.close()
				f.unlock()
				fmt.Println("Closing socket on", f.clientID)
