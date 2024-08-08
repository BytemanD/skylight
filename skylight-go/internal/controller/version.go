package controller

import (
	"skylight/internal/consts"

	"github.com/gogf/gf/v2/net/ghttp"
)

type Version struct{}

func (c *Version) Get(req *ghttp.Request) {
	respBody := struct {
		Version map[string]string `json:"version"`
	}{
		Version: map[string]string{
			"version":       consts.Version,
			"goVersion":     consts.GoVersion,
			"buildDate":     consts.BuildDate,
			"buildPlatform": consts.BuildPlatform,
		},
	}
	req.Response.WriteJsonExit(respBody)
}
