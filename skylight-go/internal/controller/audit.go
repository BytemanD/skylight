package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"skylight/internal/model/entity"
	"skylight/internal/service"
)

type AuditsController struct{}

func (c *AuditsController) Get(req *ghttp.Request) {
	loginInfo, err := service.OSService.GetLogInfo(req)
	if err != nil {
		req.Response.WriteStatusExit(400, HttpError{Error: err.Error()})
	}
	audits, err := service.AuditService.GetByProjectId(loginInfo.Project.Id)
	if err != nil {
		req.Response.WriteStatusExit(400, HttpError{Error: err.Error()})
	}
	respBody := struct {
		Audits []entity.Audit `json:"audits"`
	}{Audits: audits}
	req.Response.WriteStatusExit(200, respBody)
}
