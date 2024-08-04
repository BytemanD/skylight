package service

import (
	"context"
	"fmt"
	"path/filepath"
	"skylight/internal/model"
	"skylight/internal/service/internal/dao"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB gdb.DB

func GetClusters() ([]model.Cluster, error) {
	return dao.GetClusters(DB)
}
func GetClusterByName(name string) (*model.Cluster, error) {
	clusters, err := dao.GetClustersByName(DB, name)
	if err != nil {
		return nil, err
	}
	if len(clusters) == 0 {
		return nil, fmt.Errorf("cluster '%s' not found", name)
	}
	return &clusters[0], nil
}

func CreatCluster(name string, authUrl string) (*model.Cluster, error) {
	clusters, err := dao.GetClustersByName(DB, name)
	if err != nil {
		return nil, err
	}
	if len(clusters) > 0 {
		return nil, fmt.Errorf("cluster '%s' already exists", name)
	}
	return dao.CreateCluster(DB, name, authUrl)
}

func DBInit(ctx context.Context, dbDriver, dbLink string) error {
	if dbDriver != "sqlite" {
		return fmt.Errorf("invalid databse driver: %s", dbDriver)
	}
	dir := filepath.Dir(dbLink)
	if !gfile.Exists(dir) {
		logging.Info("create dir '%s'", dir)
		if err := gfile.Mkdir(dir); err != nil {
			return fmt.Errorf("create dir '%s' failed: %s", dir, err)
		}
	}
	db, err := gorm.Open(sqlite.Open(dbLink), &gorm.Config{})
	if err != nil {
		return err
	}
	logging.Info("migrate db ...")
	if err := db.AutoMigrate(&model.Cluster{}); err != nil {
		return err
	}
	DB = g.DB().Ctx(ctx)
	return nil
}
