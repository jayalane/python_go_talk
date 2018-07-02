			// TODO add list
		default:
			if stringInList(prefixSoFar+k, okKeys) { // whitelisted ok keys
				newJSON[k] = v
			} // default
		} // switch on type
	} // range jsonObj
	return newJSON, nil
}


func stringInList(str string, aList []string) bool {
	for _, b := range aList {
		if b == str {
			return true
		}
	}
	return false
}
