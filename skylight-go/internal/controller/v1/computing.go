package v1

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

type ServersController struct{}

// func (c *ServersController) Get(req *ghttp.Request) {
// 	session := GetAuthSession(req)
// 	if session == nil {
// 		req.Response.WriteStatus(403, "no login")
// 		return
// 	}

// 	api := session.OpenstackClient.NovaV2().Server().ResourceApi
// 	resp, err := api.AppendUrl("detail").SetQuery(req.URL.Query()).Get(nil)
// 	c.WriteProxyResponse(req.GetCtx(), req.Response, resp, err)
// }

func (c *ServersController) WriteProxyResponse(ctx context.Context, resp *ghttp.Response, proxyResp *resty.Response, err error) {
	if resp == nil {
		resp.WriteStatusExit(400, err)
	} else {
		glog.Infof(ctx, "openstack response status: %d, content-length: %s", proxyResp.StatusCode(), proxyResp.Header().Get("Content-Length"))
		resp.WriteJson(proxyResp.Body())
	}
}
