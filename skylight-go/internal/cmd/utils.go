package cmd

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func getMigraionsPath(ctx context.Context) string {
	migrationsPaths := []string{
		filepath.Join(g.Cfg().MustGet(ctx, "server.dataPath").String(), "migrations"),
		filepath.Join("../migrations"),
	}
	for _, path := range migrationsPaths {
		pathAbs, err := filepath.Abs(path)
		if err != nil {
			continue
		}
		if gfile.IsDir(pathAbs) {
			return pathAbs
		}
	}
	return ""

}
func sqliteMigrate(ctx context.Context, dbConf *gdb.ConfigNode) {
	sourcePath := getMigraionsPath(ctx)
	if sourcePath == "" {
		panic("source path is empty")
	}
	m, err := migrate.New(
		fmt.Sprintf("file://%s", sourcePath),
		fmt.Sprintf("sqlite3://%s", dbConf.Name))
	if err != nil {
		panic(err)
	}
	defer m.Close()

	version, dirty, _ := m.Version()
	g.Log().Info(ctx, "MIGRATE currently version=%d, dirty=%t", version, dirty)
	g.Log().Infof(ctx, "MIGRATE sourceUrl: %s, databaseUrl: %s", sourcePath, dbConf.Name)
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}
	}
}

func InitDB(ctx context.Context) {
	dbConf := g.DB().GetConfig()

	g.Log().Infof(ctx, "DB type: %v", dbConf.Type)
	switch dbConf.Type {
	case "sqlite":
		g.Log().Infof(ctx, "DB name: %s", dbConf.Name)
		sqliteMigrate(ctx, dbConf)
	default:
		panic("invalid db type, only support: sqlite")
	}
}
