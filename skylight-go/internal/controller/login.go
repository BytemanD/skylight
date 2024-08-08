package controller

import (
	"encoding/json"
	"skylight/internal/model"
	"skylight/internal/service"
	"skylight/internal/service/openstack"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/gogf/gf/v2/net/ghttp"
)

type LoginController struct{}
type PostLoginController struct{}

func (c *PostLoginController) Post(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")
	sessionId, err := req.Session.Id()
	if err != nil {
		req.Response.WriteStatusExit(403, HttpError{Code: 500, Message: err.Error()})
	}
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
		req.Session.Set("cluster", cluster.Name)
		req.Session.Set("authUrl", cluster.AuthUrl)
		req.Session.Set("project", reqBody.Auth.Project)
		req.Session.Set("user", reqBody.Auth.User)
		req.Session.Set("password", reqBody.Auth.Password)
	}

	logging.Info("login success")
	req.Response.Header().Set("Session-Id", sessionId)
	req.Response.WriteStatusExit(200, HttpError{Code: 200, Message: "login success"})
}
func (c *LoginController) Get(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")
	user, _ := req.Session.Get("user", nil)
	project, _ := req.Session.Get("project", "")
	cluster, _ := req.Session.Get("cluster", "")
	authInfo := struct {
		Auth model.AuthInfo `json:"auth"`
	}{
		Auth: model.AuthInfo{
			Cluster: cluster.String(),
			Project: project.String(),
			User:    user.String(),
		},
	}
	req.Response.WriteStatusExit(202, authInfo)
}
func (c *LoginController) Put(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")
	reqBody := struct{ Auth model.AuthInfo }{}
	if err := json.Unmarshal(req.GetBody(), &reqBody); err != nil {
		req.Response.WriteStatusExit(400, HttpError{Code: 400, Message: "invalid body"})
	}
	if reqBody.Auth.Region == "" {
		req.Response.WriteStatusExit(400, HttpError{Code: 400, Message: "region is empty"})
	}
	sessionId, err := req.Session.Id()
	if err != nil {
		req.Response.WriteStatusExit(500, HttpError{Code: 500, Message: err.Error()})
	}
	manager, err := openstack.GetManager(sessionId, req)
	if err != nil {
		req.Response.WriteStatusExit(500, HttpError{Code: 500, Message: "get manager failed", Data: err.Error()})
		return
	}
	manager.SetRegion(reqBody.Auth.Region)
	req.Response.WriteStatusExit(200, HttpError{Code: 200, Message: "update success"})
}
func (c *LoginController) Delete(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")
	if err := req.Session.RemoveAll(); err != nil {
		req.Response.WriteStatusExit(400, HttpError{Code: 400, Message: "logout failed"})
	}
	req.Response.WriteStatusExit(200, HttpError{Code: 200, Message: "logout success"})
}
