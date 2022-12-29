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
	"strconv"
	"strings"
	"sync"
	"time"
)

type IosHandle struct {
	Locker  sync.Mutex
	TaskIds []int
}

var IosHandler = NewIosHandler()

func NewIosHandler() *IosHandle {
	hand := &IosHandle{
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

func (h *IosHandle) Add(taskId int) {
	h.Locker.Lock()
	defer h.Locker.Unlock()
	h.TaskIds = append(h.TaskIds, taskId)
}
func (h *IosHandle) GetTaskIds() []int {
	h.Locker.Lock()
	defer h.Locker.Unlock()
	return h.TaskIds
}
func (h *IosHandle) RemoveTask(taskId int) {
	h.Locker.Lock()
	defer h.Locker.Unlock()
	for i, v := range h.TaskIds {
		if v == taskId {
			h.TaskIds = append(h.TaskIds[:i], h.TaskIds[i+1:]...)
			break
		}
	}
}
func (h *IosHandle) Check(hand *IosHandle) {
	if len(h.TaskIds) == 0 {
		return
	}
	for _, taskId := range h.GetTaskIds() {
		CheckIpaInfo(taskId, hand)
	}
}

type IosBinCheckRequest struct {
	CallbackUrl string `form:"callback_url"`
	TaskType    string `form:"task_type"`
	AppName     string `form:"app_name"`
	TemplateId  int    `form:"template_id"`
	FilePath    string `form:"file_path" binding:"required"`
}

// IosBinCheck 3.1.上传ipa并发送检测接口
func IosBinCheck(c *gin.Context) {

	var request = IosBinCheckRequest{}
	valid, errs := app.BindAndValid(c, &request)
	if !valid {
		log.Error("err:", errs.Error())
		response.FailWithMessage(errs.Error(), c)
		return
	}
	var templateInfo model.Template
	if request.TemplateId == 0 {
		templateInfo.Items = GetAllItems("ios")
	} else {
		templateInfo.ID = uint(request.TemplateId)
		if err := app.DB.Model(&model.Template{}).First(&templateInfo).Error; err != nil {
			log.Error("err:", err)
			response.FailWithMessage("未查询到当前模板", c)
			return
		}
	}
	// 1.解析内容
	var iosResponse struct {
		utils.ResponseJson
		SecInfos string `json:"sec_infos"`
		ItemKnum int    `json:"item_knum"`
		TaskId   int    `json:"task_id"`
	}
	adClient := utils.NewClient(app.Conf.CePing.Ip, app.Conf.CePing.UserName, app.Conf.CePing.Password)
	adClient.Token = app.Conf.CePing.Token
	adClient.Signature = app.Conf.CePing.Signature
	err := adClient.BinCheckApkClient(&iosResponse, templateInfo.Items, request.CallbackUrl, map[string]string{
		"ipa": request.FilePath,
	})

	if iosResponse.State != 200 {
		if iosResponse.Msg == "签名验证失败" || iosResponse.Msg == "token验证失败" {
			// 1.尝试是否可以获取到token
			_, _, err := app.GetCpToken(app.Conf.CePing.UserName, app.Conf.CePing.Password, app.Conf.CePing.Ip)
			if err != nil {
				// 如果获取不到就返回错误
				response.FailWithMessage("token获取失败，请检查配置", c)
				return
			}
			// 2.获取到token便重新调用该方法
			app.Conf = app.LoadConfig()
			IosBinCheck(c)
			return
		}
		log.Error("调用上传ios接口失败信息", iosResponse.Msg)
		response.FailWithMessage("调用测评平台上传ios接口失败,"+err.Error(), c)
		return
	}

	userId, _ := c.Get("userId")
	userID := userId.(uint)

	info := model.CePingUserTask{}
	info.TaskType = 2
	info.TaskID = strconv.Itoa(iosResponse.TaskId)
	info.AppName = request.AppName
	info.TemplateID = uint(request.TemplateId)
	info.Status = "测评中"
	info.UserID = userID
	info.FilePath = request.FilePath

	app.DB.Model(&model.CePingUserTask{}).Create(&info)

	IosHandler.Add(iosResponse.TaskId)

	response.OkWithData(iosResponse, c)
}

// CheckIpaInfo 获取当前正在检测ipa任务的信息
func CheckIpaInfo(taskId int, hand *IosHandle) {
	var responses struct {
		utils.ResponseJson
		FinishItem int    `json:"finish_item"`
		QueueNumer int    `json:"queue_numer"`
		Status     string `json:"status"`
		TotalItem  int    `json:"total_item"`
	}
	iosClient := utils.NewClient(app.Conf.CePing.Ip, app.Conf.CePing.UserName, app.Conf.CePing.Password)
	iosClient.Token = app.Conf.CePing.Token
	iosClient.Signature = app.Conf.CePing.Signature
	if err := iosClient.TaskProgress(&responses, "iOS", taskId); err != nil {
		log.Error(err)
	}

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
		CheckIpaInfo(taskId, hand)
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
		GetIpaInfo(taskId)
		hand.RemoveTask(taskId)
	default:
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).
			Updates(&taskInfo).Error; err != nil {
			log.Error(err.Error())
			return
		}
	}

}

