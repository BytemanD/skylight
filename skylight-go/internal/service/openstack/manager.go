package openstack

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"skylight/internal/model"
	"skylight/internal/model/entity"
	"skylight/internal/service"
	"skylight/utility/easyhttp"
	"strconv"
	"strings"
	"time"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/go-resty/resty/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
)

type OpenstackManager struct {
	AuthUrl  string
	Region   string
	AuthInfo Auth

	session2        *easyhttp.Client
	token           string
	tokenAlive      time.Duration
	expiredAt       time.Time
	tokenData       TokenBody
	serviceEndpoint map[string]string
	microVesrion    map[string]string
}

func (c *OpenstackManager) isTokenExpired() (expired bool) {
	defer func() {
		if expired {
			logging.Warning("token expired")
		}
	}()
	if c.token == "" {
		expired = true
	} else {
		expired = c.expiredAt.Before(time.Now())
	}
	return expired
}

func (c *OpenstackManager) sendToBackend2(req *easyhttp.Request) (*easyhttp.Response, error) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	logging.Debug("-----proxy---> %s %s ?%s Headers: %s",
		req.Method, req.URL, req.QueryValues.Encode(), c.safeHeader(req.Header))
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}
	proxyRespBody := "<Body>"
	if resp.GetContentType() == easyhttp.APPLICATION_JSON {
		if resp.IsError() {
			if reqBody, err := req.GetBytesBody(); err == nil {
				logging.Debug("req body: %s", string(reqBody))
				proxyRespBody = string(reqBody)
			} else {
				logging.Error("get bytes body failed: %s", err)
			}
		}
	}
	logging.Debug("proxy Resp [%d] %s", resp.StatusCode(), proxyRespBody)
	return resp, err
}
func (c *OpenstackManager) tokenIssue() error {
	req := c.session2.NewRequest()

	req.SetJsonBody(map[string]Auth{"auth": c.AuthInfo})
	req.Method = resty.MethodPost
	if reqUrl, err := url.JoinPath(c.AuthUrl, "/auth/tokens"); err != nil {
		return err
	} else {
		req.URL = reqUrl
	}
	resp, err := c.sendToBackend2(req)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("login failed: [%d] %s", resp.StatusCode(), resp.Body())
	}
	c.token = resp.GetHeader("X-Subject-Token")
	respBody := struct{ Token TokenBody }{}
	if err := resp.UNmarshal(&respBody); err != nil {
		return err
	}
	c.tokenData = respBody.Token
	c.expiredAt = time.Now().Add(c.tokenAlive)
	return nil
}

func (c *OpenstackManager) GetToken() (string, error) {
	if !c.isTokenExpired() {
		return c.token, nil
	}
	return c.token, c.tokenIssue()
}

