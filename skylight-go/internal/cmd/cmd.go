package cmd

import (
	"context"
	"fmt"
	"time"

	"skylight/internal/controller"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

func MiddlewareCORS(req *ghttp.Request) {
	req.Response.CORSDefault()
	req.Response.Header().Set("Access-Control-Expose-Headers", "X-Auth-Token")
	req.Middleware.Next()
}
func MiddlewareLogResponse(r *ghttp.Request) {
	startTime := time.Now()
	r.Middleware.Next()
	spentTime := time.Since(startTime).Seconds()
	if r.Response.Status < 400 {
		logging.Info("%s %s -> [%d] (%fs)", r.Method, r.URL, r.Response.Status, spentTime)
	} else {
		logging.Error("%s %s -> [%d] (%fs)\n    Resp: %s",
			r.Method, r.URL, r.Response.Status, spentTime,
			r.Response.BufferString())
	}
}
func MiddlewareAuth(req *ghttp.Request) {
	token := req.Header.Get("X-Auth-Token")
	if token == "" {
		logging.Error("no auth")
		req.Response.WriteStatusExit(403, "not auth")
		return
	}
	req.Middleware.Next()
}

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			port := parser.GetOpt("port", "8091").String()
			debug := parser.ContainsOpt("debug")
			level := logging.INFO
			if debug {
				level = logging.DEBUG
			}
			logging.BasicConfig(logging.LogConfig{Level: level, EnableColor: true})

			s := g.Server()
			if port != "" {
				s.SetAddr(fmt.Sprintf(":%s", port))
			}
			s.BindMiddlewareDefault(
				MiddlewareCORS,
				ghttp.MiddlewareHandlerResponse, MiddlewareLogResponse,
			)
			// 注册路由
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Bind(new(controller.Login))
			})
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Bind(new(controller.Version))
			})
			s.BindObjectRest("/networking/*", new(controller.Networking))
			s.BindObjectRest("/computing/*", new(controller.Computing))

			logging.Info("starting server")
			s.Run()
			return nil
		},
	}
)

func init() {
	Main.Arguments = append(
		Main.Arguments,
		gcmd.Argument{Name: "port", Short: "p", Brief: "The port of server"},
		gcmd.Argument{Name: "debug", Short: "d", Orphan: true, Brief: "Show debug message"},
	)
}
