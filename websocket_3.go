
func (m *FrontierClient) startConnect() {
	select {
	case m.doConnect <- true:
		// good
	default:
		// bad but ok
	}
}
