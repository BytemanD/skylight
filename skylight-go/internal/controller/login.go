package controller

import (
	"encoding/json"
	"skylight/internal/model"
	"skylight/internal/service"
	"skylight/internal/service/openstack"

	"github.com/gogf/gf/v2/net/ghttp"
)

type LoginController struct{}

func (c *LoginController) Post(req *ghttp.Request) {
	if sessionId, err := req.Session.Id(); err != nil {
		req.Response.WriteStatusExit(403, HttpError{Code: 500, Message: err.Error()})
	} else {
		reqBody := struct{ Auth model.AuthInfo }{}
		if err := json.Unmarshal(req.GetBody(), &reqBody); err != nil {
			req.Response.WriteStatusExit(400, HttpError{Code: 403, Message: "invalid auth info"})
		}
		if reqBody.Auth.Cluster == "" {
			req.Response.WriteStatusExit(403, HttpError{Code: 400, Message: "bad request", Data: "cluster is empty"})
		}
		cluster, err := service.GetClusterByName(reqBody.Auth.Cluster)
		if err != nil {
			req.Response.WriteStatusExit(403, HttpError{Code: 400, Message: "bad request", Data: err.Error()})
		}
		if _, err := openstack.NewManager(sessionId, cluster.AuthUrl,
			reqBody.Auth.Project, reqBody.Auth.User, reqBody.Auth.Password); err != nil {
			req.Response.WriteStatusExit(403, HttpError{Code: 400, Message: "bad request", Data: err.Error()})
		} else {
			req.Session.Set("authUrl", cluster.AuthUrl)
			req.Session.Set("project", reqBody.Auth.Project)
			req.Session.Set("user", reqBody.Auth.User)
			req.Session.Set("password", reqBody.Auth.Password)
		}
	}
	req.Response.WriteStatusExit(200, HttpError{Code: 200, Message: "login success"})
}
func (c *LoginController) Get(req *ghttp.Request) {
	user, _ := req.Session.Get("user", nil)
	project, _ := req.Session.Get("project", "")
	authInfo := struct{ Auth model.AuthInfo }{
		Auth: model.AuthInfo{
			Project: project.String(),
			User:    user.String(),
		},
	}
	req.Response.WriteStatusExit(202, authInfo)
}
func (c *LoginController) Delete(req *ghttp.Request) {
	if err := req.Session.RemoveAll(); err != nil {
		req.Response.WriteStatusExit(400, HttpError{Code: 400, Message: "logout failed"})
	}
	req.Response.WriteStatusExit(200, HttpError{Code: 200, Message: "logout success"})
}
