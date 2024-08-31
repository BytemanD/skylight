package service

import (
	"skylight/internal/model/entity"
	"skylight/internal/service/internal/dao"
	"skylight/internal/service/internal/do"
)

type sImageUploadTask struct{}

func parseImageUploadTask(item do.ImageUploadTask) entity.ImageUploadTask {
	return entity.ImageUploadTask{
		Id:        item.Id,
		ProjectId: item.ProjectId,
		ImageId:   item.ImageId,
		ImageName: item.ImageName,
		Size:      item.Size,
		Cached:    item.Cached,
		Uploaded:  item.Uploaded,
	}
}
func parseImageUploadTasks(items []do.ImageUploadTask) []entity.ImageUploadTask {
	tasks := []entity.ImageUploadTask{}
	for _, item := range items {
		tasks = append(tasks, parseImageUploadTask(item))
	}
	return tasks
}

func (s *sImageUploadTask) Create(projectId string, imageId string, size int) error {
	_, err := dao.CreateImageUploadTask(projectId, imageId, "", size)
	return err
}
func (s *sImageUploadTask) GetByImageId(imageId string) (*entity.ImageUploadTask, error) {
	item, err := dao.GetImageUploadTaskByImageId(imageId)
	if err != nil {
		return nil, err
	}
	task := parseImageUploadTask(*item)
	return &task, nil
}
func (s *sImageUploadTask) GetByProjectId(projectId string) ([]entity.ImageUploadTask, error) {
	items, err := dao.GetImageUploadTasksByProjectId(projectId)
	if err != nil {
		return nil, err
	}
	return parseImageUploadTasks(items), nil
}
func (s *sImageUploadTask) Delete(id int) error {
	return dao.DeleteImageUploadTask(id)
}

func (s *sImageUploadTask) IncrementCached(imageId string, v int) error {
	return dao.IncrementImageUploadCached(imageId, v)
}
func (s *sImageUploadTask) IncrementUploaded(imageId string, v int) error {
	return dao.IncrementImageUploadTaskUploaded(imageId, v)
}
func (s *sImageUploadTask) UpdateUploaded(imageId string, v int) error {
	return dao.UpdateImageUploadTaskUploaded(imageId, v)
}

var ImageUploadTaskService *sImageUploadTask

func init() {
	ImageUploadTaskService = &sImageUploadTask{}
}
