package controller

import (
	"context"
	"encoding/json"
	"skylight/apiv1"
	"skylight/internal/service/openstack"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"
)

type Login struct{}

type Token struct {
	UUID string
}

var tokens []Token = []Token{}

func (c *Login) Post(ctx context.Context, apiReq *apiv1.LoginPostReq) (res *apiv1.LoginPostRes, err error) {
	req := g.RequestFromCtx(ctx)
	uuid := guid.S()
	req.Response.Header().Set("X-Auth-Token", uuid)
	tokens = append(tokens, Token{UUID: uuid})
	req.Response.WriteStatusExit(200)
	return
}
func (c *Login) Get(ctx context.Context, apiReq *apiv1.LoginGetReq) (res *apiv1.LoginGetRes, err error) {
	req := g.RequestFromCtx(ctx)
	token := req.Header.Get("X-Auth-Token")
	if token == "" {
		req.Response.WriteStatusExit(403, "not auth")
	}
	req.Response.WriteJson(nil)
	return
}

type LoginController struct{}

type AuthInfo struct {
	Project  string
	User     string
	Password string
}

func (c *LoginController) Post(req *ghttp.Request) {
	if sessionId, err := req.Session.Id(); err != nil {
		req.Response.WriteStatusExit(403, HttpError{Code: 500, Message: err.Error()})
	} else {
		authBody := struct{ Auth AuthInfo }{}
		if err := json.Unmarshal(req.GetBody(), &authBody); err != nil {
			req.Response.WriteStatusExit(400, HttpError{Code: 403, Message: "invalid auth info"})
		}
		if _, err := openstack.NewManager(sessionId,
			authBody.Auth.Project, authBody.Auth.User, authBody.Auth.Password); err != nil {
			req.Response.WriteStatusExit(403, HttpError{Code: 403, Message: "bad request", Data: err.Error()})
		} else {
			req.Session.Set("project", authBody.Auth.Project)
			req.Session.Set("user", authBody.Auth.User)
			req.Session.Set("password", authBody.Auth.Password)
		}
	}
	req.Response.WriteStatusExit(200, HttpError{Code: 200, Message: "login success"})
}
func (c *LoginController) Get(req *ghttp.Request) {
	user, err := req.Session.Get("user", "")
	if err != nil || user.String() == "" {
		req.Response.WriteStatusExit(403, HttpError{Code: 403, Message: "not login"})
	}
	logging.Info("user already login")

	project, _ := req.Session.Get("project", "")
	authInfo := struct{ Auth AuthInfo }{
		Auth: AuthInfo{
			Project: project.String(),
			User:    user.String(),
		},
	}
	req.Response.WriteStatusExit(202, authInfo)
}
