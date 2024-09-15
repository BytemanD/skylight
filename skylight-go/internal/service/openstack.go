package service

import (
	"fmt"
	"os"
	"path/filepath"
	"skylight/internal/model/entity"
	"skylight/internal/service/openstack"
	"skylight/utility/easyhttp"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

type openstackService struct {
	managers map[string]*openstack.OpenstackManager
}

func (s openstackService) GetLogInfo(req *ghttp.Request) (*openstack.LoginInfo, error) {
	sessionLoginInfo, _ := req.Session.Get("loginInfo", nil)
	loginInfo := openstack.LoginInfo{}

	if err := sessionLoginInfo.Struct(&loginInfo); err != nil {
		return nil, err
	}
	return &loginInfo, nil
}

func (s openstackService) GetManager(req *ghttp.Request) (*openstack.OpenstackManager, error) {
	sessionId := req.GetSessionId()
	if client, ok := s.managers[sessionId]; ok {
		return client, nil
	}
	if loginInfo, err := s.GetLogInfo(req); err != nil {
		return nil, fmt.Errorf("get session auth info falied: %v", err)
	} else {
		cluster, err := ClusterService.GetClusterByName(loginInfo.Cluster)
		if err != nil {
			return nil, fmt.Errorf("get cluster %s failed: %s", loginInfo.Cluster, err)
		}
		manager, err := openstack.NewManager(
			cluster.AuthUrl, loginInfo.Project.Name, loginInfo.User.Name, loginInfo.Password)
		if err != nil {
			return nil, err
		}
		manager.SetRegion(loginInfo.Region)
		s.managers[sessionId] = manager
		return manager, nil
	}
}

func (s *openstackService) RemoveManager(sessionId string) {
	delete(s.managers, sessionId)
}

func (s *openstackService) SetRegion(req *ghttp.Request, region string) error {
	if loginInfo, err := s.GetLogInfo(req); err != nil {
		return fmt.Errorf("get session auth info falied: %v", err)
	} else {
		loginInfo.Region = region
		req.Session.Set("loginInfo", loginInfo)
	}
	manager, err := s.GetManager(req)
	if err != nil {
		return err
	}
	manager.SetRegion(region)
	return nil
}
func (s *openstackService) DoProxy(req *ghttp.Request, prefix string) (*easyhttp.Response, error) {
	var resp *easyhttp.Response
	var err error
	if err != nil {
		return nil, fmt.Errorf("get session failed: %s", err)
	}
	manager, err := s.GetManager(req)
	if err != nil {
		return nil, fmt.Errorf("get manager failed: %s", err)
	}

	proxyUrl := strings.TrimPrefix(req.URL.Path, prefix)
	switch prefix {
	case "/computing":
		resp, err = manager.ProxyComputing(req.Method, proxyUrl, req.URL.Query(), req.GetBody())
	case "/networking":
		resp, err = manager.ProxyNetworking(req.Method, proxyUrl, req.URL.Query(), req.GetBody())
	case "/volume":
		resp, err = manager.ProxyVolume(req.Method, proxyUrl, req.URL.Query(), req.GetBody())
	case "/image":
		resp, err = manager.ProxyImage(proxyUrl, req)
	case "/identity":
		resp, err = manager.ProxyIdentity(req.Method, proxyUrl, req.URL.Query(), req.GetBody())
	default:
		err = fmt.Errorf("invalid prefix %s", prefix)
	}
	return resp, err
}

func (s *openstackService) SaveImageCache(proxyUrl string, req *ghttp.Request) (*entity.ImageUploadTask, error) {
	metaSize := req.Header.Get("x-image-meta-size")
	imageSize, err := strconv.Atoi(metaSize)
	if err != nil {
		return nil, fmt.Errorf("image size not found in body")
	}
	imageId, err := openstack.GetImageIdFromProxyUrl(proxyUrl)
	if err != nil {
		return nil, err
	}
	dataPath, _ := g.Cfg().Get(gctx.New(), "server.dataPath")
	cacheFile := filepath.Join(dataPath.String(), "image_cache", imageId)
	if !gfile.Exists(cacheFile) {
		loginInfo, err := s.GetLogInfo(req)
		if err != nil {
			return nil, fmt.Errorf("get login info failed: %s", err)
		}
		err = openstack.ImageUploadTaskService.Create(loginInfo.Project.Id, imageId, imageSize)
		if err != nil {
			return nil, fmt.Errorf("create image upload task failed: %s", err)
		}
	}
	file, err := os.OpenFile(cacheFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data := req.GetBody()
	_, err = file.Write(data)
	if err != nil {
		return nil, fmt.Errorf("save data to cache failed: %s", err)
	}
	err = openstack.ImageUploadTaskService.IncrementCached(imageId, len(data))
	if err != nil {
		return nil, fmt.Errorf("increment %s cached failed: %s", imageId, err)
	}
	task, err := openstack.ImageUploadTaskService.GetByImageId(imageId)
	if err != nil {
		return nil, fmt.Errorf("get task for %s failed: %s", imageId, err)
	}
	return task, err
}

var SESSION_MANAGER_MAP map[string]*openstack.OpenstackManager
var OSService *openstackService

func init() {
	OSService = &openstackService{
		managers: map[string]*openstack.OpenstackManager{},
	}
}
