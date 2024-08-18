package openstack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"skylight/internal/model"
	"skylight/internal/service"
	"strings"
	"time"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/go-resty/resty/v2"
	"github.com/gogf/gf/v2/net/ghttp"
)

type OpenstackManager struct {
	AuthUrl         string
	Region          string
	AuthInfo        Auth
	session         *resty.Client
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

func (c *OpenstackManager) sendToBackend(req *resty.Request) (*resty.Response, error) {
	logging.Debug("-----proxy---> %s %s ?%s Headers: %s",
		req.Method, req.URL, req.QueryParam.Encode(), c.safeHeader(req.Header))
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}
	proxyRespBody := "<Body>"
	if resp.Header().Get("Content-Type") == "application/json" {
		if resp.IsError() {
			proxyRespBody = string(resp.Body())
			// err = fmt.Errorf("reqeust failed: [%d] %s", resp.StatusCode(), resp.Body())
		}
	}
	logging.Debug("proxy Resp [%d] %s", resp.StatusCode(), proxyRespBody)
	return resp, err
}

func (c *OpenstackManager) tokenIssue() error {
	req := c.session.NewRequest()

	req.SetBody(map[string]Auth{"auth": c.AuthInfo})
	req.Method = resty.MethodPost
	if reqUrl, err := url.JoinPath(c.AuthUrl, "/auth/tokens"); err != nil {
		return err
	} else {
		req.URL = reqUrl
	}
	resp, err := c.sendToBackend(req)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("login success: [%d] %s", resp.StatusCode(), resp.Body())
	}
	c.token = resp.Header().Get("X-Subject-Token")
	respBody := struct{ Token TokenBody }{}
	if err := json.Unmarshal(resp.Body(), &respBody); err != nil {
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
func (c *OpenstackManager) doProxy(endpoint string, method string,
	u string, q url.Values, body []byte) (*resty.Response, error) {
	return c.doProxyWithHeaders(endpoint, method, u, q, nil, body)
}
func (c *OpenstackManager) doProxyWithHeaders(endpoint string, method string,
	u string, q url.Values, headers map[string]string, body []byte) (*resty.Response, error) {
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
	req := c.session.NewRequest().SetHeader("X-Auth-Token", token).
		SetQueryParamsFromValues(q).
		SetBody(body)
	if headers != nil {
		req = req.SetHeaders(headers)
	}
	req.Method, req.URL = method, reqUrl
	return c.sendToBackend(req)
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
func (c *OpenstackManager) ProxyIdentity(method string, url string, q url.Values, body []byte) (*resty.Response, error) {
	if err := c.makeSureEndpoint("identity", "v3"); err != nil {
		return nil, err
	}
	return c.doProxy(c.serviceEndpoint["identity"], method, url, q, body)
}
func (c *OpenstackManager) ProxyNetworking(method string, url string, q url.Values, body []byte) (*resty.Response, error) {
	if err := c.makeSureEndpoint("neutron", "v2.0"); err != nil {
		return nil, err
	}
	return c.doProxy(c.serviceEndpoint["neutron"], method, url, q, body)
}

func (c *OpenstackManager) getMicroVersion(service string) (string, error) {
	resp, err := c.doProxy(c.serviceEndpoint[service], resty.MethodGet, "/", nil, nil)
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
func (c *OpenstackManager) ProxyComputing(method string, url string, q url.Values, body []byte) (*resty.Response, error) {
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
		return c.doProxyWithHeaders(c.serviceEndpoint["nova"], method, url, q, headers, body)
	} else {
		return c.doProxy(c.serviceEndpoint["nova"], method, url, q, body)
	}
}
func (c *OpenstackManager) ProxyVolume(method string, url string, q url.Values, body []byte) (*resty.Response, error) {
	if err := c.makeSureEndpoint("cinderv2", "v2"); err != nil {
		return nil, err
	}
	return c.doProxy(c.serviceEndpoint["cinderv2"], method, url, q, body)
}
func (c *OpenstackManager) ProxyImage(proxyUrl string, req *ghttp.Request) (*resty.Response, error) {
	if err := c.makeSureEndpoint("glance", "v2"); err != nil {
		return nil, err
	}
	headers := map[string]string{}
	if req.Header.Get("content-type") != "" {
		headers["content-type"] = req.Header.Get("content-type")
	}
	return c.doProxyWithHeaders(
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
func NewManager(sessionId string, authUrl, project, user, password string) (*OpenstackManager, error) {
	manager := &OpenstackManager{
		Region:          "RegionOne",
		AuthInfo:        GetAuthInfo(project, user, password),
		session:         resty.New(),
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
func GetManager(sessionId string, req *ghttp.Request) (*OpenstackManager, error) {
	if client, ok := SESSION_MANAGERS[sessionId]; ok {
		return client, nil
	}
	if loginInfo, err := GetAuthFromSession(req); err != nil {
		return nil, fmt.Errorf("get session auth info falied: %v", err)
	} else {
		cluster, err := service.GetClusterByName(loginInfo.Cluster)
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
