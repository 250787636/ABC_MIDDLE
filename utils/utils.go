package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

// 加固内置账号使用的数据
var APIKEY, APISECRET = JiaGuApi()
var H5APIKEY, H5APISECRET = H5JiaGuApi()

// 加固内置账户
func JiaGuApi() (string, string) {
	user := app.Conf.AndroidJiaGu
	return user.ApiKey, user.ApiSecret
}

// h5加固内置账号
func H5JiaGuApi() (string, string) {
	user := app.Conf.H5JiaGu
	return user.ApiKey, user.ApiSecret
}

// 获取MD5sting
func MD5String(data []byte) string {
	return hex.EncodeToString(MD5(data))
}

// 获取MD5[]byte
func MD5(data []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return md5Ctx.Sum(nil)
}

func NewFormDataRequest(url string, params map[string]interface{}, fileParams map[string]interface{}) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	// 创建 string类型
	for filed, val := range params {
		switch v := val.(type) {
		case string:
			err := writer.WriteField(filed, v)
			if err != nil {
				return nil, err
			}
		}
	}

	// 创建file类型
	for filed, val := range fileParams {
		switch v := val.(type) {
		case *multipart.FileHeader:
			part, err := writer.CreateFormFile(filed, v.Filename)
			if err != nil {
				return nil, err
			}
			file, err := v.Open()
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return nil, err
			}
		case *os.File:
			part, err := writer.CreateFormFile(filed, path.Base(v.Name()))
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, v)
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

// 获取分页的 page size offNum
func GetFromDataPageSizeOffNum(c *gin.Context) (int, int, int) {
	num, pageOk := c.GetPostForm("page")  // 页数
	num2, sizeOK := c.GetPostForm("size") // 每页条数
	if !(pageOk && sizeOK) {
		response.FailResult(http.StatusInternalServerError, "", model.ReqParameterMissing, c)
		return -1, -1, -1
	}
	page, _ := strconv.Atoi(num)
	size, _ := strconv.Atoi(num2)
	// 跳过条数
	offNum := size * (page - 1)
	return page, size, offNum
}

// 判断数据是否为空字符
func IsStringEmpty(data string, isExist bool) bool {
	if isExist && data == "" {
		isExist = false
	}
	return isExist
}

// string 数组去重
func RemoveDuplicatesAndEmpty(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// 进行加固平台特有的验签
func CheckBoxSign(m map[string]interface{}, c *gin.Context) error {
	skyApiKey, sign := GetHeaderData(c)
	if skyApiKey == "" || sign == "" {
		return errors.New("签名文件参数缺失")
	}
	// 进行本机加密
	res := GetSignByApiKey(skyApiKey, m)
	if res != sign {
		return errors.New("验签失败,参数或者api_key有误")
	}
	return nil
}

// 请求头中获取验签数据并进行返回
func GetHeaderData(c *gin.Context) (string, string) {
	apiKey := c.GetHeader("api_key")
	sign := c.GetHeader("sign")
	return apiKey, sign
}

// 参数验签
func CheckSign(apiKey, sign string, m map[string]interface{}) error {
	res := GetSignByApiKey(apiKey, m)
	if res == sign {
		return nil
	} else {
		return errors.New("验签失败,参数或者api_key有误")
	}
}

// 获取配置文件存储的 密钥获取签名(sign
func GetSign(m map[string]interface{}) (result string) {
	result = hmacSha1(APISECRET, concatParam(m, APIKEY))
	return result
}

// h5加固获取配置文件存储的 密钥获取签名(sign
func H5GetSign(m map[string]interface{}) (result string) {
	result = hmacSha1(H5APISECRET, concatParam(m, H5APIKEY))
	return result
}

// 通过远端给予的 密钥获取签名
func GetSignByApiKey(apiKey string, m map[string]interface{}) (result string) {
	result = hmacSha1(APISECRET, concatParam(m, apiKey))
	return result
}

// 获取sign的两个工具方法
func hmacSha1(secret, text string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(text))
	return hex.EncodeToString(mac.Sum(nil))
}
func concatParam(m map[string]interface{}, apiKey string) string {
	result := apiKey
	keyList := make([]string, 0)
	for k, _ := range m {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	for _, k := range keyList {
		result = result + fmt.Sprintf("%v", m[k])
	}
	return strings.Trim(result, "&")
}