func (c *OpenstackManager) GetEndpoint(service string) (string, error) {
	if _, err := c.GetToken(); err != nil {
		return "", err
	}
	for _, catalog := range c.tokenData.Catalogs {
		if catalog.Type != service && catalog.Name != service {
			continue
		}
		for _, endpoint := range catalog.Endpoints {
			if endpoint.Interface == "public" && endpoint.Region == c.Region {
				return endpoint.Url, nil
			}
		}
	}
	return "", fmt.Errorf("endpoint for service %s:%s not found", c.Region, service)
}
func (c *OpenstackManager) clearEndpoints() {
	c.serviceEndpoint = map[string]string{}
}
func (c *OpenstackManager) makeSureEndpoint(service, version string) (err error) {
	if _, ok := c.serviceEndpoint[service]; ok {
		return
	}
	endpoint, err := c.GetEndpoint(service)
	if err == nil {
		if u, err2 := url.Parse(endpoint); err2 != nil {
			return err2
		} else {
			if u.Path == "" || u.Path == "/" {
				c.serviceEndpoint[service], err = url.JoinPath(endpoint, version)
			} else {
				c.serviceEndpoint[service], err = endpoint, nil
			}
		}
	}
	return
}
func (c *OpenstackManager) safeHeader(h http.Header) http.Header {
	headers := http.Header{}
	for k, v := range h {
		if k == "X-Auth-Token" {
			headers.Set(k, "<TOKEN>")
		} else {
			headers.Set(k, strings.Join(v, ","))
		}
	}
	return headers
}
func (c *OpenstackManager) doProxy2(endpoint string, method string,
	u string, q url.Values, body []byte) (*easyhttp.Response, error) {
	return c.doProxyWithHeaders2(endpoint, method, u, q, nil, body)
}
func (c *OpenstackManager) doProxyWithHeaders2(endpoint string, method string,
	u string, q url.Values, headers map[string]string, body interface{}) (*easyhttp.Response, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	var reqUrl string
	if u != "" && u != "/" {
		reqUrl, err = url.JoinPath(endpoint, u)
	} else {
		reqUrl = endpoint
	}
	if err != nil {
		return nil, err
	}
	req := c.session2.NewRequest().SetMethod(method).SetURL(reqUrl).
		SetHeader("X-Auth-Token", token).
		AddQueryValuesFromValues(q).
		SetHeaders(headers).
		SetJsonBody(body)
	return c.sendToBackend2(req)
}
func (c *OpenstackManager) doProxyWithHeaders2BodyReader(endpoint string, method string,
	u string, q url.Values, headers map[string]string, body *bufio.Reader) (*easyhttp.Response, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	var reqUrl string
	if u != "" && u != "/" {
		if reqUrl, err = url.JoinPath(endpoint, u); err != nil {
			return nil, err
		}
	} else {
		reqUrl = endpoint
	}
	req := c.session2.NewRequest().SetMethod(method).SetURL(reqUrl).
		SetHeader("X-Auth-Token", token).
		AddQueryValuesFromValues(q).
		SetHeaders(headers).
		SetReaderBody(body)
	return c.sendToBackend2(req)
}

func (c *OpenstackManager) SetAuthUrl(authUrl string) {
	u, _ := url.Parse(authUrl)
	if u.Path == "" || u.Path == "/" {
		authUrl, _ = url.JoinPath(authUrl, "v3")
	}
	c.AuthUrl = authUrl
}
func (c *OpenstackManager) SetRegion(req *ghttp.Request, region string) error {
	if loginInfo, err := GetAuthFromSession(req); err != nil {
		return fmt.Errorf("get session auth info falied: %v", err)
	} else {
		loginInfo.Region = region
		req.Session.Set("loginInfo", loginInfo)
	}

	c.Region = region
	c.clearEndpoints()
	return nil
}
func (c *OpenstackManager) GetUser() User {
	return c.tokenData.User
}
func (c *OpenstackManager) GetProject() Project {
	return c.tokenData.Project
}
func (c *OpenstackManager) GetRoles() []Role {
	return c.tokenData.Roles
}
func (c *OpenstackManager) ProxyIdentity(method string, url string, q url.Values, body []byte) (*easyhttp.Response, error) {
	if err := c.makeSureEndpoint("identity", "v3"); err != nil {
		return nil, err
	}
	return c.doProxy2(c.serviceEndpoint["identity"], method, url, q, body)
}
func (c *OpenstackManager) ProxyNetworking(method string, url string, q url.Values, body []byte) (*easyhttp.Response, error) {
	if err := c.makeSureEndpoint("neutron", "v2.0"); err != nil {
		return nil, err
	}
	return c.doProxy2(c.serviceEndpoint["neutron"], method, url, q, body)
}

