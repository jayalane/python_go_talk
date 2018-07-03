// in production file
var genID = realGenID // allows for test

func realGenID() string {
	hn := HostName
	h := sha1.New()
	io.WriteString(h, "go")
	io.WriteString(h, fmt.Sprintf("%f", float64(theContext.startTime.UnixNano())/1000000000.0))
	io.WriteString(h, hn)
	io.WriteString(h, strconv.Itoa(seqNum))
	r := fmt.Sprintf("%x", h.Sum(nil))

	return r
}

// in test file
func newGenID() string {
	return "someData"
}

func TestFixupApache(t *testing.T) {
	genID = newGenID
...
