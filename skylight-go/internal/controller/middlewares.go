package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"skylight/internal/service"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func MiddlewareCORS(req *ghttp.Request) {
	req.Response.CORSDefault()
	req.Response.Header().Set("Access-Control-Expose-Headers", "X-Auth-Token")
	req.Middleware.Next()
}

func MiddlewareLogResponse(req *ghttp.Request) {
	g.Log().Infof(req.GetCtx(), "%s %s", req.Method, req.URL)
	startTime := time.Now()
	req.Middleware.Next()
	spentTime := time.Since(startTime).Seconds()
	if req.Response.Status < 400 {
		g.Log().Infof(req.GetCtx(), "%s %s -> [%d] (%fs)", req.Method, req.URL, req.Response.Status, spentTime)
	} else {
		g.Log().Errorf(req.GetCtx(), "%s %s -> [%d] (%fs)\n    Resp: %s",
			req.Method, req.URL, req.Response.Status, spentTime,
			req.Response.BufferString())
	}
}

type NoAuthRule struct {
	Method string
	Path   string
	Router string
}

func MiddlewareAuth(req *ghttp.Request) {
	if _, err := req.Session.Id(); err != nil {
		g.Log().Errorf(req.GetCtx(), "get session id failed: %s", err)
		req.Response.WriteStatusExit(400, NewHttpIntervalError())
	}
	if user, err := req.Session.Get("loginInfo", nil); err != nil || user.IsNil() {
		g.Log().Errorf(req.GetCtx(), "invalid request: auth info not found in session")
		req.Response.WriteStatusExit(403, HttpError{Error: "no login"})
	}
	req.Middleware.Next()
}

func MiddlewareGlanceImageUploadCache(req *ghttp.Request) {
	uploadFileReg, _ := regexp.Compile("/images/.+/file")
	proxyUrl := strings.TrimPrefix(req.URL.Path, "/image")
	if req.Method == http.MethodPut && uploadFileReg.MatchString(proxyUrl) {
		task, err := service.OSService.SaveImageCache(proxyUrl, req)
		if err != nil {
			g.Log().Errorf(req.GetCtx(), "save image cache failed: %s", err)
			req.Response.WriteStatusExit(400, HttpError{Error: "save image cache failed"})
		}
		if err != nil {
			req.Response.WriteStatusExit(400, HttpError{
				Error: fmt.Sprintf("get task for %s failed: %s", task.ImageId, err),
			})
		}
		if task.Cached < task.Size {
			req.Response.WriteStatusExit(204)
		} else {
			g.Log().Infof(req.GetCtx(), "image %s all cached", task.ImageId)
		}
	}
	req.Middleware.Next()
}
