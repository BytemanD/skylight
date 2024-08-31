package service

import (
	"context"
	"fmt"
	"path/filepath"
	"skylight/internal/model/entity"
	"skylight/internal/service/internal/dao"
	"skylight/internal/service/internal/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB gdb.DB

type clusterService struct{}

func parseCluster(item do.Cluster) entity.Cluster {
	return entity.Cluster{Id: item.Id, Name: item.Name, AuthUrl: item.AuthUrl}
}
func parseClusters(items []do.Cluster) []entity.Cluster {
	clusters := []entity.Cluster{}
	for _, item := range items {
		clusters = append(clusters, parseCluster(item))
	}
	return clusters
}

// cluster
func (s clusterService) GetClusters() ([]entity.Cluster, error) {
	items, err := dao.GetClusters()
	if err != nil {
		return nil, err
	}
	return parseClusters(items), nil
}
func (s clusterService) GetClusterByName(name string) (*entity.Cluster, error) {
	items, err := dao.GetClustersByName(name)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("cluster '%s' not found", name)
	}
	cluster := parseCluster(items[0])
	return &cluster, nil
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
		return nil, fmt.Errorf("create cluster '%s' failed", name)
	}
	cluster := parseCluster(*item)
	return &cluster, nil
}
func (s clusterService) DeleteCluster(id int) error {
	return dao.DeleteClusterById(id)
}

func DBInit(ctx context.Context) error {
	dbLink, _ := g.Cfg().Get(ctx, "database.link", "/var/lib/skylight/skylight.db")
	glog.Infof(ctx, "database.link: %s", dbLink.String())

	dir := filepath.Dir(dbLink.String())
	if !gfile.Exists(dir) {
		glog.Infof(ctx, "create dir '%s'", dir)
		if err := gfile.Mkdir(dir); err != nil {
			return fmt.Errorf("create dir '%s' failed: %s", dir, err)
		}
	}
	db, err := gorm.Open(sqlite.Open(dbLink.String()), &gorm.Config{})
	if err != nil {
		return err
	}
	glog.Infof(ctx, "migrate db ...")
	orms := []interface{}{
		&do.Cluster{},
		&do.ImageUploadTask{},
	}
	if err := db.AutoMigrate(orms...); err != nil {
		return err
	}
	DB = g.DB().Ctx(ctx)
	return nil
}

var ClusterService *clusterService

func init() {
	ClusterService = &clusterService{}
}
