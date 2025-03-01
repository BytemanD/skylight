package dao

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"skylight/internal/model/entity"
)

const TABLE_AUDITS = "audits"

func modelAudit() *gdb.Model {
	return g.Model(TABLE_AUDITS)
}
func GetAudits() ([]entity.Audit, error) {
	items := []entity.Audit{}
	err := modelAudit().Order("created_at desc").Scan(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
func GetAuditsByProjectId(projectId string) ([]entity.Audit, error) {
	items := []entity.Audit{}
	err := modelAudit().Where("project_id = ?", projectId).Order("created_at desc").Scan(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
func CreateAudit(projectId, projectName, userId, userName, action string) (*entity.Audit, error) {
	item := entity.Audit{
		ProjectId:   projectId,
		ProjectName: projectName,
		UserId:      userId,
		UserName:    userName,
		Action:      action,
	}
	if id, err := modelAudit().Data(item).InsertAndGetId(); err != nil {
		return nil, err
	} else {
		item.Id = id
		return &item, nil
	}
}
