package dao

import (
	"fmt"
	"skylight/internal/model/entity"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

const TABLE_CLUSTER = "clusters"

func modelCluster() *gdb.Model {
	return g.Model(TABLE_CLUSTER)
}

func GetClusters() ([]entity.Cluster, error) {
	clusters := []entity.Cluster{}
	err := modelCluster().Scan(&clusters)
	// err := queryCluster(db).Scan(&clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}
func GetClustersByName(name string) ([]entity.Cluster, error) {
	clusters := []entity.Cluster{}
	err := modelCluster().Where("name = ?", name).Scan(&clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}
func CreateCluster(name, authUrl string) (*entity.Cluster, error) {
	cluster := entity.Cluster{Name: name, AuthUrl: authUrl}
	if !strings.HasPrefix(authUrl, "https://") && !strings.HasPrefix(authUrl, "http://") {
		return nil, fmt.Errorf("invalid auth url: %s, it must starts with https:// or http://", authUrl)
	}
	if result, err := modelCluster().Insert(cluster); err != nil {
		return nil, err
	} else {
		id, _ := result.LastInsertId()
		cluster.Id = int(id)
	}
	return &cluster, nil
}
func DeleteClusterById(id int) error {
	_, err := modelCluster().Delete("id = ?", id)
	return err
}
