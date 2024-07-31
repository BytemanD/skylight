package controller

import (
	"context"
	"encoding/json"
	"skylight/apiv1"
	"skylight/internal/service/openstack"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

type BadRequest struct {
	Error string `json:"error"`
}

type Networking struct{}

func (c *Networking) Get(ctx context.Context, apiReq *apiv1.NetworkingGetReq) (res *apiv1.NetworkingGetRes, err error) {
	req := g.RequestFromCtx(ctx)
	proxyUrl := strings.TrimPrefix(req.URL.String(), "/networking")
	manager := openstack.GetManager()
	resp, err := manager.ProxyNetworking(proxyUrl)

	if err != nil {
		data, _ := json.Marshal(BadRequest{Error: err.Error()})
		req.Response.WriteStatus(400, data)
	} else {
		req.Response.WriteStatus(resp.StatusCode(), resp.Body())
		req.Response.Header().Set("Content-Type", resp.Header().Get("Content-Type"))
	}
	return
}
