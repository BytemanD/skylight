package utility

import "encoding/json"

func StructToJson(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), err
}
