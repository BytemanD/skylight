package dao

import (
	"fmt"
	"skylight/internal/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

const TABLE_IMAGE_UPLOAD_TASKS = "image_upload_tasks"

func modelImageUploadTasks() *gdb.Model {
	return g.DB().Model(TABLE_IMAGE_UPLOAD_TASKS)
}

func GetImageUploadTasks() ([]entity.ImageUploadTask, error) {
	items := []entity.ImageUploadTask{}
	err := modelImageUploadTasks().Scan(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
func GetImageUploadTasksByProjectId(projectId string) ([]entity.ImageUploadTask, error) {
	items := []entity.ImageUploadTask{}
	err := modelImageUploadTasks().Where("project_id", projectId).Scan(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
func GetImageUploadTaskByImageId(imageId string) (*entity.ImageUploadTask, error) {
	items := []entity.ImageUploadTask{}
	err := modelImageUploadTasks().Where("image_id", imageId).Scan(&items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("task not found")
	}
	return &(items[0]), nil
}
func CreateImageUploadTask(projectId, imageId, imageName string, size int) (*entity.ImageUploadTask, error) {
	item := entity.ImageUploadTask{
		ProjectId: projectId,
		ImageId:   imageId,
		ImageName: imageName,
		Size:      size,
		Cached:    0,
		Uploaded:  0,
	}
	if _, err := modelImageUploadTasks().Insert(item); err != nil {
		return nil, err
	}
	return &item, nil
}
func DeleteImageUploadTask(id int) error {
	_, err := modelImageUploadTasks().Delete("id = ?", id)
	return err
}
func DeleteImageUploadTaskByImageId(imageId string) error {
	_, err := modelImageUploadTasks().Delete("image_id = ?", imageId)
	return err
}
func IncrementImageUploadCached(imageId string, cached int) error {
	_, err := modelImageUploadTasks().Where("image_id", imageId).Increment("cached", cached)
	return err
}
func IncrementImageUploadTaskUploaded(imageId string, uploaded int) error {
	if _, err := modelImageUploadTasks().Where("image_id", imageId).Increment("uploaded", uploaded); err != nil {
		return err
	}
	return nil
}
func UpdateImageUploadTaskUploaded(imageId string, uploaded int) error {
	_, err := modelImageUploadTasks().Data(
		g.Map{"uploaded": uploaded},
	).Where("image_id", imageId).Update()
	return err
}
