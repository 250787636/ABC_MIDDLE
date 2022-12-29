package ceping

import (
	"bytes"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MpBinCheckRequest struct {
	AppName     string `form:"file_name"` // 文件名称
	PkgName     string `form:"pkg_name"`  // 应用名称
	CallBackUrl string `form:"callback_url"`
	TemplateId  int    `form:"template_id"`
	FilePath    string `form:"file_path" binding:"required"`
}
type MpHandle struct {
	Locker  sync.Mutex
	TaskIds []int
}

var MpHandler = NewMpHandler()

func NewMpHandler() *MpHandle {
	hand := &MpHandle{
		TaskIds: make([]int, 0),
	}
	go func() {
		for {
			hand.Check(hand)
			time.Sleep(20 * time.Second)
		}
	}()
	return hand
}

func (h *MpHandle) Add(taskId int) {
	h.Locker.Lock()
	defer h.Locker.Unlock()
	h.TaskIds = append(h.TaskIds, taskId)
}
func (h *MpHandle) GetTaskIds() []int {
	h.Locker.Lock()
	defer h.Locker.Unlock()
	return h.TaskIds
}
func (h *MpHandle) RemoveTask(taskId int) {
	h.Locker.Lock()
	defer h.Locker.Unlock()
	for i, v := range h.TaskIds {
		if v == taskId {
			h.TaskIds = append(h.TaskIds[:i], h.TaskIds[i+1:]...)
			break
		}
	}
}
func (h *MpHandle) Check(hand *MpHandle) {
	if len(h.TaskIds) == 0 {
		return
	}
	for _, taskId := range h.GetTaskIds() {
		GetMpInfo(float64(taskId), hand)
	}
}

// MpBinCheck 7.1．提交小程序任务
func MpBinCheck(c *gin.Context) {
	var req = MpBinCheckRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}
	var templateInfo model.Template
	templateInfo.ID = uint(req.TemplateId)
	if err := app.DB.Model(&model.Template{}).First(&templateInfo).Error; err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage("未查询到当前模板", c)
		return
	}
	//fmt.Println("找出来的测评项", templateInfo.Items)
	// 1.1.读取文件
	open, err := os.Open(req.FilePath)
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	part, err := writer.CreateFormFile("file", path.Base(open.Name()))
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	_, err = io.Copy(part, open)
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	paramMap := make(map[string]interface{})
	paramMap["token"] = app.Conf.CePing.Token
	paramMap["signature"] = app.Conf.CePing.Signature
	paramMap["name"] = req.PkgName
	paramMap["items"] = templateInfo.Items
	paramMap["callback_url"] = req.CallBackUrl

	value, err := jsoniter.Marshal(paramMap)
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = writer.WriteField("param", string(value))
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = writer.Close()
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	clientURL := IP + "/v4/mp/bin_check"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	//注意别忘了设置header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	//Do方法发送请求
	resp, err := client.Do(request)
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	post, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	reponse := make(map[string]interface{})
	err = jsoniter.Unmarshal(post, &reponse)
	if err != nil {
		return
	}
	err = Check(reponse, post, MpBinCheck, c)
	if err != nil {
		log.Error("err:", err.Error())
		return
	}
	taskId, _ := reponse["task_id"].(string)
	userId, _ := c.Get("userId")
	userID := userId.(uint)
	fmt.Println("userid", userID)

	var info model.CePingUserTask
	info.TaskType = 3
	info.Status = "测评中"
	info.AppName = req.AppName
	info.PkgName = req.PkgName
	info.TemplateID = uint(req.TemplateId)
	info.TaskID = taskId
	info.UserID = userID
	info.FilePath = req.FilePath

	if err := app.DB.Model(&model.CePingUserTask{}).Create(&info).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	tId, _ := strconv.Atoi(taskId)
	MpHandler.Add(tId)
	response.OkWithData(reponse, c)
	return
}

