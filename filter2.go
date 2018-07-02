func filterDoMessage(msg string) bool {
	for _, v := range literalBadStrings {
		if strings.Contains(msg, v) {
			return false
		}
	}
	for _, v := range literalDoubleBadStrings {
		if strings.Contains(msg, v.needsOne) && (strings.Contains(msg, v.needsTwo)) {
			return false
		}
	}
	for _, v := range literalAndNotBadStrings {
		if strings.Contains(msg, v.needsOne) && (!strings.Contains(msg, v.needsTwo)) {
			return false
		}
	}
	return true
}
