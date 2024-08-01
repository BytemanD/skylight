package controller

import (
	"encoding/json"
	"skylight/internal/service/openstack"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gogf/gf/v2/net/ghttp"
)

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type OpenstackProxy struct {
	Prefix string
}

func (c *OpenstackProxy) doProxy(req *ghttp.Request) {
	manager := openstack.GetManager()
	var resp *resty.Response
	var err error

	proxyUrl := strings.TrimPrefix(req.URL.Path, c.Prefix)
	switch c.Prefix {
	case "/computing":
		resp, err = manager.ProxyComputing(req.Method, proxyUrl, req.URL.Query(), req.GetBody())
	case "/networking":
		resp, err = manager.ProxyNetworking(req.Method, proxyUrl, req.URL.Query(), req.GetBody())
	case "/volume":
		resp, err = manager.ProxyVolume(req.Method, proxyUrl, req.URL.Query(), req.GetBody())
	case "/image":
		resp, err = manager.ProxyImage(req.Method, proxyUrl, req.URL.Query(), req.GetBody())
	}
	if err != nil {
		data, _ := json.Marshal(HttpError{Code: 400, Message: err.Error()})
		req.Response.WriteStatus(400, data)
	} else {
		req.Response.WriteStatus(resp.StatusCode(), resp.Body())
		req.Response.Header().Set("Content-Type", resp.Header().Get("Content-Type"))
	}
}

func (c *OpenstackProxy) Get(req *ghttp.Request) {
	c.doProxy(req)
}
func (c *OpenstackProxy) Delete(req *ghttp.Request) {
	c.doProxy(req)
}
func (c *OpenstackProxy) Post(req *ghttp.Request) {
	c.doProxy(req)
}
func (c *OpenstackProxy) Put(req *ghttp.Request) {
	c.doProxy(req)
}
