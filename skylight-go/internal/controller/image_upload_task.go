package controller

import (
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"skylight/internal/model/entity"
	"skylight/internal/service"
	"skylight/internal/service/openstack"
)

type ImageUploadTasksController struct{}

func (c *ImageUploadTasksController) Get(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")
	loginInfo, err := service.OSService.GetLogInfo(req)
	if err != nil {
		req.Response.WriteStatusExit(400, HttpError{Error: err.Error()})
		return
	}
	tasks, err := openstack.ImageUploadTaskService.GetByProjectId(loginInfo.Project.Id)
	if err != nil {
		req.Response.WriteStatusExit(400, HttpError{Error: err.Error()})
		return
	}
	respBody := struct {
		ImageUploadTasks []entity.ImageUploadTask `json:"image_upload_tasks"`
	}{ImageUploadTasks: tasks}
	req.Response.WriteStatusExit(200, respBody)
}

type ImageUploadTaskController struct{}

func (c *ImageUploadTaskController) Delete(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")

	routerId := req.GetRouterMap()["id"]
	if routerId == "" {
		req.Response.WriteStatusExit(400, HttpError{Message: "invalid request"})
	}
	id, err := strconv.Atoi(routerId)
	if err != nil {
		req.Response.WriteStatusExit(400, HttpError{Message: "invalid cluster id"})
	}
	if openstack.ImageUploadTaskService.Delete(id) != nil {
		g.Log().Errorf(req.GetCtx(), "delete image upload task failed: %s", err)
		req.Response.WriteStatusExit(403, HttpError{Error: "delete image upload task failed"})
	}
	req.Response.WriteStatusExit(204)
}
