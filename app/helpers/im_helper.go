package helpers

import (
	"encoding/json"
	"os"
	"sort"
)

func ImSignEncryption(args interface{}) (string, error) {
	var convertText map[string]string
	jsonText, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(jsonText, &convertText)
	if err != nil {
		return "", err
	}
	keys := make([]string, 0, len(convertText))
	for k := range convertText {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	query := ""
	for _, val := range keys {
		if query == "" {
			query += val + "=" + convertText[val]
		} else {
			query += "&" + val + "=" + convertText[val]
		}
	}
	salt := os.Getenv("IM_SALT")
	query += "&pri_key=" + salt
	return GetMD5Hash(query), nil
}
