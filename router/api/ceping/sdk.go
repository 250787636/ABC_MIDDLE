package ceping

import (
	"bytes"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/utils"
	"example.com/m/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SdkHandle struct {
	Locker  sync.Mutex
	TaskIds []int
}

var SdkHandler = NewSdkHandler()

func NewSdkHandler() *SdkHandle {
	hand := &SdkHandle{
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

func (h *SdkHandle) Add(taskId int) {
	h.Locker.Lock()
	defer h.Locker.Unlock()
	h.TaskIds = append(h.TaskIds, taskId)
}
func (h *SdkHandle) GetTaskIds() []int {
	h.Locker.Lock()
	defer h.Locker.Unlock()
	return h.TaskIds
}
func (h *SdkHandle) RemoveTask(taskId int) {
	h.Locker.Lock()
	defer h.Locker.Unlock()
	for i, v := range h.TaskIds {
		if v == taskId {
			h.TaskIds = append(h.TaskIds[:i], h.TaskIds[i+1:]...)
			break
		}
	}
}
func (h *SdkHandle) Check(hand *SdkHandle) {
	if len(h.TaskIds) == 0 {
		return
	}
	for _, taskId := range h.GetTaskIds() {
		CheckSdkInfo(taskId, hand)
	}
}

type SdkBinCheckRequest struct {
	CallBackUrl string `form:"callback_url"`
	TaskType    string `form:"task_type"`
	TemplateId  int    `form:"template_id"`
	FilePath    string `form:"file_path" binding:"required"`
	SdkName     string `form:"sdk_name" binding:"required"`
	SdkVersion  string `form:"sdk_version" binding:"required"`
	SdkTag      string `form:"sdk_tag" binding:"required"`
}

// SdkBinCheck 14.1.上传sdk并发送检测接口
func SdkBinCheck(c *gin.Context) {

	var request = SdkBinCheckRequest{}
	valid, errs := app.BindAndValid(c, &request)
	if !valid {
		log.Error("err:", errs.Error())
		response.FailWithMessage(errs.Error(), c)
		return
	}
	var templateInfo model.Template
	if request.TemplateId == 0 {
		templateInfo.Items = GetAllItems("sdk")
	} else {
		templateInfo.ID = uint(request.TemplateId)
		if err := app.DB.Model(&model.Template{}).First(&templateInfo).Error; err != nil {
			log.Error("err:", err)
			response.FailWithMessage("未查询到当前模板", c)
			return
		}
	}
	var sdkResponse struct {
		utils.ResponseJson
		SecInfos string `json:"sec_infos"`
		ItemKnum int    `json:"item_knum"`
		TaskId   int    `json:"task_id"`
	}
	sdkClient := utils.NewClient(app.Conf.CePing.Ip, app.Conf.CePing.UserName, app.Conf.CePing.Password)
	sdkClient.Token = app.Conf.CePing.Token
	sdkClient.Signature = app.Conf.CePing.Signature
	err := sdkClient.BinCheckApkClient(&sdkResponse, templateInfo.Items, request.CallBackUrl, map[string]string{
		"sdk": request.FilePath,
	})

	if sdkResponse.State != 200 {
		if sdkResponse.Msg == "签名验证失败" || sdkResponse.Msg == "token验证失败" {
			log.Error(sdkResponse.Msg)
			// 1.尝试是否可以获取到token
			_, _, err := app.GetCpToken(app.Conf.CePing.UserName, app.Conf.CePing.Password, app.Conf.CePing.Ip)
			if err != nil {
				// 如果获取不到就返回错误
				response.FailWithMessage("token获取失败，请检查配置", c)
				return
			}
			// 2.获取到token便重新调用该方法
			app.Conf = app.LoadConfig()
			SdkBinCheck(c)
			return
		}
		log.Error("调用上传sdk接口失败信息", sdkResponse.Msg)
		response.FailWithMessage("调用测评平台上传sdk接口失败"+sdkResponse.Msg, c)
		return
	}
	//sdkclient := utils.NewClient(app.Conf.CePing.Ip, app.Conf.CePing.UserName, app.Conf.CePing.Password)
	//sdkclient.Token = app.Conf.CePing.Token
	//sdkclient.Signature = app.Conf.CePing.Signature
	var sdkUrl struct {
		utils.ResponseJson
		SdKUrl   string `json:"sdk_url"`
		ItemKnum int    `json:"item_knum"`
	}

	if err = sdkClient.SdkSearchOneDetail(&sdkUrl, sdkResponse.TaskId); err != nil {
		log.Error(err)
	}

	userId, _ := c.Get("userId")
	userID := userId.(uint)

	info := model.CePingUserTask{}
	info.TaskType = 4
	info.TaskID = strconv.Itoa(sdkResponse.TaskId)
	info.PkgName = path.Base(request.FilePath)
	info.AppName = request.SdkName
	info.TemplateID = uint(request.TemplateId)
	info.Version = request.SdkVersion
	info.Status = "测评中"
	info.UserID = userID
	info.FilePath = request.FilePath
	info.ViewUrl = sdkUrl.SdKUrl
	info.ItemsNum = sdkUrl.ItemKnum

	app.DB.Model(&model.CePingUserTask{}).Create(&info)

	sdkUse := model.SdkUse{}
	sdkUse.TaskID = info.TaskID
	sdkUse.TaskTag = request.SdkTag
	app.DB.Model(&model.SdkUse{}).Create(&sdkUse)

	SdkHandler.Add(sdkResponse.TaskId)

	response.OkWithData(sdkResponse, c)
}

// CheckSdkInfo 获取当前正在检测sdk任务的信息
func CheckSdkInfo(taskId int, hand *SdkHandle) {
	var responses struct {
		utils.ResponseJson
		FinishItem int    `json:"finish_item"`
		QueueNumer int    `json:"queue_numer"`
		Status     string `json:"status"`
		TotalItem  int    `json:"total_item"`
	}

	sdkclient := utils.NewClient(app.Conf.CePing.Ip, app.Conf.CePing.UserName, app.Conf.CePing.Password)
	sdkclient.Token = app.Conf.CePing.Token
	sdkclient.Signature = app.Conf.CePing.Signature

	if err := sdkclient.TaskProgress(&responses, "SDK", taskId); err != nil {
		log.Error(err)
	}
	//if err := responses.Unmarshal(&responses); err != nil {
	//	log.Error(err)
	//}

	errMessage := ""
	if responses.State != 200 {
		errMessage = responses.Msg
	}
	fmt.Println("err", errMessage)
	if errMessage == "签名验证失败" || errMessage == "token验证失败" {
		// 1.尝试是否可以获取到token
		_, _, err := app.GetCpToken(app.Conf.CePing.UserName, app.Conf.CePing.Password, app.Conf.CePing.Ip)
		if err != nil {
			log.Error("err:", err.Error())
			return
		}
		// 2.获取到token便重新调用该方法
		app.Conf = app.LoadConfig()
		CheckSdkInfo(taskId, hand)
		return
	}
	if errMessage != "" {
		log.Error("err", errMessage)
		return
	}

	var taskInfo model.CePingUserTask
	taskInfo.Status = "测评中"
	taskInfo.ItemsNum = responses.TotalItem
	taskInfo.FinishItem = responses.FinishItem

	switch responses.Status {
	case "EXCEPTION":
		hand.RemoveTask(taskId)
		taskInfo.Status = "测评失败"
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).Updates(&taskInfo).Error; err != nil {
			log.Error(err)
		}
		return
	case "REPORT_GENERATING":
		taskInfo.Status = "测评报告生成中"
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).
			Updates(&taskInfo).Error; err != nil {
			fmt.Println("修改失败", err)
			log.Error(err)
			return
		}
	case "FINISHED":
		taskInfo.Status = "测评完成"
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).
			Updates(&taskInfo).Error; err != nil {
			log.Error(err.Error())
			return
		}
		// 如果检测完毕 就获取检测信息
		GetSdkInfo(taskId, hand)
		hand.RemoveTask(taskId)
	default:
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).
			Updates(&taskInfo).Error; err != nil {
			log.Error(err.Error())
			return
		}
	}
}

