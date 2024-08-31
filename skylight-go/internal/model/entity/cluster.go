package entity

type Cluster struct {
	Id      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	AuthUrl string `json:"auth_url,omitempty"`
}
