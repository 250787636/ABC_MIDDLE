package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ResponseJson struct {
	Raw   string `json:"-"`
	State int    `json:"state"`
	Msg   string `json:"msg"`
}

func (r *ResponseJson) Handle(resp *http.Response) error {
	byteArr, err := readHttpResponse(resp)
	if err != nil {
		return err
	}
	r.Raw = string(byteArr)
	if err = json.Unmarshal(byteArr, &r); err != nil {
		return fmt.Errorf("json parsing error:%w,data:%s", err, string(byteArr))
	}
	if r.State != 200 {
		return fmt.Errorf("接口响应错误:state:%d,msg:%s", r.State, r.Msg)
	}
	return nil
}
func (r *ResponseJson) Unmarshal(v interface{}) error {
	return json.Unmarshal([]byte(r.Raw), v)
}

type ResponseFile struct {
	LocalFilePath string
}

func (r *ResponseFile) Handle(resp *http.Response) error {
	if r.LocalFilePath == "" {
		return fmt.Errorf("文件保存路径不能为空")
	}
	byteArr, err := readHttpResponse(resp)
	if err != nil {
		return err
	}
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		return fmt.Errorf("下载文件失败,响应内容:%s", string(byteArr))
	}
	_ = os.MkdirAll(filepath.Dir(r.LocalFilePath), os.ModePerm)
	if err = ioutil.WriteFile(r.LocalFilePath, byteArr, os.ModePerm); err != nil {
		return fmt.Errorf("文件写入失败:%w", err)
	}
	return nil
}

func readHttpResponse(resp *http.Response) ([]byte, error) {
	defer func() { _ = resp.Body.Close() }()
	byteArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error:%w", err)
	}
	if len(byteArr) == 0 {
		return nil, fmt.Errorf("响应数据为空,StatusCode:%d", resp.StatusCode)
	}
	return byteArr, nil
}
