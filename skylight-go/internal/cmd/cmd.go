package cmd

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"skylight/internal/consts"
	"skylight/internal/controller"
	_ "skylight/internal/packed"
	"skylight/internal/service"
	"skylight/utility"

	"github.com/BytemanD/easygo/pkg/global/logging"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
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
		Brief: "start skylight server",
		Arguments: []gcmd.Argument{
			{Name: "port", Short: "p", Brief: "The port of server"},
			{Name: "debug", Short: "d", Orphan: true, Brief: "Show debug message"},
			{Name: "static", Short: "S", Brief: "The path of static"},
			{Name: "template", Short: "T", Brief: "The path of template"},
		},
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			port := parser.GetOpt("port")
			debug := !parser.GetOpt("debug").IsNil()

			// 初始化日志
			enableDebug, logLevel := false, logging.INFO
			if debug {
				enableDebug, logLevel = debug, logging.DEBUG
			}
			g.Log().SetDebug(enableDebug)
			g.Log().SetStack(enableDebug)
			logging.BasicConfig(logging.LogConfig{Level: logLevel, EnableColor: true})

			s := g.Server()
			if !port.IsNil() && port.Uint() > 0 {
				g.Log().Debugf(ctx, "use port %d", port.Uint())
				s.SetAddr(fmt.Sprintf(":%d", port.Uint()))
			}
			// 初始化静态资源
			if gres.IsEmpty() {
				if gfile.Exists(DEV_TEMPLATE) {
					g.Log().Infof(ctx, "add template path: %s", DEV_TEMPLATE)
					s.AddSearchPath(DEV_TEMPLATE)
				} else {
					g.Log().Warningf(ctx, "template %s not exists", DEV_TEMPLATE)
				}
				if gfile.Exists(DEV_STATIC) {
					g.Log().Infof(ctx, "static path: %s", DEV_STATIC)
					s.AddStaticPath("/static", DEV_STATIC)
				} else {
					g.Log().Warningf(ctx, "static %s not exists", DEV_STATIC)
				}
			}

			if !g.Cfg().Available(ctx) {
				return errors.New("config is not available")
			}
			// 初始化DB
			// if err := service.DBInit(ctx); err != nil {
			// 	return errors.Join(errors.New("init db error"), err)
			// }
			if gres.Contains("resources") {
				s.AddSearchPath("resources")
			}

			dataPath := g.Cfg().MustGet(ctx, "server.dataPath")
			g.Log().Infof(ctx, "data path: %s", dataPath.String())
			utility.MakesureDir(dataPath.String())
			utility.MakesureDir(filepath.Join(dataPath.String(), "image_cache"))

			// 初始化 session 路径
			service.InitSessionStorage(ctx)
			s.SetSessionStorage(service.SessionStorage)
			// s.SetSessionStorage(gsession.NewStorageFile(gsessionPath.String(), time.Hour*5))
			s.SetSessionCookieMaxAge(time.Hour)
			s.SetSessionMaxAge(time.Hour)

			InitDB(ctx)
			// core := g.DB().GetCore()
			s.BindMiddlewareDefault(
				controller.MiddlewareCORS,
				ghttp.MiddlewareHandlerResponse, controller.MiddlewareLogResponse,
			)
			// 注册路由
			s.BindObjectRest("/login", controller.PostLoginController{})
			s.BindObjectRest("/version", controller.Version{})
			s.BindObjectRest("/clusters", controller.ClustersController{})
			s.BindObjectRest("/clusters/:id", controller.ClusterController{})
			s.BindObjectRest("/image_upload_tasks", controller.ImageUploadTasksController{})
			s.BindObjectRest("/image_upload_tasks/:id", controller.ImageUploadTaskController{})

			// v1.RegiterRouters(s)

			s.Group("", func(group *ghttp.RouterGroup) {
				group.Middleware(controller.MiddlewareAuth)

				s.BindObjectRest("/audits", controller.AuditsController{})
				group.REST("/login", controller.LoginController{})
				for _, prefix := range PROXY_PREFIXY {
					if prefix == "/image" {
						group.Middleware(controller.MiddlewareGlanceImageUploadCache)
					}
					group.REST(prefix+"/*", controller.OpenstackProxy{Prefix: prefix})
				}
			})

			// SSE
			s.BindHandler("/sse", func(req *ghttp.Request) {
				session := req.GetQuery("session", "")
				if session.IsEmpty() {
					req.Response.WriteStatusExit(400, controller.HttpError{
						Error: "session is missing",
					})
				}
				if !service.OSService.IsLogin(session.String()) {
					req.Response.WriteStatusExit(403, controller.HttpError{
						Error: "session is not login",
					})
				}
				service.SseService.Register(session.String(), req)
			})

			g.Log().Info(ctx, "starting server")
			s.Run()
			return nil
		},
	}
	Main = gcmd.Command{}
)

func init() {
	Main.AddCommand(&ServeCmd, &VersionCmd)
}
