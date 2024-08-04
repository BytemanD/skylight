package cmd

import (
	"context"
	"fmt"

	"skylight/internal/controller"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var PROXY_PREFIXY = []string{"/networking", "/computing", "/volume", "/image"}

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
			// s.SetSessionStorage(gsession.NewStorageFile("/tmp"))
			if port != "" {
				s.SetAddr(fmt.Sprintf(":%s", port))
			}
			s.BindMiddlewareDefault(
				controller.MiddlewareCORS,
				ghttp.MiddlewareHandlerResponse, controller.MiddlewareLogResponse,
				controller.MiddlewareAuth,
			)
			// 注册路由
			s.BindObjectRest("/login", controller.LoginController{})

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Bind(new(controller.Version))
			})
			for _, prefix := range PROXY_PREFIXY {
				s.BindObjectRest(prefix+"/*", controller.OpenstackProxy{Prefix: prefix})
			}

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