// GetMpInfo 获取小程序任务状态
func GetMpInfo(taskId float64, hand *MpHandle) {
	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	paramMap := make(map[string]interface{})
	paramMap["token"] = app.Conf.CePing.Token
	paramMap["signature"] = app.Conf.CePing.Signature
	paramMap["taskid"] = taskId

	value, err := jsoniter.Marshal(paramMap)
	if err != nil {
		return
	}
	err = writer.WriteField("param", string(value))
	if err != nil {
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}

	clientURL := IP + "/v4/mp/search_mp"

	// 生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		return
	}

	// 注意别忘了设置header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	//Do方法发送请求
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	post, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	reponse := make(map[string]interface{})
	err = jsoniter.Unmarshal(post, &reponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	errMessage := ""
	if len(reponse) == 0 {
		errMessage = strings.Trim(string(post), `"`)
	} else if key, ok := reponse["state"].(float64); ok && key != 200 {
		errMessage = reponse["msg"].(string)
	}
	fmt.Println("err", errMessage)
	if errMessage == "签名验证失败" || errMessage == "token验证失败" {
		app.Conf = app.LoadConfig()
		GetMpInfo(taskId, hand)
		return
	}
	if errMessage != "" {
		//fmt.Println("app.conf", app.Conf.CePing.Signature)
		fmt.Println("err", errMessage)
	}
	//fmt.Println("date", reponse["data"])

	data := reponse["data"].(map[string]interface{})
	overview, _ := data["overview"].(map[string]interface{})
	finish_item, _ := overview["finish_item"].(float64)
	total_item, _ := overview["total_item"].(float64)
	//fmt.Println("risk", item_statistic["risk"])

	var dataInfo map[string]interface{}
	marshal, err := jsoniter.Marshal(reponse["data"])
	if err != nil {
		fmt.Println(err)
		return
	}

	err = jsoniter.Unmarshal(marshal, &dataInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	riskNum, err := jsoniter.Marshal(dataInfo["risk_statistic"])
	if err != nil {
		fmt.Println(err)
		return
	}

	var riskNumMap map[string]interface{}
	err = jsoniter.Unmarshal(riskNum, &riskNumMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("111", riskNumMap)
	high, _ := riskNumMap["high"].(float64)
	low, _ := riskNumMap["low"].(float64)
	medium, _ := riskNumMap["medium"].(float64)

	var finishInfo model.CePingUserTask
	finishInfo.FinishItem = int(finish_item)
	finishInfo.ItemsNum = int(total_item)
	if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).Updates(&finishInfo).Error; err != nil {
		return
	}

	var info model.CePingUserTask
	info.Status = "测评完成"
	times := time.Now()
	info.FinishedTime = &times
	info.HighNum = int(high)
	info.MiddleNum = int(medium)
	info.LowNum = int(low)

	var FindInfo model.CePingUserTask
	if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).First(&FindInfo).Error; err != nil {
		if err != nil {
			return
		}
	}
	if FindInfo.CreatedAt.Add(1*time.Hour).Unix() < time.Now().Unix() {
		hand.RemoveTask(int(taskId))
		FindInfo.Status = "测评失败"
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).Updates(&FindInfo).Error; err != nil {
			return
		}
		return
	}

	fmt.Println("finish_item", finish_item)
	fmt.Println("total_item", total_item)

	if total_item-finish_item == 1 {
		var taskInfo model.CePingUserTask
		taskInfo.Status = "测评报告生成中"
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).
			Updates(&taskInfo).Error; err != nil {
			fmt.Println("修改失败", err)
			return
		}
	}
	if finish_item == total_item {
		// 如果完成检测就 就写小程序详细数据
		fmt.Println("小程序测评完成开始获取详细信息")
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).Updates(&info).Error; err != nil {
			return
		}
		hand.RemoveTask(int(taskId))
		return
	}

}

type MpSearchRequest struct {
	TaskId int `form:"task_id" binding:"required"`
}

