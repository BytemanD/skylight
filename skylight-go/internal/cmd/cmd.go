package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"skylight/internal/controller"
	"skylight/internal/service"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gsession"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Version       string
	GoVersion     string
	BuildDate     string
	BuildPlatform string
)

var PROXY_PREFIXY = []string{
	"/identity",
	"/networking", "/computing", "/volume", "/image",
}

var (
	VersionCmd = gcmd.Command{
		Name:  "version",
		Usage: "version",
		Brief: "show version",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			fmt.Println("Version: ", Version)
			fmt.Println("GoVersion: ", GoVersion)
			fmt.Println("BuildDate: ", BuildDate)
			fmt.Println("BuildPlatform: ", BuildPlatform)
			return nil
		},
	}

	ServeCmd = gcmd.Command{
		Name:  "serve",
		Usage: "serve",
		Brief: "start http server",
		Arguments: []gcmd.Argument{
			{Name: "port", Short: "p", Brief: "The port of server"},
			{Name: "debug", Short: "d", Orphan: true, Brief: "Show debug message"},
			{Name: "static", Short: "d", Orphan: true, Brief: "The path of static"},
			{Name: "template", Short: "d", Orphan: true, Brief: "The path of template"},
		},
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 初始化日志
			port := parser.GetOpt("port", "8081").String()
			debug := parser.ContainsOpt("debug")
			templatePath := parser.GetOpt("template", "skylight-web/dist")
			staticPath := parser.GetOpt("static", "skylight-web/dist/static")

			level := logging.INFO
			if debug {
				level = logging.DEBUG
			}
			logging.BasicConfig(logging.LogConfig{Level: level, EnableColor: true})

			s := g.Server()
			// 初始化DB
			logging.Info("init db driver ...")
			dbDriver, _ := g.Cfg().Get(ctx, "database.type", "sqlite")
			dbLink, _ := g.Cfg().Get(ctx, "database.link", "/var/lib/skylight/skylight.db")
			logging.Info("dadtabase link: %s", dbLink.String())
			service.DBInit(ctx, dbDriver.String(), dbLink.String())

			if gfile.Exists(staticPath.String()) {
				s.AddStaticPath("/static", staticPath.String())
			}
			if gfile.Exists(templatePath.String()) {
				s.AddSearchPath(templatePath.String())
			}
			if port != "" {
				s.SetAddr(fmt.Sprintf(":%s", port))
			}

			// 初始化 session 驱动
			logging.Info("init session driver ...")
			sessionDriver, _ := g.Cfg().Get(ctx, "session.type", "file")
			sessionPath, _ := g.Cfg().Get(ctx, "session.path", "/var/lib/skylight")
			gsessionPath := filepath.Join(sessionPath.String(), "gsessions")
			if !gfile.Exists(gsessionPath) {
				logging.Info("create dir '%s'", gsessionPath)
				if err := gfile.Mkdir(gsessionPath); err != nil {
					return fmt.Errorf("create dir '%s' failed: %s", gsessionPath, err)
				}
			}

			if sessionDriver.String() != "file" {
				return fmt.Errorf("invalid session driver: %s", sessionDriver.String())
			}
			s.SetSessionStorage(gsession.NewStorageFile(gsessionPath))

			s.BindMiddlewareDefault(
				controller.MiddlewareCORS,
				ghttp.MiddlewareHandlerResponse, controller.MiddlewareLogResponse,
				controller.MiddlewareAuth,
			)
			// 注册路由
			s.BindObjectRest("/login", controller.LoginController{})
			s.BindObjectRest("/clusters", controller.ClusterController{})

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
	Main = gcmd.Command{}
)

func init() {
	Main.AddCommand(&ServeCmd, &VersionCmd)
}