// GetIpaInfo 获取ipa检测信息
func GetIpaInfo(taskId int) {
	times := time.Now()
	var infoNum struct {
		utils.ResponseJson
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
	}
	iosClient := utils.NewClient(app.Conf.CePing.Ip, app.Conf.CePing.UserName, app.Conf.CePing.Password)
	iosClient.Token = app.Conf.CePing.Token
	iosClient.Signature = app.Conf.CePing.Signature
	if err := iosClient.GetIpaInfoClient(&infoNum, taskId); err != nil {
		log.Error(err)
	}

	if infoNum.Msg == "签名验证失败" || infoNum.Msg == "token验证失败" {
		app.Conf = app.LoadConfig()

		return
	}
	var info model.CePingUserTask
	info.PkgName = infoNum.Data.AppInfo.PkgName
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

type IosSearchOneDetailRequest struct {
	TaskId int `form:"task_id" binding:"required"`
}

// IosSearchOneDetail 3.2.查询ipa检测任务的结果接口
func IosSearchOneDetail(c *gin.Context) {
	var req = IosSearchOneDetailRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	var responses struct {
		utils.ResponseJson
		AppMd5     string   `json:"app_md5"`
		AppName    string   `json:"app_name"`
		AppPname   string   `json:"app_pname"`
		AppSize    int      `json:"app_size"`
		AppVersion string   `json:"app_version"`
		DownUrl    string   `json:"down_url"`
		Id         int      `json:"id"`
		ItemKeys   []string `json:"item_keys"`
		ItemKnum   int      `json:"item_knum"`
		ItemNum    int      `json:"item_num"`
		ResCode    int      `json:"res_code"`
		Status     string   `json:"status"`
		TCommit    string   `json:"t_commit"`
		TUpdate    string   `json:"t_update"`
		UserId     int      `json:"user_id"`
	}
	iosClient := utils.NewClient(app.Conf.CePing.Ip, app.Conf.CePing.UserName, app.Conf.CePing.Password)
	iosClient.Token = app.Conf.CePing.Token
	iosClient.Signature = app.Conf.CePing.Signature
	if err := iosClient.IosSearchOneDetailClient(&responses, req.TaskId); err != nil {
		log.Error(err)
	}
	response.OkWithData(responses, c)
}

type IosBatchStatisticsResultRequest struct {
	TaskId int `form:"task_id" binding:"required"`
}

// IosIpaReport 3.4.下载测评ipa的word或pdf报告接口
func IosIpaReport(c *gin.Context) {
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

	clientURL := IP + "/v4/ios/ipa_report"

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
		err := Check(responseMap, reponse, IosIpaReport, c)
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