func (c *OpenstackManager) getMicroVersion(service string) (string, error) {
	resp, err := c.doProxy2(c.serviceEndpoint[service], resty.MethodGet, "/", nil, nil)
	if err != nil {
		return "", err
	}
	versionBody := struct {
		Version model.Version
	}{Version: model.Version{}}
	if err := json.Unmarshal(resp.Body(), &versionBody); err != nil {
		return "", err
	}
	return versionBody.Version.Version, nil
}
func (c *OpenstackManager) ProxyComputing(method string, url string, q url.Values, body []byte) (*easyhttp.Response, error) {
	if err := c.makeSureEndpoint("nova", "v2.1"); err != nil {
		return nil, err
	}
	if _, ok := c.microVesrion["nova"]; !ok {
		microVersion, err := c.getMicroVersion("nova")
		if err != nil {
			logging.Warning("get microversion failed: %s", err)
			c.microVesrion["nova"] = ""
		} else {
			c.microVesrion["nova"] = microVersion
		}
	}
	if microVersion, ok := c.microVesrion["nova"]; ok && microVersion != "" {
		headers := map[string]string{
			"X-OpenStack-Nova-API-Version": microVersion,
			"OpenStack-API-Version":        microVersion,
		}
		return c.doProxyWithHeaders2(c.serviceEndpoint["nova"], method, url, q, headers, body)
	} else {
		return c.doProxy2(c.serviceEndpoint["nova"], method, url, q, body)
	}
}
func (c *OpenstackManager) ProxyVolume(method string, url string, q url.Values, body []byte) (*easyhttp.Response, error) {
	if err := c.makeSureEndpoint("cinderv2", "v2"); err != nil {
		return nil, err
	}
	return c.doProxy2(c.serviceEndpoint["cinderv2"], method, url, q, body)
}

func getImageIdFromProxyUrl(proxyUrl string) (string, error) {
	reg := regexp.MustCompile("/images/(.*)/file")
	matched := reg.FindStringSubmatch(proxyUrl)
	logging.Debug("matched %v", matched)
	if len(matched) >= 2 {
		return matched[1], nil
	}
	return "", fmt.Errorf("get image id failed")
}

