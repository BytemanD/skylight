package do

type ImageUploadTask struct {
	Id        int    `gorm:"id,primary,autoinc" json:"id,omitempty"`
	ProjectId string `gorm:"project_id"           `
	ImageId   string `gorm:"image_id,primary"   `
	ImageName string `gorm:"image_name" `
	Size      int    `gorm:"total"      `
	Cached    int    `gorm:"cached"             `
	Uploaded  int    `gorm:"uploaded"           `
}

func (ImageUploadTask) TableName() string {
	return "image_upload_tasks"
}
