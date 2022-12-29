package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type EvaluationClient struct {
	HostUrl   string
	Username  string
	Password  string
	Token     string
	Signature string
	lock      sync.RWMutex
}

type ApiResponseHandler interface {
	Handle(resp *http.Response) error
}

type UnmarshalData interface {
	Unmarshal(v interface{}) error
}

func NewClient(host string, user string, pwd string) *EvaluationClient {
	client := new(EvaluationClient)
	client.HostUrl = host + "/v4"
	client.Username = user
	client.Password = pwd
	return client
}

func (c *EvaluationClient) SendFormDataRequest(url string, params map[string]string, fileParam ...map[string]string) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for filed, val := range params {
		err := writer.WriteField(filed, val)
		if err != nil {
			return nil, err
		}
	}
	if len(fileParam) > 0 && fileParam[0] != nil {
		for fileField, filePath := range fileParam[0] {
			_, filename := filepath.Split(filePath)
			part, err := writer.CreateFormFile(fileField, filename)
			if err != nil {
				return nil, err
			}
			file, err := os.Open(filePath)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return nil, err
			}
		}
	}
	_ = writer.Close()
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, nil
}

func (c *EvaluationClient) requestApi(api string, params map[string]interface{}, fileParam ...map[string]string) (resp *http.Response, err error) {
	if params == nil {
		params = make(map[string]interface{}, 2)
	}
	switch api {
	case "/apply_auth", "/apply_access_token":
	default:
		c.lock.RLock()
		params["token"] = c.Token
		params["signature"] = c.Signature
		c.lock.RUnlock()
	}
	url := c.HostUrl + api
	strBytes, _ := json.Marshal(params)
	req, err := c.SendFormDataRequest(url, map[string]string{
		"param": string(strBytes),
	}, fileParam...)
	if err != nil {
		return nil, fmt.Errorf("构造请求%s错误:%w", api, err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求接口%s错误:%w", api, err)
	}
	return resp, nil
}

func (c *EvaluationClient) RequestApi(api string, handler ApiResponseHandler, params map[string]interface{}, fileParam ...map[string]string) error {
	resp, err := c.requestApi(api, params, fileParam...)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if err = handler.Handle(resp); err != nil {
		return fmt.Errorf("请求%s:%w", api, err)
	}
	if unmarshalHandler, ok := handler.(UnmarshalData); ok {
		if err = unmarshalHandler.Unmarshal(handler); err != nil {
			return fmt.Errorf("UnmarshalData:%w", err)
		}
	}
	return nil
}

// 获取预览地址接口
func (c *EvaluationClient) SdkSearchOneDetail(handler ApiResponseHandler, taskId int) error {
	return c.RequestApi("/sdk/search_one_detail", handler, map[string]interface{}{
		"taskid": taskId,
	})
}

// 获取任务进度接口
func (c *EvaluationClient) TaskProgress(handler ApiResponseHandler, taskType string, taskId int) error {
	return c.RequestApi("/task/progress", handler, map[string]interface{}{
		"type":   taskType,
		"taskid": taskId,
	})
}
