package openstack

import (
	"skylight/internal/model/entity"
	"skylight/internal/service/internal/dao"
)

type sImageUploadTask struct{}

func (s *sImageUploadTask) Create(projectId string, imageId string, size int) error {
	_, err := dao.CreateImageUploadTask(projectId, imageId, "", size)
	return err
}
func (s *sImageUploadTask) GetByImageId(imageId string) (*entity.ImageUploadTask, error) {
	item, err := dao.GetImageUploadTaskByImageId(imageId)
	if err != nil {
		return nil, err
	}
	return item, nil
}
func (s *sImageUploadTask) GetByProjectId(projectId string) ([]entity.ImageUploadTask, error) {
	items, err := dao.GetImageUploadTasksByProjectId(projectId)
	if err != nil {
		return nil, err
	}
	return items, nil
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
