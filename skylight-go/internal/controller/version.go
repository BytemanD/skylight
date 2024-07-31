package controller

import (
	"context"
	"skylight/apiv1"

	"github.com/gogf/gf/v2/frame/g"
)

type Version struct{}

func (c *Version) Get(ctx context.Context, apiReq *apiv1.VersionGetReq) (res *apiv1.VersionGetRes, err error) {
	req := g.RequestFromCtx(ctx)
	req.Response.WriteJson(map[string]string{"version": "dev"})
	return
}
