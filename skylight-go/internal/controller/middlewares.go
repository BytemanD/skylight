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
	if !(req.Request.URL.Path == "/login" && strings.ToUpper(req.Method) == "POST") {
		if _, err := req.Session.Id(); err != nil {
			logging.Error("get session id failed: %s", err)
			req.Response.WriteStatusExit(400, HttpError{Code: 500, Message: "internal error"})
		}
		if user, _ := req.Session.Get("user"); user == nil {
			logging.Error("invalid request: auth info not found in session")
			req.Response.WriteStatusExit(403, HttpError{Code: 403, Message: "no login"})
		}
	}

	req.Middleware.Next()
}
