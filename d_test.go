package main

import (
	"example.com/m/pkg/log"
	"example.com/m/utils/response"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func Test(t *testing.T) {

}

func aaa(r *gin.Context) {
	token := r.GetHeader("token")
	url := ip + "/dmcwebapi/api/dmc/ReqToken"
	contentType := "application/json"
	dataMap := make(map[string]interface{})
	dataMap["type"] = "sgtoken"
	dataMap["imtoken"] = token

	jsonStr, err := jsoniter.Marshal(dataMap)
	resp, err := http.Post(url, contentType, strings.NewReader(string(jsonStr)))
	if err != nil {
		log.Error("调用验证token接口失败 error:", err)
		response.FailWithMessage("调用验证token接口失败", r)
		r.Abort()
		return
	}
	defer resp.Body.Close()

	post, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("读取验证token接口内容失败 error:", err.Error())
		response.FailWithMessage("读取验证token接口内容失败", r)
		r.Abort()
		return
	}

	var TokenResponse struct {
		Status     string `json:"status"`
		FailureMsg string `json:"failure_msg"`
		Uid        string `json:"uid"`
	}
	err = jsoniter.Unmarshal(post, &TokenResponse)
	if err != nil {
		log.Error("解析接口值失败", err.Error())
		response.FailWithMessage("解析接口值失败", r)
		r.Abort()
	}
}
