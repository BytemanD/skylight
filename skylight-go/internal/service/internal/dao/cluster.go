package dao

import (
	"fmt"
	"skylight/internal/model"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

func queryCluster(db gdb.DB) *gdb.Model {
	return db.Model(model.Clusters{})
}

func GetClusters(db gdb.DB) ([]model.Cluster, error) {
	clusters := model.Clusters{}
	err := queryCluster(db).Scan(&clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}
func GetClustersByName(db gdb.DB, name string) ([]model.Cluster, error) {
	clusters := model.Clusters{}
	err := queryCluster(db).Where("name = ?", name).Scan(&clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}
func CreateCluster(db gdb.DB, name, authUrl string) (*model.Cluster, error) {
	cluster := model.Cluster{Name: name, AuthUrl: authUrl}
	if !strings.HasPrefix(authUrl, "https://") && !strings.HasPrefix(authUrl, "http://") {
		return nil, fmt.Errorf("invalid auth url: %s, it must starts with https:// or http://", authUrl)
	}
	if _, err := queryCluster(db).Insert(cluster); err != nil {
		return nil, err
	}
	return &cluster, nil
}
func DeleteClusterById(db gdb.DB, id int) error {
	_, err := queryCluster(db).Delete("id = ?", id)
	return err
}
