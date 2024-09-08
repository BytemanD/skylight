package service

import (
	"fmt"
	"skylight/internal/model/entity"
	"skylight/internal/service/internal/dao"
	"skylight/internal/service/internal/do"

	"github.com/BytemanD/easygo/pkg/global/logging"
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
func (s auditService) Login(req *ghttp.Request) {
	loginInfo, err := OSService.GetLogInfo(req)
	if err != nil {
		glog.Infof(req.GetCtx(), "get login info failed: %s", err)
		return
	}
	_, err = dao.CreateAudit(
		loginInfo.Project.Id, loginInfo.Project.Name,
		loginInfo.User.Id, loginInfo.User.Name,
		fmt.Sprintf("登录集群 %s", loginInfo.Cluster),
	)
	if err != nil {
		logging.Error("create audit failed: %s", err)
	}
}
func (s auditService) Logout(req *ghttp.Request) {
	loginInfo, err := OSService.GetLogInfo(req)
	if err != nil {
		glog.Infof(req.GetCtx(), "get login info failed: %s", err)
		return
	}
	_, err = dao.CreateAudit(
		loginInfo.Project.Id, loginInfo.Project.Name,
		loginInfo.User.Id, loginInfo.User.Name,
		fmt.Sprintf("退出集群 %s", loginInfo.Cluster),
	)
	if err != nil {
		logging.Error("create audit failed: %s", err)
	}
}

var AuditService *auditService

func init() {
	AuditService = &auditService{}
}
