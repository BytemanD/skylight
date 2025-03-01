package service

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"skylight/internal/model/entity"
	"skylight/internal/service/internal/dao"
	"skylight/utility"
	"time"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gsession"
)

type auditService struct{}

var AuditService *auditService
var SessionStorage *gsession.StorageFile

// cluster
func (s auditService) GetAll() ([]entity.Audit, error) {
	items, err := dao.GetAudits()
	if err != nil {
		return nil, err
	}
	return items, nil
}
func (s auditService) GetByProjectId(projectId string) ([]entity.Audit, error) {
	items, err := dao.GetAuditsByProjectId(projectId)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s auditService) Create(projectId, projectName, userId, userName, action string) (*entity.Audit, error) {
	return dao.CreateAudit(projectId, projectName, userId, userName, action)
}
func (s auditService) Login(req *ghttp.Request) error {
	loginInfo, err := OSService.GetLogInfo(req)
	if err != nil {
		// g.Log().Infof(req.GetCtx(), "get login info failed: %s", err)
		return errors.Join(errors.New("get login info failed"), err)
	}
	_, err = dao.CreateAudit(
		loginInfo.Project.Id, loginInfo.Project.Name,
		loginInfo.User.Id, loginInfo.User.Name,
		fmt.Sprintf("登录集群 %s", loginInfo.Cluster),
	)
	if err != nil {
		return fmt.Errorf("create audit failed: %s", err)
	}
	return nil
}
func (s auditService) Logout(req *ghttp.Request) error {
	loginInfo, err := OSService.GetLogInfo(req)
	if err != nil {
		g.Log().Errorf(req.GetCtx(), "get login info failed: %s", err)
		return err
	}
	OSService.RemoveManager(req.GetSessionId())
	if err := req.Session.RemoveAll(); err != nil {
		return err
	}
	_, err = dao.CreateAudit(
		loginInfo.Project.Id, loginInfo.Project.Name,
		loginInfo.User.Id, loginInfo.User.Name,
		fmt.Sprintf("退出集群 %s", loginInfo.Cluster),
	)
	if err != nil {
		logging.Error("create audit failed: %s", err)
	}
	return nil
}
func (s auditService) DeleteResoure(req *ghttp.Request, name string, resource string) error {
	loginInfo, err := OSService.GetLogInfo(req)
	if err != nil {
		g.Log().Errorf(req.GetCtx(), "get login info failed: %s", err)
		return err
	}
	_, err = dao.CreateAudit(
		loginInfo.Project.Id, loginInfo.Project.Name,
		loginInfo.User.Id, loginInfo.User.Name,
		fmt.Sprintf("删除 %s %s", name, resource),
	)
	if err != nil {
		logging.Error("create audit failed: %s", err)
	}
	return nil
}

func InitSessionStorage(ctx context.Context) {
	gsessionPath := g.Cfg().MustGet(ctx, "session.default.path")
	g.Log().Infof(ctx, "SESSION type: %s", g.Cfg().MustGet(ctx, "session.default.type").String())
	g.Log().Infof(ctx, "SESSION path: %s", gsessionPath.String())
	utility.MakesureDir(gsessionPath.String())
	SessionStorage = gsession.NewStorageFile(gsessionPath.String(), time.Hour*3)

	// TODO: 周期任务
	filepath.Walk(filepath.Join(gsessionPath.String()),
		func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if !info.ModTime().Before(time.Now().AddDate(0, 0, -1)) {
				return nil
			}
			g.Log().Infof(ctx, "session %s is expired, cleanup", info.Name())
			os.Remove(path)
			return nil
		},
	)

}

func init() {
	AuditService = &auditService{}
}
