package dao

import (
	"skylight/internal/service/internal/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func queryAudit() *gdb.Model {
	return g.DB().Model(do.Audit{})
}
func GetAudits() ([]do.Audit, error) {
	items := []do.Audit{}
	err := queryAudit().Order("created_at desc").Scan(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
func GetAuditsByProjectId(projectId string) ([]do.Audit, error) {
	items := do.Audits{}
	err := queryAudit().Where("project_id = ?", projectId).Order("created_at desc").Scan(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
func CreateAudit(projectId, projectName, userId, userName, action string) (*do.Audit, error) {
	item := do.Audit{
		ProjectId:   projectId,
		ProjectName: projectName,
		UserId:      userId,
		UserName:    userName,
		Action:      action,
	}
	if result, err := queryAudit().Insert(item); err != nil {
		return nil, err
	} else {
		id, _ := result.LastInsertId()
		item.Id = int(id)
		return &item, nil
	}
}
