package controller

import (
	"encoding/json"
	"skylight/internal/model"
	"skylight/internal/service"
	"skylight/internal/service/openstack"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type LoginController struct{}
type PostLoginController struct{}

func (c *PostLoginController) Post(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")
	sessionId, err := req.Session.Id()
	if err != nil || sessionId == "" {
		req.Response.WriteStatusExit(500, NewHttpIntervalError())
	}
	reqBody := struct{ Auth model.AuthInfo }{}
	if err := json.Unmarshal(req.GetBody(), &reqBody); err != nil {
		req.Response.WriteStatusExit(403, HttpError{Error: "invalid auth info"})
	}
	if reqBody.Auth.Cluster == "" {
		req.Response.WriteStatusExit(400, HttpError{Error: "cluster is empty"})
	}
	cluster, err := service.ClusterService.GetClusterByName(reqBody.Auth.Cluster)
	if err != nil {
		req.Response.WriteStatusExit(403, HttpError{Error: err.Error()})
	}
	if manager, err := openstack.NewManager(cluster.AuthUrl,
		reqBody.Auth.Project, reqBody.Auth.User, reqBody.Auth.Password); err != nil {
		req.Response.WriteStatusExit(403, HttpError{Error: err.Error()})
	} else {
		loginInfo := openstack.LoginInfo{
			Cluster:  cluster.Name,
			Region:   reqBody.Auth.Region,
			Project:  manager.GetProject(),
			User:     manager.GetUser(),
			Roles:    manager.GetRoles(),
			Password: reqBody.Auth.Password,
		}
		req.Session.Set("loginInfo", loginInfo)
		if err := service.AuditService.Login(req); err != nil {
			g.Log().Errorf(req.GetCtx(), "login failed: %s", err)
			req.Response.WriteStatusExit(400, HttpError{Message: "login failed"})
		}
		regions, err := manager.GetRegionFromCatalog()
		if err != nil {
			req.Response.WriteStatusExit(400, HttpError{Message: "get regions failed"})
		} else {
			g.Log().Infof(req.GetCtx(), "login success")
			// v1.AddAuthSession(
			// 	req.GetSessionId(), cluster.AuthUrl,
			// 	reqBody.Auth.User, reqBody.Auth.Project,
			// 	reqBody.Auth.Password, "RegionOne",
			// )
			req.Response.WriteStatusExit(
				200, map[string]interface{}{"regions": regions},
			)
		}
	}
}
func (c *LoginController) Get(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")

	sessionLoginInfo, _ := req.Session.Get("loginInfo", nil)
	loginInfo := openstack.LoginInfo{}

	if err := sessionLoginInfo.Struct(&loginInfo); err != nil {
		req.Response.WriteStatusExit(500, HttpError{Error: "get login info failed"})
	}
	loginInfo.Password = ""
	authInfo := struct {
		Auth openstack.LoginInfo `json:"auth"`
	}{
		Auth: loginInfo,
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
	err := service.OSService.SetRegion(req, reqBody.Auth.Region)
	if err != nil {
		req.Response.WriteStatusExit(500, HttpError{Error: "set region failed", Message: err.Error()})
		return
	}

	// session := v1.GetAuthSession(req)
	// session.SetRegion(reqBody.Auth.Region)

	req.Response.WriteStatusExit(200, HttpError{Message: "update success"})
}
func (c *LoginController) Delete(req *ghttp.Request) {
	defer service.SseService.Unregister(req.GetSessionId())

	req.Response.Header().Set("Content-Type", "application/json")
	if err := service.AuditService.Logout(req); err != nil {
		g.Log().Errorf(req.Context(), "logout failed: %s", err)
		req.Response.WriteStatusExit(400, HttpError{Message: "logout failed"})
	} else {
		req.Response.WriteStatusExit(200, HttpError{Message: "logout success"})
	}
}
