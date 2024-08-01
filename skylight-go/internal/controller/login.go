package controller

import (
	"context"
	"encoding/json"
	"skylight/apiv1"

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
	authInfo := struct{ Auth AuthInfo }{}
	err := json.Unmarshal(req.GetBody(), &authInfo)
	if err != nil {
		req.Response.WriteStatusExit(400, HttpError{Code: 403, Message: "invalid auth info"})
	}
	req.Session.Set("user", authInfo.Auth.User)
	req.Session.Set("project", authInfo.Auth.Project)
	req.Session.Set("password", authInfo.Auth.Password)
	req.Response.WriteStatusExit(200)
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