// GetSdkInfo 获取sdk检测信息
func GetSdkInfo(taskId int, hand *SdkHandle) {
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

	clientURL := IP + "/v4/sdk/preview"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		log.Error(err)

		return
	}

	//注意别忘了设置header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	//Do方法发送请求
	resp, err := client.Do(request)
	defer resp.Body.Close()
	post, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	times := time.Now()

	var infoNum struct {
		Data struct {
			AppInfo struct {
				PkgName string `json:"app_name"`
			} `json:"app_info"`
			RiskStatistic struct {
				Low    int `json:"low"`
				Medium int `json:"medium"`
				High   int `json:"high"`
			} `json:"risk_statistic"`
		} `json:"data"`
		Msg   string `json:"msg"`
		State int    `json:"state"`
	}

	err = jsoniter.Unmarshal(post, &infoNum)
	if err != nil {
		fmt.Println("解析失败")
		log.Error(err)

		return
	}
	if infoNum.Msg == "签名验证失败" || infoNum.Msg == "token验证失败" {
		app.Conf = app.LoadConfig()

		return
	}

	var info model.CePingUserTask
	//info.PkgName = infoNum.Data.AppInfo.PkgName
	info.FinishedTime = &times
	info.LowNum = infoNum.Data.RiskStatistic.Low
	info.MiddleNum = infoNum.Data.RiskStatistic.Medium
	info.HighNum = infoNum.Data.RiskStatistic.High
	info.RiskNum = infoNum.Data.RiskStatistic.Low + infoNum.Data.RiskStatistic.Medium + infoNum.Data.RiskStatistic.High
	info.Status = "测评完成"
	if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).Updates(&info).Error; err != nil {
		log.Error(err)

		fmt.Println(err)
	}
}

