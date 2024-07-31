package openstack

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/go-resty/resty/v2"
)

type OpenstackManager struct {
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
	logging.Debug("proxy GET %s\n    Headers: %s", req.URL, req.Header)
	resp, err := req.Send()

	proxyRespBody := "<...>"
	if resp.Header().Get("Content-Type") == "application/json" {
		if resp.IsError() {
			proxyRespBody = string(resp.Body())
		}
	}
	logging.Debug("proxy Resp [%d]:\n     %s", resp.StatusCode(), proxyRespBody)
	return resp, err
}

func (c *OpenstackManager) tokenIssue() error {
	req := c.session.NewRequest()

	req.SetBody(map[string]Auth{"auth": getAuth()})
	req.Method = resty.MethodPost
	req.URL, _ = url.JoinPath("http://keystone-server:35357/v3", "/auth/tokens")
	resp, err := c.sendToBackend(req)
	if err != nil {
		return nil
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
		if strings.HasSuffix(endpoint, version) {
			c.serviceEndpoint[service], err = endpoint, nil
		} else {
			c.serviceEndpoint[service], err = url.JoinPath(endpoint, version)
		}
	}
	return
}

func (c *OpenstackManager) doProxy(endpoint string, u string) (*resty.Response, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}
	req := c.session.NewRequest().SetHeader("X-Auth-Token", token)
	reqUrl, err := url.JoinPath(endpoint, u)
	if err != nil {
		return nil, err
	}
	logging.Debug("proxy GET %s\n    Headers: %s", reqUrl, req.Header)

	resp, err := req.Get(reqUrl)
	proxyRespBody := "<...>"
	if err != nil {
		return nil, err
	}
	if resp.Header().Get("Content-Type") == "application/json" {
		if resp.IsError() {
			proxyRespBody = string(resp.Body())
		}
	}
	logging.Debug("proxy Resp [%d]:\n     %s", resp.StatusCode(), proxyRespBody)
	return resp, nil
}
func (c *OpenstackManager) ProxyNetworking(u string) (*resty.Response, error) {
	if err := c.makeSureEndpoint("neutron", "v2.0"); err != nil {
		return nil, err
	}
	return c.doProxy(c.serviceEndpoint["neutron"], u)
}
func (c *OpenstackManager) ProxyComputing(u string) (*resty.Response, error) {
	if err := c.makeSureEndpoint("nova", "v2.1"); err != nil {
		return nil, err
	}
	return c.doProxy(c.serviceEndpoint["nova"], u)
}

var TOKEN_CACHE map[string]*OpenstackManager

func GetManager() *OpenstackManager {
	// TODO use cookie
	if client, ok := TOKEN_CACHE["cookie"]; ok {
		return client
	} else {
		TOKEN_CACHE["cookie"] = &OpenstackManager{
			session:    resty.New(),
			tokenAlive: time.Minute * 30,
		}
	}
	return TOKEN_CACHE["cookie"]
}

func init() {
	TOKEN_CACHE = map[string]*OpenstackManager{}
}
