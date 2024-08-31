package openstack

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"skylight/internal/service"

	"github.com/BytemanD/easygo/pkg/global/logging"
)

type ImageUploadProgress struct {
	ImageId    string
	Total      int
	writedSize int
	percent    int
}

func (progress *ImageUploadProgress) Percent() int {
	return progress.percent
}
func (progress *ImageUploadProgress) Write(p []byte) (int, error) {
	progress.writedSize += len(p)
	percent := progress.writedSize * 100 / progress.Total
	if progress.percent != percent {
		progress.percent = percent
		err := service.ImageUploadTaskService.UpdateUploaded(progress.ImageId, progress.writedSize)
		if err != nil {
			logging.Error("update %s uploaded failed: %s", progress.ImageId, err)
		}
	}
	return len(p), nil
}

func ImageUploadBufReader(imageFile string) (*bufio.Reader, error) {
	fileInfo, err := os.Stat(imageFile)
	if err != nil {
		return nil, fmt.Errorf("get image stat failed: %s", err)
	}
	reader, err := os.Open(imageFile)
	if err != nil {
		return nil, err
	}
	wc := &ImageUploadProgress{ImageId: filepath.Base(imageFile), Total: int(fileInfo.Size())}
	return bufio.NewReader(io.TeeReader(reader, wc)), nil
}