// MpSearch 7.2．查询小程序任务
func MpSearch(c *gin.Context) {
	var req = MpSearchRequest{}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	paramMap := make(map[string]interface{})
	paramMap["token"] = app.Conf.CePing.Token
	paramMap["signature"] = app.Conf.CePing.Signature
	paramMap["taskid"] = req.TaskId

	value, err := jsoniter.Marshal(paramMap)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"错误": err.Error(),
		})
		return
	}
	err = writer.WriteField("param", string(value))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"错误": err.Error(),
		})
		return
	}
	err = writer.Close()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"错误": err.Error(),
		})
		return
	}

	clientURL := IP + "/v4/mp/search_mp"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"错误": 1,
		})
		return
	}

	//注意别忘了设置header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	//Do方法发送请求
	resp, err := client.Do(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	post, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	reponse := make(map[string]interface{})
	err = jsoniter.Unmarshal(post, &reponse)
	if err != nil {
		return
	}
	err = Check(reponse, post, MpSearch, c)
	if err != nil {
		return
	}
	response.OkWithData(reponse, c)
}

// MpMiniReport 7.3．下载小程序word或pdf报告
func MpMiniReport(c *gin.Context) {
	num := c.Query("task_id")
	taskId, _ := strconv.Atoi(num)
	downloadType := c.Query("type")

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	paramMap := make(map[string]interface{})
	paramMap["token"] = app.Conf.CePing.Token
	paramMap["signature"] = app.Conf.CePing.Signature
	paramMap["taskid"] = taskId
	paramMap["type"] = downloadType

	value, err := jsoniter.Marshal(paramMap)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = writer.WriteField("param", string(value))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = writer.Close()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	clientURL := IP + "/v4/mp/mini_report"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	//注意别忘了设置header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	//Do方法发送请求
	resp, err := client.Do(request)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		reponse, _ := ioutil.ReadAll(resp.Body)
		responseMap := make(map[string]interface{})
		jsoniter.Unmarshal(reponse, &responseMap)
		err := Check(responseMap, reponse, MpMiniReport, c)
		if err != nil {
			return
		}
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	// 控制用户请求所得的内容存为一个文件的时候提供一个默认的文件名
	c.Writer.Header().Set("Content-Disposition", contentDisposition)
	_, _ = io.Copy(c.Writer, resp.Body)
}

type MpReBinCheckRequest struct {
	TaskId int `form:"task_id" binding:"required"`
}

// MpReBinCheck 7.4 重新扫描小程序任务
func MpReBinCheck(c *gin.Context) {
	var req = MpReBinCheckRequest{}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	paramMap := make(map[string]interface{})
	paramMap["token"] = app.Conf.CePing.Token
	paramMap["signature"] = app.Conf.CePing.Signature
	paramMap["taskid"] = req.TaskId

	value, err := jsoniter.Marshal(paramMap)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"错误": err.Error(),
		})
		return
	}
	err = writer.WriteField("param", string(value))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"错误": err.Error(),
		})
		return
	}
	err = writer.Close()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"错误": err.Error(),
		})
		return
	}

	clientURL := IP + "/v4/mp/re_bin_check"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"错误": 1,
		})
		return
	}

	//注意别忘了设置header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	//Do方法发送请求
	resp, err := client.Do(request)
	post, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"错误": err.Error(),
		})
		return
	}

	reponse, err := checkError(post)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"info": gin.H{
				"datalist": err.Error(),
			},
			"code": reponse["state"],
			"msg":  reponse["msg"],
		})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"info": gin.H{
				"datalist": reponse,
			},
			"code": reponse["state"],
			"msg":  reponse["msg"],
		},
	)
}

type MpBatcFileDeleteRequest struct {
	TaskId string `form:"task_id" binding:"required"`
}

// MpBatcFileDelete 删除小程序任务
func MpBatcFileDelete(c *gin.Context) {
	var req = MpBatcFileDeleteRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	str1 := strings.ReplaceAll(req.TaskId, "[", "")
	str2 := strings.ReplaceAll(str1, "]", "")
	idArray := strings.Split(str2, ",")
	var info model.CePingUserTask
	for _, taskId := range idArray {
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id =  ?", taskId).Delete(&info).Error; err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}
	response.OkWithData("删除成功", c)

}
