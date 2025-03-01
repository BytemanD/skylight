package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"skylight/internal/consts"
	"skylight/internal/controller"
	_ "skylight/internal/packed"
	"skylight/internal/service"

	"github.com/BytemanD/easygo/pkg/global/logging"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/os/gsession"
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
			debug := !parser.GetOpt("debug").IsNil()

			// 初始化日志
			enableDebug, logLevel := false, logging.INFO
			if debug {
				enableDebug, logLevel = true, logging.DEBUG
			}
			glog.SetDebug(enableDebug)
			glog.SetStack(enableDebug)
			logging.BasicConfig(logging.LogConfig{Level: logLevel, EnableColor: true})

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
				return errors.New("config is not available")
			}
			// 初始化DB
			// if err := service.DBInit(ctx); err != nil {
			// 	return errors.Join(errors.New("init db error"), err)
			// }
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

			s.SetSessionStorage(gsession.NewStorageFile(gsessionPath.String(), time.Hour*5))
			s.SetSessionCookieMaxAge(time.Hour)
			s.SetSessionMaxAge(time.Hour)

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

			s.BindHandler("/ws", func(req *ghttp.Request) {
				service.PublishService.RegisterPublisher(req)
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
			// TODO: 周期任务
			filepath.Walk(filepath.Join(gsessionPath.String()),
				func(path string, info fs.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}
					if !info.ModTime().Before(time.Now().AddDate(0, 0, -1)) {
						return nil
					}
					glog.Infof(ctx, "session %s is expired, cleanup", info.Name())
					os.Remove(filepath.Join(gsessionPath.String(), info.Name()))
					return nil
				},
			)

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
