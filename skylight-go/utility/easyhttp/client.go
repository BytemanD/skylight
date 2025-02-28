package easyhttp

import (
	"context"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"golang.org/x/net/publicsuffix"
)

type Client struct {
	DefaultContentType string

	httpClient *http.Client
	SafeHeader func(header http.Header) http.Header
}

func (c *Client) getSafeHeader(h http.Header) http.Header {
	if c.SafeHeader == nil {
		return h
	}
	return c.SafeHeader(h)
}
func (c *Client) SetSafeHeader(f func(h http.Header) http.Header) *Client {
	c.SafeHeader = f
	return c
}
func (c *Client) SetDefaultContentType(value string) *Client {
	c.DefaultContentType = value
	return c
}
func (c *Client) Execute(req *Request) (*Response, error) {
	rawReq := req.getRawRequest()
	if req.GetContentType() == APPLICATION_OCTET_STREAM {
		logging.Debug("easyhttp Req: %s %s\n    Headers: %s\n    Body: <steam>", rawReq.Method, rawReq.URL,
			c.getSafeHeader(req.Header),
		)
	} else {
		logging.Debug("easyhttp Req: %s %s\n    Headers: %s\n    Body: %v", rawReq.Method, rawReq.URL,
			c.getSafeHeader(req.Header), req.body,
		)
	}
	rawResp, err := c.httpClient.Do(rawReq)
	if err != nil {
		return nil, err
	}
	resp := Response{rawResponse: rawResp}
	if resp.IsError() {
		logging.Error("easyhttp Resp: [%s]    ContentLength: %d", resp.Status(), resp.rawResponse.ContentLength)
	} else if resp.GetContentType() == APPLICATION_OCTET_STREAM {
		logging.Debug("easyhttp Resp: [%s]    Body: <stream>", resp.Status())
	} else {
		logging.Debug("easyhttp Resp: [%s]    ContentLength: %d", resp.Status(), resp.rawResponse.ContentLength)
	}
	return &resp, nil
}

func (c *Client) NewRequest() *Request {
	req := Request{
		Header:      http.Header{},
		QueryValues: url.Values{},
		client:      c,
	}
	if c.DefaultContentType != "" {
		req.SetContentType(c.DefaultContentType)
	}
	return &req
}

func New() *Client {
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	return &Client{
		httpClient: &http.Client{
			Jar: cookieJar,
		},
	}
}
func transportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
func DefaultClient() *Client {
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	rawHttpClient := http.Client{
		Jar: cookieJar,
		Transport: &http.Transport{
			DialContext:           transportDialContext(dialer),
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			// MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}
	return &Client{httpClient: &rawHttpClient}
}
