package controller

import (
	"encoding/json"
	"strconv"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"

	"skylight/internal/model/entity"
	"skylight/internal/service"
)

type ClustersController struct{}

func (c *ClustersController) Get(req *ghttp.Request) {
	clusters, err := service.ClusterService.GetClusters()
	req.Response.Header().Set("Content-Type", "application/json")
	if err != nil {
		req.Response.WriteStatusExit(400, HttpError{Error: err.Error()})
	}
	respBody := struct {
		Clusters []entity.Cluster `json:"clusters"`
	}{Clusters: clusters}
	req.Response.WriteStatusExit(200, respBody)
}

func (c *ClustersController) Post(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")
	reqBody := struct{ Cluster entity.Cluster }{}
	if err := json.Unmarshal(req.GetBody(), &reqBody); err != nil {
		req.Response.WriteStatusExit(400, HttpError{Error: "invalid cluster"})
	}
	cluster, err := service.ClusterService.CreatCluster(reqBody.Cluster.Name, reqBody.Cluster.AuthUrl)
	if err != nil {
		glog.Errorf(req.GetCtx(), "create cluster failed: %s", err)
		req.Response.WriteStatusExit(400, HttpError{Error: "create cluster failed"})
	}
	req.Response.WriteStatusExit(200, cluster)
}

type ClusterController struct{}

func (c *ClusterController) Delete(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")

	routerId := req.GetRouterMap()["id"]
	if routerId == "" {
		req.Response.WriteStatusExit(400, HttpError{Message: "invalid request"})
	}
	id, err := strconv.Atoi(routerId)
	if err != nil {
		req.Response.WriteStatusExit(400, HttpError{Message: "invalid cluster id"})
	}
	if service.ClusterService.DeleteCluster(id) != nil {
		glog.Errorf(req.GetCtx(), "delete cluster failed: %s", err)
		req.Response.WriteStatusExit(403, HttpError{Error: "delete cluster failed"})
	}
	req.Response.WriteStatusExit(204)
}
