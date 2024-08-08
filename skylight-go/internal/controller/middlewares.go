package controller

import (
	"time"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/gogf/gf/v2/net/ghttp"
)

func MiddlewareCORS(req *ghttp.Request) {
	req.Response.CORSDefault()
	req.Response.Header().Set("Access-Control-Expose-Headers", "X-Auth-Token")
	req.Middleware.Next()
}

func MiddlewareLogResponse(req *ghttp.Request) {
	startTime := time.Now()
	req.Middleware.Next()
	spentTime := time.Since(startTime).Seconds()
	if req.Response.Status < 400 {
		logging.Info("%s %s -> [%d] (%fs)", req.Method, req.URL, req.Response.Status, spentTime)
	} else {
		logging.Error("%s %s -> [%d] (%fs)\n    Resp: %s",
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
		logging.Error("get session id failed: %s", err)
		req.Response.WriteStatusExit(400, HttpError{Code: 500, Message: "internal error"})
	}
	if user, err := req.Session.Get("user", nil); err != nil || user.IsNil() {
		logging.Error("invalid request: auth info not found in session")
		req.Response.WriteStatusExit(403, HttpError{Code: 403, Message: "no login"})
	}
	req.Middleware.Next()
}
