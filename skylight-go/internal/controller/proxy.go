package controller

import (
	"fmt"
	"skylight/internal/service/openstack"
	"skylight/utility/easyhttp"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

type HttpError struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewHttpIntervalError() HttpError {
	return HttpError{Error: "internal error"}
}

type OpenstackProxy struct {
	Prefix string
}

func (c *OpenstackProxy) doProxy(req *ghttp.Request) {
	var resp *easyhttp.Response
	var err error
	sessionId, err := req.Session.Id()
	if err != nil {
		glog.Error(req.GetCtx(), "get session failed: %s", err)
		req.Response.WriteStatusExit(500, NewHttpIntervalError())
	}

	manager, err := openstack.GetManager(sessionId, req)
	if err != nil {
		glog.Errorf(req.GetCtx(), "get manager failed: %s", err)
		req.Response.WriteStatusExit(500, NewHttpIntervalError())
	}

	proxyUrl := strings.TrimPrefix(req.URL.Path, c.Prefix)
	switch c.Prefix {
	case "/computing":
		resp, err = manager.ProxyComputing(req.Method, proxyUrl, req.URL.Query(), req.GetBody())
	case "/networking":
		resp, err = manager.ProxyNetworking(req.Method, proxyUrl, req.URL.Query(), req.GetBody())
	case "/volume":
		resp, err = manager.ProxyVolume(req.Method, proxyUrl, req.URL.Query(), req.GetBody())
	case "/image":
		resp, err = manager.ProxyImage(proxyUrl, req)
	case "/identity":
		resp, err = manager.ProxyIdentity(req.Method, proxyUrl, req.URL.Query(), req.GetBody())
	default:
		err = fmt.Errorf("invalid prefix %s", c.Prefix)
	}
	if err != nil {
		req.Response.WriteStatusExit(400, HttpError{Error: err.Error()})
	} else if resp == nil {
		req.Response.WriteStatusExit(204, "ok")
	} else {
		req.Response.WriteStatusExit(resp.StatusCode(), resp.Body())
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
func (c *OpenstackProxy) Patch(req *ghttp.Request) {
	c.doProxy(req)
}
