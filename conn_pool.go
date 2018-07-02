// -*- tab-width: 2 -*-

package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	connPool      = map[string]chan net.Conn{}
	connPoolMutex = sync.RWMutex{}
	NextConnList  []string
	NextConnIndex int
)

// parsedMessageListener is channel listener for parsed syslog messages to send for front end -- there are >1 of these running
func parsedMessageListenerFront() {
	fmt.Println("PML-front started")
	for {
		select {
		case msg := <-theContext.parsedMessages:
			theContext.recvParsedMessagesDequeued++
			sendOneMessage(msg)
		case <-time.After(time.Second * 60):
			fmt.Println("60 seconds with no messages")
		}
	}
}

func getConnPool(hostName string) (conn net.Conn) {
	connPoolMutex.Lock() // hash not thread safe
	defer connPoolMutex.Unlock()

	poolChan, ok := connPool[hostName] // check for pool

	if ok { // found pool check for conn
		conn = nil
		select {
		case conn = <-poolChan:
			theContext.connectionsReused++
		case <-time.After(time.Millisecond * 10):
			conn = nil
		}
		return conn
	}
	connPool[hostName] = make(chan net.Conn, theConfig["maxConnToUse"].IntVal) // for return
	return nil
}

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

// fiddle with a parsed message to make it have the extra data

func tweakMessage(theMsg string) string {
	if len(theMsg) < 10 {
		return theMsg
	}
	var newMsg string
	new_str := fmt.Sprintf("t=%f r=%s ",
		float64(time.Now().UnixNano())/1000000.0,
		genId())

	newMsg = theMsg[0:4] + new_str + theMsg[4:]
	return newMsg
}

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
