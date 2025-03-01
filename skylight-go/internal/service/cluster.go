package service

import (
	"fmt"
	"skylight/internal/model/entity"
	"skylight/internal/service/internal/dao"
)

type clusterService struct{}

var ClusterService *clusterService

func (s clusterService) GetClusters() ([]entity.Cluster, error) {
	return dao.GetClusters()
}
func (s clusterService) GetClusterByName(name string) (*entity.Cluster, error) {
	items, err := dao.GetClustersByName(name)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("cluster '%s' not found", name)
	}
	return &(items[0]), nil
}

func (s clusterService) CreatCluster(name string, authUrl string) (*entity.Cluster, error) {
	items, err := dao.GetClustersByName(name)
	if err != nil {
		return nil, err
	}
	if len(items) > 0 {
		return nil, fmt.Errorf("cluster '%s' already exists", name)
	}
	item, err := dao.CreateCluster(name, authUrl)
	if err != nil {
		return nil, fmt.Errorf("create cluster '%s' failed: %s", name, err)
	}
	return item, nil
}
func (s clusterService) DeleteCluster(id int) error {
	return dao.DeleteClusterById(id)
}

func init() {
	ClusterService = &clusterService{}
}
