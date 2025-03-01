package service

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"skylight/internal/model/entity"
	"skylight/internal/service/openstack"
	"skylight/utility/easyhttp"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

type openstackService struct {
	managers map[string]*openstack.OpenstackManager
}

var RESOURCE_MAP = map[string]string{
	"projects": "项目",
	"users":    "用户",
	"roles":    "角色",

	"servers":     "实例",
	"os-keypairs": "密钥",
	"flavors":     "规格",
	"os-services": "服务", "services": "服务",

	"images": "镜像",

	"volumes":      "卷",
	"volume-types": "卷类型", "volume_types": "卷类型",
	"backups":   "备份",
	"snpashots": "快照",

	"routers":  "路由",
	"networks": "网络",
	"subnets":  "子网",
	"ports":    "端口",
}

func getResourceName(resource string) string {
	if name, ok := RESOURCE_MAP[resource]; ok {
		return name
	} else {
		return resource
	}
}
func (s openstackService) GetLogInfo(req *ghttp.Request) (*openstack.LoginInfo, error) {
	sessionLoginInfo, _ := req.Session.Get("loginInfo", nil)
	loginInfo := openstack.LoginInfo{}

	if err := sessionLoginInfo.Struct(&loginInfo); err != nil {
		return nil, err
	}
	return &loginInfo, nil
}
func (s openstackService) IsLogin(sessionId string) bool {
	_, ok := s.managers[sessionId]
	return ok
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
func (s *openstackService) addAudit(req *ghttp.Request, proxyUrl string) {
	if strings.ToUpper(req.Method) == "DELETE" {
		reg, _ := regexp.Compile("/(.+)/(.+)")
		found := reg.FindStringSubmatch(proxyUrl)
		if len(found) < 3 {
			return
		}
		if err := AuditService.DeleteResoure(req, getResourceName(found[1]), found[2]); err != nil {
			g.Log().Infof(req.GetCtx(), "add audit failed: %s", err)
		}
		return
	}
}

func (s *openstackService) DoProxy(req *ghttp.Request, prefix string) (*easyhttp.Response, error) {
	var (
		resp *easyhttp.Response
		err  error
	)
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
	if err == nil {
		s.addAudit(req, proxyUrl)
		if prefix == "/computing" && strings.ToUpper(req.Method) == "DELETE" {
			go s.watchComputeDeleted(manager, req, proxyUrl)
		}
		if prefix == "/computing" && strings.ToUpper(req.Method) == "POST" {
			body := openstack.Server{}
			if err := resp.UNmarshal(&body); err == nil && body.Server.Id != "" {
				go s.watchComputeCreated(req, manager, fmt.Sprintf("/servers/%s", body.Server.Id))
			}
		}
	}
	return resp, err
}
func (s *openstackService) watchComputeDeleted(manager *openstack.OpenstackManager, req *ghttp.Request, proxyUrl string) {
	if !strings.HasPrefix(proxyUrl, "/servers") {
		return
	}
	values := strings.Split(proxyUrl, "/")
	for {
		resp, _ := manager.ProxyComputing("GET", proxyUrl, nil, nil)
		if resp == nil {
			break
		}
		if resp.IsError() {
			if resp.StatusCode() == 404 {
				SseService.Send(req.GetSessionId(), "success", "实例删除成功", values[2])
				return
			}
		} else {
			data := string(resp.Body())
			SseService.Send(req.GetSessionId(), "info", "更新实例", data)

			body := openstack.Server{}
			if resp.UNmarshal(&body) == nil {
				if strings.ToUpper(body.Server.Status) == "ERROR" {
					SseService.Send(req.GetSessionId(), "error", "实例删除失败", values[2])
					return
				}
			}
		}
		time.Sleep(time.Second * 2)
	}
}
func (s *openstackService) watchComputeCreated(req *ghttp.Request, manager *openstack.OpenstackManager, proxyUrl string) {
	if !strings.HasPrefix(proxyUrl, "/servers") {
		return
	}
	for {
		resp, err := manager.ProxyComputing("GET", proxyUrl, nil, nil)
		if err != nil || resp.IsError() {
			break
		}
		body := openstack.Server{}
		if err := resp.UNmarshal(&body); err != nil {
			break
		}
		data := string(resp.Body())
		SseService.Send(req.GetSessionId(), "info", "更新实例", data)

		switch strings.ToUpper(body.Server.Status) {
		case "ACTIVE":
			SseService.Send(req.GetSessionId(), "success", "实例创建成功", body.Server.Id)
			return
		case "ERROR":
			SseService.Send(req.GetSessionId(), "error", "实例创建失败", body.Server.Id)
			return
		}
		time.Sleep(time.Second * 2)
	}
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
