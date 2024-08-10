package controller

import (
	"encoding/json"
	"skylight/internal/model"
	"skylight/internal/service"
	"skylight/internal/service/openstack"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

type LoginController struct{}
type PostLoginController struct{}

func (c *PostLoginController) Post(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")
	sessionId := req.GetSessionId()
	if sessionId == "" {
		req.Response.WriteStatusExit(500, NewHttpIntervalError())
	}
	reqBody := struct{ Auth model.AuthInfo }{}
	if err := json.Unmarshal(req.GetBody(), &reqBody); err != nil {
		req.Response.WriteStatusExit(403, HttpError{Error: "invalid auth info"})
	}
	if reqBody.Auth.Cluster == "" {
		req.Response.WriteStatusExit(400, HttpError{Error: "cluster is empty"})
	}
	cluster, err := service.GetClusterByName(reqBody.Auth.Cluster)
	if err != nil {
		req.Response.WriteStatusExit(403, HttpError{Error: err.Error()})
	}
	if _, err := openstack.NewManager(sessionId, cluster.AuthUrl,
		reqBody.Auth.Project, reqBody.Auth.User, reqBody.Auth.Password); err != nil {
		req.Response.WriteStatusExit(403, HttpError{Error: err.Error()})
	} else {
		req.Session.Set("cluster", cluster.Name)
		req.Session.Set("authUrl", cluster.AuthUrl)
		req.Session.Set("region", reqBody.Auth.Region)
		req.Session.Set("project", reqBody.Auth.Project)
		req.Session.Set("user", reqBody.Auth.User)
		req.Session.Set("password", reqBody.Auth.Password)
	}

	glog.Infof(req.GetCtx(), "login success")
	req.Response.Header().Set("Session-Id", sessionId)
	req.Response.WriteStatusExit(200, HttpError{Message: "login success"})
}
func (c *LoginController) Get(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")
	user, _ := req.Session.Get("user", nil)
	project, _ := req.Session.Get("project", "")
	cluster, _ := req.Session.Get("cluster", "")
	region, _ := req.Session.Get("region", "")
	authInfo := struct {
		Auth model.AuthInfo `json:"auth"`
	}{
		Auth: model.AuthInfo{
			Cluster: cluster.String(),
			Region:  region.String(),
			Project: project.String(),
			User:    user.String(),
		},
	}
	req.Response.WriteJson(authInfo)
}
func (c *LoginController) Put(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")
	reqBody := struct{ Auth model.AuthInfo }{}
	if err := json.Unmarshal(req.GetBody(), &reqBody); err != nil {
		req.Response.WriteStatusExit(400, HttpError{Error: "invalid body"})
	}
	if reqBody.Auth.Region == "" {
		req.Response.WriteStatusExit(400, HttpError{Error: "region is empty"})
	}
	sessionId, err := req.Session.Id()
	if err != nil {
		req.Response.WriteStatusExit(500, HttpError{Error: err.Error()})
	}
	manager, err := openstack.GetManager(sessionId, req)
	if err != nil {
		req.Response.WriteStatusExit(500, HttpError{Error: "get manager failed", Message: err.Error()})
		return
	}
	req.Session.Set("region", reqBody.Auth.Region)
	manager.SetRegion(reqBody.Auth.Region)
	req.Response.WriteStatusExit(200, HttpError{Message: "update success"})
}
func (c *LoginController) Delete(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")
	if err := req.Session.RemoveAll(); err != nil {
		req.Response.WriteStatusExit(400, HttpError{Error: "logout failed"})
	}
	req.Response.WriteStatusExit(200, HttpError{Message: "logout success"})
}
