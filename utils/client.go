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
		return nil, fmt.Errorf("????????????%s??????:%w", api, err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("????????????%s??????:%w", api, err)
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
		return fmt.Errorf("??????%s:%w", api, err)
	}
	if unmarshalHandler, ok := handler.(UnmarshalData); ok {
		if err = unmarshalHandler.Unmarshal(handler); err != nil {
			return fmt.Errorf("UnmarshalData:%w", err)
		}
	}
	return nil
}

// ????????????????????????
func (c *EvaluationClient) SdkSearchOneDetail(handler ApiResponseHandler, taskId int) error {
	return c.RequestApi("/sdk/search_one_detail", handler, map[string]interface{}{
		"taskid": taskId,
	})
}

// ????????????????????????  Android ios sdk ???????????????
func (c *EvaluationClient) TaskProgress(handler ApiResponseHandler, taskType string, taskId int) error {
	return c.RequestApi("/task/progress", handler, map[string]interface{}{
		"type":   taskType,
		"taskid": taskId,
	})
}

// BinCheckApk 2.1 ??????apk?????????????????????  android
func (c *EvaluationClient) BinCheckApkClient(handler ApiResponseHandler, items, callbackUrl string, fileParam map[string]string) error {
	return c.RequestApi("/bin_check_apk", handler, map[string]interface{}{
		"items":        items,
		"callback_url": callbackUrl,
	}, fileParam)
}

// ??????android???????????????????????????
func (c *EvaluationClient) GetAdInfoClient(handler ApiResponseHandler, taskId int) error {
	return c.RequestApi("/ad/preview", handler, map[string]interface{}{
		"taskid": taskId,
	})
}

// SearchOneDetail 2.5??????????????????apk???????????????
func (c *EvaluationClient) SearchOneDetailClient(handler ApiResponseHandler, taskId int) error {
	return c.RequestApi("/search_one_detail", handler, map[string]interface{}{
		"taskid": taskId,
	})
}

// IosBinCheck 2.1 ??????ios?????????????????????  ios
func (c *EvaluationClient) IosBinCheckClient(handler ApiResponseHandler, items, callbackUrl string, fileParam map[string]string) error {
	return c.RequestApi("/ios/bin_check_apk", handler, map[string]interface{}{
		"items":        items,
		"callback_url": callbackUrl,
	}, fileParam)
}

// ??????ipa???????????????????????????
func (c *EvaluationClient) GetIpaInfoClient(handler ApiResponseHandler, taskId int) error {
	return c.RequestApi("/ios/preview", handler, map[string]interface{}{
		"taskid": taskId,
	})
}

// IosSearchOneDetail 3.2.??????ipa???????????????????????????
func (c *EvaluationClient) IosSearchOneDetailClient(handler ApiResponseHandler, taskId int) error {
	return c.RequestApi("/ios/search_one_detail", handler, map[string]interface{}{
		"taskid": taskId,
	})
}

// SdkBinCheck 14.1.??????sdk?????????????????????
func (c *EvaluationClient) SdkBinCheckClient(handler ApiResponseHandler, items, callbackUrl string, fileParam map[string]string) error {
	return c.RequestApi("/sdk/bin_check_apk", handler, map[string]interface{}{
		"items":        items,
		"callback_url": callbackUrl,
	}, fileParam)
}

// SdkSearchOneDetail 14.2.??????sdk???????????????????????????
func (c *EvaluationClient) SdkSearchOneDetailClient(handler ApiResponseHandler, taskId int) error {
	return c.RequestApi("/sdk/search_one_detail", handler, map[string]interface{}{
		"taskid": taskId,
	})
}
