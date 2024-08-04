package controller

import (
	"strings"
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
func mustAuth(req *ghttp.Request) bool {
	if req.Request.URL.Path == "/clusters" || req.Request.URL.Path == "/version" {
		return false
	}
	if req.Request.URL.Path == "/favicon.ico" {
		return false
	}
	if req.Request.URL.Path == "/login" && strings.ToUpper(req.Method) == "POST" {
		return false
	}
	return true
}
func MiddlewareAuth(req *ghttp.Request) {
	if mustAuth(req) {
		if _, err := req.Session.Id(); err != nil {
			logging.Error("get session id failed: %s", err)
			req.Response.WriteStatusExit(400, HttpError{Code: 500, Message: "internal error"})
		}
		if user, err := req.Session.Get("user", nil); err != nil || user.IsNil() {
			logging.Error("invalid request: auth info not found in session")
			req.Response.WriteStatusExit(403, HttpError{Code: 403, Message: "no login"})
		}
	}

	req.Middleware.Next()
}
