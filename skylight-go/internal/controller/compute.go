package controller

import (
	"encoding/json"
	"skylight/internal/service/openstack"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

type Computing struct{}

func (c *Computing) Get(req *ghttp.Request) {
	proxyUrl := strings.TrimPrefix(req.URL.String(), "/computing")
	manager := openstack.GetManager()
	resp, err := manager.ProxyComputing(proxyUrl)

	if err != nil {
		data, _ := json.Marshal(HttpError{Code: 400, Message: err.Error()})
		req.Response.WriteStatus(400, data)
	} else {
		req.Response.WriteStatus(resp.StatusCode(), resp.Body())
		req.Response.Header().Set("Content-Type", resp.Header().Get("Content-Type"))
	}
}
func (c *Computing) Delete(req *ghttp.Request) {
	proxyUrl := strings.TrimPrefix(req.URL.String(), "/computing")
	manager := openstack.GetManager()
	resp, err := manager.ProxyComputing(proxyUrl)

	if err != nil {
		data, _ := json.Marshal(HttpError{Code: 400, Message: err.Error()})
		req.Response.WriteStatus(400, data)
	} else {
		req.Response.WriteStatus(resp.StatusCode(), resp.Body())
		req.Response.Header().Set("Content-Type", resp.Header().Get("Content-Type"))
	}
}
