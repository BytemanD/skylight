package controller

import (
	"encoding/json"
	"skylight/internal/model"
	"skylight/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

type ClusterController struct{}

func (c *ClusterController) Get(req *ghttp.Request) {
	clusters, err := service.GetClusters()
	req.Response.Header().Set("Content-Type", "application/json")
	if err != nil {
		req.Response.WriteStatusExit(400, HttpError{Code: 400, Message: err.Error()})
	}
	respBody := struct {
		Clusters model.Clusters `json:"clusters"`
	}{Clusters: clusters}
	req.Response.WriteStatusExit(200, respBody)
}

func (c *ClusterController) Post(req *ghttp.Request) {
	req.Response.Header().Set("Content-Type", "application/json")
	reqBody := struct{ Cluster model.Cluster }{}
	if err := json.Unmarshal(req.GetBody(), &reqBody); err != nil {
		req.Response.WriteStatusExit(400, HttpError{Code: 400, Message: "invalid cluster"})
	}
	cluster, err := service.CreatCluster(reqBody.Cluster.Name, reqBody.Cluster.AuthUrl)
	if err != nil {
		req.Response.WriteStatusExit(403, HttpError{Code: 400, Message: "create cluster failed", Data: err.Error()})
	}
	req.Response.WriteStatusExit(200, cluster)
}
