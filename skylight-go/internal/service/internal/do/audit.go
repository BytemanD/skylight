package do

import (
	"time"
)

type Audit struct {
	Id          int       `gorm:"primary_key,autoinc" json:"id,omitempty"`
	ProjectId   string    `gorm:"column:project_id;not null"`
	ProjectName string    `gorm:"column:project_name"`
	UserId      string    `gorm:"column:user_id;not null"`
	UserName    string    `gorm:"column:user_name"`
	Action      string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at;type:datetime;not null"`
}

type Audits []Audit

func (Audit) TableName() string {
	return "audits"
}

func (Audits) TableName() string {
	return "audits"
}
