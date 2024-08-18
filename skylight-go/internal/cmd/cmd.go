package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"skylight/internal/consts"
	"skylight/internal/controller"
	_ "skylight/internal/packed"
	"skylight/internal/service"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/os/gsession"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DEV_TEMPLATE = "../skylight-web/dist"
	DEV_STATIC   = "../skylight-web/dist/static"
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
			fmt.Println("Version: ", consts.Version)
			fmt.Println("GoVersion: ", consts.GoVersion)
			fmt.Println("BuildDate: ", consts.BuildDate)
			fmt.Println("BuildPlatform: ", consts.BuildPlatform)
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
			{Name: "static", Short: "S", Brief: "The path of static"},
			{Name: "template", Short: "T", Brief: "The path of template"},
		},
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			port := parser.GetOpt("port", "8081").String()
			debug := parser.ContainsOpt("debug")

			// 初始化日志
			level := logging.INFO
			if debug {
				level = logging.DEBUG
			} else {
				glog.SetDebug(true)
			}
			logging.BasicConfig(logging.LogConfig{Level: level, EnableColor: true})
			glog.SetStack(false)

			s := g.Server()
			// 初始化静态资源
			if gres.IsEmpty() {
				if gfile.Exists(DEV_TEMPLATE) {
					glog.Infof(ctx, "add template path: %s", DEV_TEMPLATE)
					s.AddSearchPath(DEV_TEMPLATE)
				} else {
					glog.Warningf(ctx, "template %s not exists", DEV_TEMPLATE)
				}
				if gfile.Exists(DEV_STATIC) {
					glog.Infof(ctx, "static path: %s", DEV_STATIC)
					s.AddStaticPath("/static", DEV_STATIC)
				} else {
					glog.Warningf(ctx, "static %s not exists", DEV_STATIC)
				}
			}

			if !g.Cfg().Available(ctx) {
				return fmt.Errorf("config is not available")
			}

			// 初始化DB
			glog.Infof(ctx, "init db driver ...")
			dbDriver, _ := g.Cfg().Get(ctx, "database.type", "sqlite")
			dbLink, _ := g.Cfg().Get(ctx, "database.link", "/var/lib/skylight/skylight.db")
			glog.Infof(ctx, "dadtabase link: %s", dbLink.String())
			service.DBInit(ctx, dbDriver.String(), dbLink.String())

			if gres.Contains("resources") {
				s.AddSearchPath("resources")
			}

			if port != "" {
				s.SetAddr(fmt.Sprintf(":%s", port))
			}
			// 初始化 session 驱动
			glog.Infof(ctx, "init session driver ...")
			sessionPath, _ := g.Cfg().Get(ctx, "session.path", "/var/lib/skylight")
			gsessionPath := filepath.Join(sessionPath.String(), "gsessions")
			if !gfile.Exists(gsessionPath) {
				glog.Infof(ctx, "create dir '%s'", gsessionPath)
				if err := gfile.Mkdir(gsessionPath); err != nil {
					return fmt.Errorf("create dir '%s' failed: %s", gsessionPath, err)
				}
			}
			glog.Infof(ctx, "session path: %s", gsessionPath)
			s.SetSessionStorage(gsession.NewStorageFile(gsessionPath))

			s.BindMiddlewareDefault(
				controller.MiddlewareCORS,
				ghttp.MiddlewareHandlerResponse, controller.MiddlewareLogResponse,
			)
			// 注册路由
			s.BindObjectRest("/version", controller.Version{})
			s.BindObjectRest("/clusters", controller.ClustersController{})
			s.BindObjectRest("/clusters/:id", controller.ClusterController{})
			s.BindObjectRest("/login", controller.PostLoginController{})
			s.Group("", func(group *ghttp.RouterGroup) {
				group.Middleware(controller.MiddlewareAuth)
				group.REST("/login", controller.LoginController{})
				for _, prefix := range PROXY_PREFIXY {
					group.REST(prefix+"/*", controller.OpenstackProxy{Prefix: prefix})
				}
			})
			glog.Info(ctx, "starting server")
			s.Run()
			return nil
		},
	}
	Main = gcmd.Command{}
)

func init() {
	Main.AddCommand(&ServeCmd, &VersionCmd)
}
