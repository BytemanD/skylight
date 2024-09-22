package openstack

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BytemanD/easygo/pkg/global/logging"
)

type Domain struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type Scope struct {
	Project Project `json:"project,omitempty"`
}
type Project struct {
	Id          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Domain      Domain   `json:"domain,omitempty"`
	Description string   `json:"description,omitempty"`
	Enabled     bool     `json:"enabled,omitempty"`
	DomainId    string   `json:"domain_id,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	IsDomain    bool     `json:"is_domain,omitempty"`
	ParentId    string   `json:"parent_id,omitempty"`
}
type User struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Password    string `json:"password,omitempty"`
	Project     string `json:"project,omitempty"`
	Description string `json:"description,omitempty"`
	Email       string `json:"email,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
	Domain      Domain `json:"domain,omitempty"`
	DomainId    string `json:"domain_id,omitempty"`
}
type Password struct {
	User User `json:"user"`
}
type Identity struct {
	Methods  []string `json:"methods,omitempty"`
	Password Password `json:"password,omitempty"`
}

type Auth struct {
	Identity Identity `json:"identity,omitempty"`
	Scope    Scope    `json:"scope,omitempty"`
}

type Endpoint struct {
	Id        string `json:"id"`
	Region    string `json:"region"`
	Url       string `json:"url"`
	Interface string `json:"interface"`
	RegionId  string `json:"region_id"`
	ServiceId string `json:"service_id"`
}

type Catalog struct {
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Id        string     `json:"id"`
	Endpoints []Endpoint `json:"endpoints"`
}
type Role struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TokenBody struct {
	IsDomain bool      `json:"is_domain"`
	User     User      `json:"user"`
	Project  Project   `json:"project"`
	Roles    []Role    `json:"roles"`
	Catalogs []Catalog `json:"catalog"`
}

type LoginInfo struct {
	Cluster  string  `json:"cluster,omitempty"`
	Region   string  `json:"region,omitempty"`
	Project  Project `json:"project,omitempty"`
	User     User    `json:"user,omitempty"`
	Roles    []Role  `json:"roles,omitempty"`
	Password string  `json:"password,omitempty"`
}

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
	logging.Debug("upload image %s to glance, process: %d/%d",
		progress.ImageId, progress.writedSize, progress.Total)
	progress.writedSize += len(p)
	percent := progress.writedSize * 100 / progress.Total
	if progress.percent != percent {
		progress.percent = percent
		err := ImageUploadTaskService.UpdateUploaded(progress.ImageId, progress.writedSize)
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

type Resource struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
type Server struct {
	Server Resource `json:"server"`
}
