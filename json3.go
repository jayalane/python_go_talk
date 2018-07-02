
// EditJSON returns whitelisted bits of JSON as per InfoSec rules
func EditJSON(msg string) (string, error) {

	var fullDatamap map[string]interface{}

	json.Unmarshal([]byte(msg), &fullDatamap)

	editedDatamap, err := editJSONRecur(fullDatamap, "")

	if err != nil {
		return "", err
	}
	editedDatamapBytes, err := json.Marshal(editedDatamap)
	if err != nil {
		log.Println("Error in JSON parse", err)
		return "", err
	}
	return string(editedDatamapBytes), nil
}
