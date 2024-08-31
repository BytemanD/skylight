package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

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
			service.DBInit(ctx)

			if gres.Contains("resources") {
				s.AddSearchPath("resources")
			}

			if port != "" {
				s.SetAddr(fmt.Sprintf(":%s", port))
			}
			dataPath, _ := g.Cfg().Get(ctx, "server.dataPath", "/var/lib/skylight")
			glog.Infof(ctx, "data path: %s", dataPath.String())
			MakesureDir(dataPath.String())
			MakesureDir(filepath.Join(dataPath.String(), "image_cache"))

			// 初始化 session 驱动
			glog.Infof(ctx, "init session driver ...")
			gsessionPath, _ := g.Cfg().Get(ctx, "session.path", "/var/lib/skylight/gsessions")
			glog.Infof(ctx, "session path: %s", gsessionPath.String())
			MakesureDir(gsessionPath.String())

			s.SetSessionStorage(gsession.NewStorageFile(gsessionPath.String()))
			s.SetSessionCookieMaxAge(time.Hour)
			s.SetSessionMaxAge(time.Hour)

			s.BindMiddlewareDefault(
				controller.MiddlewareCORS,
				ghttp.MiddlewareHandlerResponse, controller.MiddlewareLogResponse,
			)
			// 注册路由
			s.BindObjectRest("/version", controller.Version{})
			s.BindObjectRest("/clusters", controller.ClustersController{})
			s.BindObjectRest("/clusters/:id", controller.ClusterController{})
			s.BindObjectRest("/login", controller.PostLoginController{})
			s.BindObjectRest("/image_upload_tasks", controller.ImageUploadTasksController{})
			s.BindObjectRest("/image_upload_tasks/:id", controller.ImageUploadTaskController{})
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