func (c *OpenstackManager) uploadImage(proxyUrl string, req *ghttp.Request) (*easyhttp.Response, error) {
	imageId, err := getImageIdFromProxyUrl(proxyUrl)
	if err != nil {
		return nil, err
	}
	dataPath, _ := g.Cfg().Get(gctx.New(), "server.dataPath")
	cacheFile := filepath.Join(dataPath.String(), "image_cache", imageId)

	// 上传到后端
	go func(imageId string, imageFile string) {
		imageBuff, err := ImageUploadBufReader(imageFile)
		if err != nil {
			glog.Errorf(req.GetCtx(), "load image from file failed: %s", err)
			return
		}
		glog.Infof(req.GetCtx(), "uploading image %s to backend, path: %s", imageId, imageFile)
		_, err = c.doProxyWithHeaders2BodyReader(
			c.serviceEndpoint["glance"], req.Method, proxyUrl, req.URL.Query(),
			map[string]string{
				"content-type":      easyhttp.APPLICATION_OCTET_STREAM,
				"x-image-meta-size": req.Header.Get("x-image-meta-size"),
			},
			imageBuff,
		)
		if err != nil {
			glog.Errorf(req.GetCtx(), "upload image %s failed: %s", imageId, err)
		} else {
			glog.Infof(req.GetCtx(), "uploaded image %s to backend", imageId)
			if err := os.Remove(imageFile); err != nil {
				glog.Warningf(req.GetCtx(), "remove image file %s failed: %s", imageFile, err)
			}
		}
	}(imageId, cacheFile)
	return nil, nil
}
func (c *OpenstackManager) SaveImageCache(proxyUrl string, req *ghttp.Request) (*entity.ImageUploadTask, error) {
	metaSize := req.Header.Get("x-image-meta-size")
	imageSize, err := strconv.Atoi(metaSize)
	if err != nil {
		return nil, fmt.Errorf("image size not found in body")
	}
	imageId, err := getImageIdFromProxyUrl(proxyUrl)
	if err != nil {
		return nil, err
	}
	dataPath, _ := g.Cfg().Get(gctx.New(), "server.dataPath")
	cacheFile := filepath.Join(dataPath.String(), "image_cache", imageId)
	if !gfile.Exists(cacheFile) {
		projectId, err := GetSessionProjectId(req)
		if err != nil {
			return nil, fmt.Errorf("get session project id failed: %s", err)
		}
		err = service.ImageUploadTaskService.Create(projectId, imageId, imageSize)
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
	err = service.ImageUploadTaskService.IncrementCached(imageId, len(data))
	if err != nil {
		return nil, fmt.Errorf("increment %s cached failed: %s", imageId, err)
	}
	task, err := service.ImageUploadTaskService.GetByImageId(imageId)
	if err != nil {
		return nil, fmt.Errorf("get task for %s failed: %s", imageId, err)
	}
	return task, err
}
func (c *OpenstackManager) ProxyImage(proxyUrl string, req *ghttp.Request) (*easyhttp.Response, error) {
	if err := c.makeSureEndpoint("glance", "v2"); err != nil {
		return nil, err
	}
	headers := map[string]string{}
	uploadFileReg, _ := regexp.Compile("/images/.+/file")
	if strings.ToUpper(req.Method) == "PUT" && uploadFileReg.MatchString(proxyUrl) {
		return c.uploadImage(proxyUrl, req)
	}
	if req.Header.Get("content-type") != "" {
		headers["content-type"] = req.Header.Get("content-type")
	}
	return c.doProxyWithHeaders2(
		c.serviceEndpoint["glance"], req.Method, proxyUrl, req.URL.Query(),
		headers, req.GetBody(),
	)
}

var SESSION_MANAGERS map[string]*OpenstackManager

func GetAuthInfo(project, user, password string) Auth {
	return Auth{
		Scope: Scope{
			Project: Project{
				Name:   project,
				Domain: Domain{Name: "default"},
			},
		},
		Identity: Identity{
			Methods: []string{"password"},
			Password: Password{
				User: User{Name: user, Password: password, Domain: Domain{Name: "default"}},
			},
		},
	}
}
func HideTokenHeader(header http.Header) http.Header {
	safeHeader := http.Header{}
	for h, v := range header {
		if strings.ToLower(h) == "x-auth-token" || strings.ToLower(h) == "x-subject-token" {
			continue
		}
		safeHeader[h] = v
	}
	return safeHeader
}
func NewManager(sessionId string, authUrl, project, user, password string) (*OpenstackManager, error) {
	manager := &OpenstackManager{
		Region:   "RegionOne",
		AuthInfo: GetAuthInfo(project, user, password),
		session2: easyhttp.DefaultClient().SetDefaultContentType(easyhttp.APPLICATION_JSON).
			SetSafeHeader(HideTokenHeader),
		tokenAlive:      time.Minute * 30,
		serviceEndpoint: map[string]string{},
		microVesrion:    map[string]string{},
	}
	manager.SetAuthUrl(authUrl)
	if err := manager.tokenIssue(); err != nil {
		return nil, err
	}
	SESSION_MANAGERS[sessionId] = manager
	return SESSION_MANAGERS[sessionId], nil
}
func GetAuthFromSession(req *ghttp.Request) (*LoginInfo, error) {
	sessionLoginInfo, _ := req.Session.Get("loginInfo", nil)
	loginInfo := LoginInfo{}
	if err := sessionLoginInfo.Struct(&loginInfo); err != nil {
		return nil, fmt.Errorf("get login info failed: %s", err)
	}
	return &loginInfo, nil
}
func GetSessionProjectId(req *ghttp.Request) (string, error) {
	loginInfo, err := GetAuthFromSession(req)
	if err != nil {
		return "", err
	}
	return loginInfo.Project.Id, err
}

func GetManager(sessionId string, req *ghttp.Request) (*OpenstackManager, error) {
	if client, ok := SESSION_MANAGERS[sessionId]; ok {
		return client, nil
	}
	if loginInfo, err := GetAuthFromSession(req); err != nil {
		return nil, fmt.Errorf("get session auth info falied: %v", err)
	} else {
		cluster, err := service.ClusterService.GetClusterByName(loginInfo.Cluster)
		if err != nil {
			return nil, fmt.Errorf("get cluster %s failed: %s", loginInfo.Cluster, err)
		}
		return NewManager(
			sessionId, cluster.AuthUrl,
			loginInfo.Project.Name, loginInfo.User.Name,
			loginInfo.Password)
	}
}

func init() {
	SESSION_MANAGERS = map[string]*OpenstackManager{}
}
