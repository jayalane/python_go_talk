
func editJSONRecur(jsonObj map[string]interface{}, prefixSoFar string) (map[string]interface{}, error) {
	newJSON := make(map[string]interface{})
	for k, v := range jsonObj {
		newV := ""
		switch v.(type) {
		case map[string]interface{}:
			newMap, err := editJSONRecur(v.(map[string]interface{}), prefixSoFar+k+".")
			if err != nil {
				return nil, err
			}
			newJSON[k] = newMap
		case string:
			if stringInList(prefixSoFar+k, okKeys) { // whitelisted ok keys
				newV = v.(string)
				if stringInList(k, stripKeys) {
					splitV := strings.SplitN(v.(string), "?", 2)
					newV = splitV[0]
				}
				if stringInList(k, blanketKeys) {
					newV = "BLANKET"
				}
				newJSON[k] = newV
			} // string
