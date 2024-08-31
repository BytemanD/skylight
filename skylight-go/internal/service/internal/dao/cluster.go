package dao

import (
	"fmt"
	"skylight/internal/service/internal/do"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func queryCluster() *gdb.Model {
	return g.DB().Model(do.Clusters{})
}

func GetClusters() ([]do.Cluster, error) {
	clusters := do.Clusters{}
	err := g.DB().Model(do.Clusters{}).Scan(&clusters)
	// err := queryCluster(db).Scan(&clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}
func GetClustersByName(name string) ([]do.Cluster, error) {
	clusters := do.Clusters{}
	err := queryCluster().Where("name = ?", name).Scan(&clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}
func CreateCluster(name, authUrl string) (*do.Cluster, error) {
	cluster := do.Cluster{Name: name, AuthUrl: authUrl}
	if !strings.HasPrefix(authUrl, "https://") && !strings.HasPrefix(authUrl, "http://") {
		return nil, fmt.Errorf("invalid auth url: %s, it must starts with https:// or http://", authUrl)
	}
	if result, err := queryCluster().Insert(cluster); err != nil {
		return nil, err
	} else {
		id, _ := result.LastInsertId()
		cluster.Id = int(id)
	}
	return &cluster, nil
}
func DeleteClusterById(id int) error {
	_, err := queryCluster().Delete("id = ?", id)
	return err
}
