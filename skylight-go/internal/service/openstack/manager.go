package openstack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/go-resty/resty/v2"
	"github.com/gogf/gf/v2/net/ghttp"
)

type OpenstackManager struct {
	AuthUrl         string
	AuthInfo        Auth
	session         *resty.Client
	token           string
	tokenAlive      time.Duration
	expiredAt       time.Time
	catalogs        []Catalog
	serviceEndpoint map[string]string
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
	logging.Debug("proxy %s %s ?%s Headers: %s",
		req.Method, req.URL, req.QueryParam.Encode(), c.safeHeader(req.Header))
	resp, err := req.Send()

	proxyRespBody := "<Body>"
	if resp.Header().Get("Content-Type") == "application/json" {
		if resp.IsError() {
			proxyRespBody = string(resp.Body())
			if resp.IsError() {
				err = fmt.Errorf("reqeust failed: [%d] %s", resp.StatusCode(), resp.Body())
			}
		}
	}
	logging.Debug("proxy Resp [%d] %s", resp.StatusCode(), proxyRespBody)
	return resp, err
}

func (c *OpenstackManager) tokenIssue() error {
	req := c.session.NewRequest()

	req.SetBody(map[string]Auth{"auth": c.AuthInfo})
	req.Method = resty.MethodPost
	// service.GetClusterByName(name)
	fmt.Println("222222222222 auth url", c.AuthUrl)

	req.URL, _ = url.JoinPath(c.AuthUrl, "/auth/tokens")
	resp, err := c.sendToBackend(req)
	if err != nil {
		return err
	}
	c.token = resp.Header().Get("X-Subject-Token")
	respBody := struct{ Token TokenBody }{}
	if err := json.Unmarshal(resp.Body(), &respBody); err != nil {
		return err
	}
	c.catalogs = respBody.Token.Catalogs
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
	for _, catalog := range c.catalogs {
		if catalog.Type != service && catalog.Name != service {
			continue
		}
		for _, endpoint := range catalog.Endpoints {
			if endpoint.Interface == "public" {
				return endpoint.Url, nil
			}
		}
	}
	return "", fmt.Errorf("endpoint for service %s not found", service)
}
func (c *OpenstackManager) makeSureEndpoint(service, version string) (err error) {
	if c.serviceEndpoint == nil {
		c.serviceEndpoint = map[string]string{}
	}
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
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	reqUrl, err := url.JoinPath(endpoint, u)
	if err != nil {
		return nil, err
	}
	req := c.session.NewRequest().SetHeader("X-Auth-Token", token).
		SetQueryParamsFromValues(q).
		SetBody(body)
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
func (c *OpenstackManager) ProxyComputing(method string, url string, q url.Values, body []byte) (*resty.Response, error) {
	if err := c.makeSureEndpoint("nova", "v2.1"); err != nil {
		return nil, err
	}
	return c.doProxy(c.serviceEndpoint["nova"], method, url, q, body)
}
func (c *OpenstackManager) ProxyVolume(method string, url string, q url.Values, body []byte) (*resty.Response, error) {
	if err := c.makeSureEndpoint("cinderv2", "v2"); err != nil {
		return nil, err
	}
	return c.doProxy(c.serviceEndpoint["cinderv2"], method, url, q, body)
}
func (c *OpenstackManager) ProxyImage(method string, url string, q url.Values, body []byte) (*resty.Response, error) {
	if err := c.makeSureEndpoint("glance", "v2"); err != nil {
		return nil, err
	}
	return c.doProxy(c.serviceEndpoint["glance"], method, url, q, body)
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
	if client, ok := SESSION_MANAGERS[sessionId]; ok {
		return client, nil
	} else {
		manager := &OpenstackManager{
			AuthInfo:   GetAuthInfo(project, user, password),
			session:    resty.New(),
			tokenAlive: time.Minute * 30,
		}
		manager.SetAuthUrl(authUrl)
		if err := manager.tokenIssue(); err != nil {
			return nil, err
		}
		SESSION_MANAGERS[sessionId] = manager
	}
	return SESSION_MANAGERS[sessionId], nil
}
func GetManager(sessionId string, req *ghttp.Request) (*OpenstackManager, error) {
	if client, ok := SESSION_MANAGERS[sessionId]; ok {
		return client, nil
	} else {
		authUrl, _ := req.Session.Get("authUrl", nil)
		project, _ := req.Session.Get("project", nil)
		user, _ := req.Session.Get("user", nil)
		password, _ := req.Session.Get("password", nil)
		if authUrl == nil || project == nil || user == nil || password == nil {
			return nil, fmt.Errorf("auth info not found")
		}
		manager, err := NewManager(sessionId, authUrl.String(), project.String(),
			user.String(), password.String())
		if err != nil {
			return nil, fmt.Errorf("create manager failed")
		} else {
			return manager, nil
		}
	}
}

func init() {
	SESSION_MANAGERS = map[string]*OpenstackManager{}
}
