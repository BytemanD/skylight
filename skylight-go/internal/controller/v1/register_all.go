package v1

import "github.com/gogf/gf/v2/net/ghttp"

func RegiterRouters(s *ghttp.Server) {
	s.BindObjectRest("/v1/servers", ServersController{})
}
