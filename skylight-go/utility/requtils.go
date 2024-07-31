package utility

import (
	"encoding/json"
	"strconv"

	"github.com/gogf/gf/v2/net/ghttp"
)

func GetReqNamespace(req *ghttp.Request) string {
	namespace := req.URL.Query().Get("namespace")
	if namespace != "" {
		return namespace
	} else {
		return "default"
	}
}
func GetReqParamString(req *ghttp.Request, key string) *string {
	if !req.URL.Query().Has(key) {
		return nil
	}
	value := req.URL.Query().Get(key)
	return &value
}
func GetReqParamInt64(req *ghttp.Request, key string) *int64 {
	value := GetReqParamString(req, key)
	if value == nil {
		return nil
	}
	val, err := strconv.Atoi(*value)
	if err != nil {
		return nil
	}
	int64Val := int64(val)
	return &int64Val
}

func GetReqBody(req *ghttp.Request, respBody interface{}) error {
	reqBody := req.GetBody()
	if err := json.Unmarshal(reqBody, &respBody); err != nil {
		return err
	}
	return nil
}
