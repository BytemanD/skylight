package controller

import (
	"encoding/json"
	"skylight/internal/service/openstack"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type Networking struct{}

func (c *Networking) Get(req *ghttp.Request) {
	proxyUrl := strings.TrimPrefix(req.URL.String(), "/networking")
	manager := openstack.GetManager()
	resp, err := manager.ProxyNetworking(proxyUrl)

	if err != nil {
		data, _ := json.Marshal(HttpError{Code: 400, Message: err.Error()})
		req.Response.WriteStatus(400, data)
	} else {
		req.Response.WriteStatus(resp.StatusCode(), resp.Body())
		req.Response.Header().Set("Content-Type", resp.Header().Get("Content-Type"))
	}
}
