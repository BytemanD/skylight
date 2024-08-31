package easyhttp

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type Response struct {
	body        []byte
	rawResponse *http.Response
}

func (resp *Response) StatusCode() int {
	return resp.rawResponse.StatusCode
}

func (resp *Response) Status() string {
	return resp.rawResponse.Status
}
func (resp *Response) GetHeaders() http.Header {
	if resp.rawResponse == nil {
		return http.Header{}
	}
	return resp.rawResponse.Header
}
func (resp *Response) GetHeader(key string) string {
	return resp.GetHeaders().Get(key)
}
func (resp *Response) GetContentType() string {
	return resp.GetHeaders().Get(CONTENT_TYPE)
}
func (resp *Response) IsError() bool {
	return resp.StatusCode() >= 400
}
func (resp *Response) GetBody() ([]byte, error) {
	if resp.rawResponse == nil {
		return []byte{}, nil
	}
	var err error
	if resp.body == nil {
		resp.body, err = io.ReadAll(resp.rawResponse.Body)
		defer resp.rawResponse.Body.Close()
	}
	return resp.body, err
}

func (resp *Response) Body() []byte {
	body, _ := resp.GetBody()
	return body
}

func (resp *Response) SaveBody(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.rawResponse.Body)
	return err
}

func (resp *Response) UNmarshal(v interface{}) error {

	body, err := resp.GetBody()

	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
