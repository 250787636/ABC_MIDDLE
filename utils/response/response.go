package response

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

type PageList struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

func CutPage(total int, page int, limit int) (pager PageList) {
	pager.Total = total
	pager.TotalPage = int(math.Ceil(float64(total) / float64(limit)))
	pager.Page = page
	pager.Limit = limit
	return
}

type Response struct {
	Code int         `json:"code"`
	Info interface{} `json:"info"`
	Msg  string      `json:"msg"`
}

type DataList struct {
	DataList interface{} `json:"datalist"`
	PageList
}

const (
	ERROR   = 500
	SUCCESS = 200
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", c)
}

func OkDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

// OkWithList 返回分页列表数据
func OkWithList(data interface{}, total, page, limit int, c *gin.Context) {
	Result(SUCCESS, DataList{data, CutPage(total, page, limit)}, "操作成功", c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	FailResult(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(code int, data interface{}, message string, c *gin.Context) {
	Result(code, data, message, c)
}

func FailResult(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, FailResponse{
		code,
		data,
		msg,
	})
}

type FailResponse struct {
	Code int         `json:"code"`
	Info interface{} `json:"info"`
	Err  string      `json:"err"`
}
