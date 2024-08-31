package entity

type ImageUploadTask struct {
	Id        int    `json:"id,omitempty"`
	ProjectId string `json:"project_id,omitempty"`
	ImageId   string `json:"image_id,omitempty"`
	ImageName string `json:"image_name,omitempty" `
	Size      int    `json:"size"`
	Cached    int    `json:"cached"`
	Uploaded  int    `json:"uploaded"`
}

func (ImageUploadTask) TableName() string {
	return "image_upload_tasks"
}
