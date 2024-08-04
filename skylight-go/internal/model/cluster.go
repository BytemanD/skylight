package model

type Cluster struct {
	Id      int    `gorm:"id,primary,autoinc" json:"id,omitempty"`
	Name    string `gorm:"name,primary"        json:"name,omitempty"`
	AuthUrl string `gorm:"auth_url"           json:"auth_url,omitempty"`
}

func (Cluster) TableName() string {
	return "clusters"
}

type Clusters []Cluster

func (Clusters) TableName() string {
	return "clusters"
}
