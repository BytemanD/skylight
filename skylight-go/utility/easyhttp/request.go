package easyhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	CONTENT_TYPE             = "Content-Type"
	APPLICATION_JSON         = "application/json"
	APPLICATION_OCTET_STREAM = "application/octet-stream"
)

type Request struct {
	URL         string
	Method      string
	QueryValues url.Values
	Header      http.Header

	body   io.Reader
	client *Client
}

func (r *Request) SetMethod(method string) *Request {
	r.Method = method
	return r
}
func (r *Request) SetURL(url string) *Request {
	r.URL = url
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	r.Header.Set(key, value)
	return r
}
func (r *Request) AddHeader(key, value string) *Request {
	r.Header.Add(key, value)
	return r
}
func (r *Request) SetHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	return r
}
func (r *Request) AddHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.Header.Add(k, v)
	}
	return r
}
func (r *Request) SetContentType(value string) *Request {
	r.SetHeader(CONTENT_TYPE, value)
	return r
}
func (r *Request) GetContentType() string {
	return r.Header.Get(CONTENT_TYPE)
}

func (r *Request) SetQueryValues(key, value string) *Request {
	r.QueryValues.Set(key, value)
	return r
}
func (r *Request) AddQueryValues(key, value string) *Request {
	r.QueryValues.Add(key, value)
	return r
}
func (r *Request) AddQueryValuesFromValues(q url.Values) *Request {
	for key, values := range q {
		for _, value := range values {
			r.AddQueryValues(key, value)
		}
	}
	return r
}
func (r *Request) GetBody() io.Reader {
	return r.body
}
func (r *Request) GetBytesBody() ([]byte, error) {
	if r.body == nil {
		return []byte{}, nil
	}
	bytes := []byte{}

	_, err := r.body.Read(bytes)
	return bytes, err
}

func (r *Request) SetJsonBody(body interface{}) *Request {
	if data, ok := body.([]byte); ok {
		if len(data) > 0 {
			r.body = bytes.NewBuffer(data)
		}
	} else {
		v, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		r.body = bytes.NewBuffer(v)
	}
	return r
}
func (r *Request) SetStringBody(body string) *Request {
	if len(body) > 0 {
		r.body = bytes.NewBufferString(body)
	}
	return r
}
func (r *Request) SetFileBody(filePath string) *Request {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	r.body = file
	return r
}
func (r *Request) SetReaderBody(reader io.Reader) *Request {
	r.body = reader
	return r
}
func (r *Request) GetMethod() string {
	if r.Method != "" {
		return r.Method
	}
	return http.MethodGet
}
func (r *Request) getReqBody() io.Reader {
	// TODO
	if r.GetMethod() == http.MethodGet {
		return nil
	} else {
		return r.body
	}
}
func (r *Request) getRawRequest() *http.Request {
	rawReq, err := http.NewRequest(r.GetMethod(), r.URL, r.body)
	if err != nil {
		panic(err)
	}
	rawReq.URL.RawQuery = r.QueryValues.Encode()
	for k, v := range r.Header {
		rawReq.Header[k] = v
	}
	return rawReq
}
func (r *Request) HasHeder(key string) bool {
	_, ok := r.Header[key]
	return ok
}

func (r *Request) Send() (*Response, error) {
	return r.client.Execute(r)
}

func (r *Request) Get() (*Response, error) {
	return r.SetMethod(http.MethodGet).Send()
}

func (r *Request) Post() (*Response, error) {
	return r.SetMethod(http.MethodPost).Send()
}
func (r *Request) Put() (*Response, error) {
	return r.SetMethod(http.MethodPut).Send()
}
func (r *Request) Delete() (*Response, error) {
	return r.SetMethod(http.MethodDelete).Send()
}
func (r *Request) Patch() (*Response, error) {
	return r.SetMethod(http.MethodPatch).Send()
}
func (r *Request) Head() (*Response, error) {
	return r.SetMethod(http.MethodHead).Send()
}
func (r *Request) Options() (*Response, error) {
	return r.SetMethod(http.MethodOptions).Send()
}