type SdkSearchOneDetailRequest struct {
	TaskId int `form:"task_id" binding:"required"`
}

// SdkSearchOneDetail 14.2.查询sdk检测任务的结果接口
func SdkSearchOneDetail(c *gin.Context) {
	var req = SdkSearchOneDetailRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	var responses struct {
		utils.ResponseJson
		DownUrl    string   `json:"down_url"`
		Id         int      `json:"id"`
		IsDeleted  bool     `json:"is_deleted"`
		ItemKeys   []string `json:"item_keys"`
		ItemKnum   int      `json:"item_knum"`
		ItemNum    int      `json:"item_num"`
		ResCode    int      `json:"res_code"`
		SdkMd5     string   `json:"sdk_md5"`
		SdkName    string   `json:"sdk_name"`
		SdkSize    int      `json:"sdk_size"`
		SdkVersion string   `json:"sdk_version"`
		TCommit    string   `json:"t_commit"`
		TUpdate    string   `json:"t_update"`
		UserId     int      `json:"user_id"`
	}

	sdkClient := utils.NewClient(app.Conf.CePing.Ip, app.Conf.CePing.UserName, app.Conf.CePing.Password)
	sdkClient.Token = app.Conf.CePing.Token
	sdkClient.Signature = app.Conf.CePing.Signature
	if err := sdkClient.SearchOneDetailClient(&responses, req.TaskId); err != nil {
		log.Error(err)
	}
	response.OkWithData(responses, c)
}

type SdkBatchStatisticsResultRequest struct {
	TaskId int `form:"task_id" binding:"required"`
}

// SdkReport 14.4.下载测评sdk的word或pdf报告接口
func SdkReport(c *gin.Context) {
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

	clientURL := IP + "/v4/sdk/sdk_report"

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
		err := Check(responseMap, reponse, SdkReport, c)
		if err != nil {
			log.Error("err:", err.Error())
			return
		}
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	// 控制用户请求所得的内容存为一个文件的时候提供一个默认的文件名
	c.Writer.Header().Set("Content-Disposition", contentDisposition)
	_, _ = io.Copy(c.Writer, resp.Body)

}
