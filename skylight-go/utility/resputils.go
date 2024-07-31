package utility

import (
	"encoding/json"
)

func JsonErrorResponse(errorMsg string) []byte {
	body, _ := json.Marshal(map[string]string{"error": errorMsg})
	return body
}
func JsonResponse(key string, data interface{}) []byte {
	body, _ := json.Marshal(map[string]interface{}{key: data})
	return body
}
