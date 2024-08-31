package dao

import (
	"fmt"
	"skylight/internal/service/internal/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func queryImageUploadTask() *gdb.Model {
	return g.DB().Model(do.ImageUploadTask{})
}

func GetImageUploadTasks() ([]do.ImageUploadTask, error) {
	items := []do.ImageUploadTask{}
	err := queryImageUploadTask().Scan(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
func GetImageUploadTasksByProjectId(projectId string) ([]do.ImageUploadTask, error) {
	items := []do.ImageUploadTask{}
	err := queryImageUploadTask().Where("project_id", projectId).Scan(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
func GetImageUploadTaskByImageId(imageId string) (*do.ImageUploadTask, error) {
	items := []do.ImageUploadTask{}
	err := queryImageUploadTask().Where("image_id", imageId).Scan(&items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("task not found")
	}
	return &(items[0]), nil
}
func CreateImageUploadTask(projectId, imageId, imageName string, size int) (*do.ImageUploadTask, error) {
	item := do.ImageUploadTask{
		ProjectId: projectId,
		ImageId:   imageId,
		ImageName: imageName,
		Size:      size,
		Cached:    0,
		Uploaded:  0,
	}
	if _, err := queryImageUploadTask().Insert(item); err != nil {
		return nil, err
	}
	return &item, nil
}
func DeleteImageUploadTask(id int) error {
	_, err := queryImageUploadTask().Delete("id = ?", id)
	return err
}
func DeleteImageUploadTaskByImageId(imageId string) error {
	_, err := queryImageUploadTask().Delete("image_id = ?", imageId)
	return err
}
func IncrementImageUploadCached(imageId string, cached int) error {
	_, err := queryImageUploadTask().Where("image_id", imageId).Increment("cached", cached)
	return err
}
func IncrementImageUploadTaskUploaded(imageId string, uploaded int) error {
	if _, err := queryImageUploadTask().Where("image_id", imageId).Increment("uploaded", uploaded); err != nil {
		return err
	}
	return nil
}
func UpdateImageUploadTaskUploaded(imageId string, uploaded int) error {
	_, err := queryImageUploadTask().Data(
		g.Map{"uploaded": uploaded},
	).Where("image_id", imageId).Update()
	return err
}
