package entity

type Audit struct {
	Id          int    `json:"id,omitempty"`
	ProjectId   string `json:"project_id,omitempty"`
	ProjectName string `json:"project_name,omitempty"`
	UserId      string `json:"user_id,omitempty"`
	UserName    string `json:"user_name,omitempty"`
	Action      string `json:"action"`
	CreatedAt   string `json:"created_at"`
}
