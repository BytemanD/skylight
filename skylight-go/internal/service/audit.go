package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"skylight/internal/model/entity"
	"skylight/internal/service/internal/dao"
	"skylight/internal/service/internal/do"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

type auditService struct{}

func parseAudit(item do.Audit) entity.Audit {
	return entity.Audit{
		Id:          item.Id,
		ProjectId:   item.ProjectId,
		ProjectName: item.ProjectName,
		UserId:      item.UserId,
		UserName:    item.UserName,
		Action:      item.Action,
		CreatedAt:   item.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
func parseAudits(items []do.Audit) []entity.Audit {
	audits := []entity.Audit{}
	for _, item := range items {
		audits = append(audits, parseAudit(item))
	}
	return audits
}

// cluster
func (s auditService) GetAll() ([]entity.Audit, error) {
	items, err := dao.GetAudits()
	if err != nil {
		return nil, err
	}
	return parseAudits(items), nil
}
func (s auditService) GetByProjectId(projectId string) ([]entity.Audit, error) {
	items, err := dao.GetAuditsByProjectId(projectId)
	if err != nil {
		return nil, err
	}
	return parseAudits(items), nil
}

func (s auditService) Create(projectId, projectName, userId, userName, action string) (*entity.Audit, error) {
	item, err := dao.CreateAudit(projectId, projectName, userId, userName, action)
	if err != nil {
		return nil, err
	}
	audit := parseAudit(*item)
	return &audit, nil
}
func (s auditService) Login(req *ghttp.Request) error {
	loginInfo, err := OSService.GetLogInfo(req)
	if err != nil {
		// glog.Infof(req.GetCtx(), "get login info failed: %s", err)
		return errors.Join(errors.New("get login info failed"), err)
	}
	glog.Info(req.GetCtx(), "============ create audit= %s")
	_, err = dao.CreateAudit(
		loginInfo.Project.Id, loginInfo.Project.Name,
		loginInfo.User.Id, loginInfo.User.Name,
		fmt.Sprintf("登录集群 %s", loginInfo.Cluster),
	)
	if err != nil {
		return errors.Join(errors.New("create audit failed"), err)
	}
	return nil
}
func (s auditService) Logout(req *ghttp.Request) error {
	loginInfo, err := OSService.GetLogInfo(req)
	if err != nil {
		glog.Errorf(req.GetCtx(), "get login info failed: %s", err)
		return err
	}
	OSService.RemoveManager(req.GetSessionId())
	if err := req.Session.RemoveAll(); err != nil {
		return err
	}
	gsessionPath, _ := g.Cfg().Get(req.GetCtx(), "session.path", "/var/lib/skylight/gsessions")
	gsessionFile := filepath.Join(gsessionPath.String(), req.GetSessionId())

	if err = os.Remove(gsessionFile); err != nil {
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
		glog.Errorf(req.GetCtx(), "get login info failed: %s", err)
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

var AuditService *auditService

func init() {
	AuditService = &auditService{}
}
