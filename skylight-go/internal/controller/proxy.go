package controller

import (
	"skylight/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
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
	resp, err := service.OSService.DoProxy(req, c.Prefix)
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
