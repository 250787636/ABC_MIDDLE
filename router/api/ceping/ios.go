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
	CallBackUrl string `form:"callback_url"`
	TaskType    string `form:"task_type"`
	AppName     string `form:"app_name"`
	TemplateId  int    `form:"template_id"`
	FilePath    string `form:"file_path" binding:"required"`
}

// IosBinCheck 3.1.上传ipa并发送检测接口
func IosBinCheck(c *gin.Context) {

	var FormReq = IosBinCheckRequest{}
	valid, errs := app.BindAndValid(c, &FormReq)
	if !valid {
		log.Error("err:", errs.Error())
		response.FailWithMessage(errs.Error(), c)
		return
	}
	var templateInfo model.Template
	if FormReq.TemplateId == 0 {
		templateInfo.Items = GetAllItems("ios")
	} else {
		templateInfo.ID = uint(FormReq.TemplateId)
		if err := app.DB.Model(&model.Template{}).First(&templateInfo).Error; err != nil {
			log.Error("err:", err)
			response.FailWithMessage("未查询到当前模板", c)
			return
		}
	}
	//fmt.Println("找出来的测评项", templateInfo.Items)
	// 1.1.读取文件
	open, err := os.Open(FormReq.FilePath)
	if err != nil {
		log.Error("err:", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	part, err := writer.CreateFormFile("ipa", path.Base(open.Name()))
	if err != nil {
		log.Error("err:", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	_, err = io.Copy(part, open)
	if err != nil {
		log.Error("err:", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	paramMap := make(map[string]interface{})
	paramMap["token"] = app.Conf.CePing.Token
	paramMap["signature"] = app.Conf.CePing.Signature
	paramMap["items"] = templateInfo.Items
	paramMap["callback_url"] = FormReq.CallBackUrl

	value, err := jsoniter.Marshal(paramMap)
	if err != nil {
		log.Error("err:", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = writer.WriteField("param", string(value))
	if err != nil {
		log.Error("err:", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = writer.Close()
	if err != nil {
		log.Error("err:", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	clientURL := IP + "/v4/ios/bin_check"
	// 发送一个POST请求
	req, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		log.Error("err:", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 设置你需要的Header（不要想当然的手动设置Content-Type）multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// 执行请求
	resp, err := Client.Do(req)
	if err != nil {
		log.Error("err:", err)
		response.FailWithMessage("调用测评平台上传ios接口失败", c)
		return
	}

	// 3.读取返回内容
	post, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		response.FailWithMessage("读取测评平台上传ios内容失败", c)
		return
	}
	//fmt.Println("psot", string(post))

	var iosResponse struct {
		Msg    string `json:"msg"`
		State  int    `json:"state"`
		TaskId int    `json:"task_id"`
	}

	err = jsoniter.Unmarshal(post, &iosResponse)
	if err != nil {
		log.Error("调用测评上传ios解析内容失败", err)
		response.FailWithMessage("调用测评平台接口解析内容失败,err:"+err.Error(), c)
		return
	}

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
	info.AppName = FormReq.AppName
	info.TemplateID = uint(FormReq.TemplateId)
	info.Status = "测评中"
	info.UserID = userID
	info.FilePath = FormReq.FilePath

	app.DB.Model(&model.CePingUserTask{}).Create(&info)

	IosHandler.Add(iosResponse.TaskId)

	response.OkWithData(iosResponse, c)
}

// CheckIpaInfo 获取当前正在检测ipa任务的信息
func CheckIpaInfo(taskId int, hand *IosHandle) {
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

	clientURL := IP + "/v4/ios/search_one_detail"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		return
	}

	//注意别忘了设置header
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
		log.Error(err)
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

	score := reponse["app_score"].(float64)

	fmt.Println("app", reponse["item_knum"])
	itemKnum := reponse["item_knum"].(float64)
	fmt.Println("down_url", reponse["view_url"])
	viewUrl, _ := reponse["view_url"].(string)
	num := reponse["item_num"].(float64)

	var getInfo model.CePingUserTask
	getInfo.ViewUrl = viewUrl
	getInfo.ItemsNum = int(itemKnum)
	getInfo.Score = int(score)
	getInfo.FinishItem = int(num)
	if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).Updates(&getInfo).Error; err != nil {
		if err != nil {
			log.Error(err)
			return
		}
	}

	var FindInfo model.CePingUserTask
	if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).First(&FindInfo).Error; err != nil {
		if err != nil {
			log.Error(err)

			return
		}
	}
	if FindInfo.CreatedAt.Add(1*time.Hour).Unix() < time.Now().Unix() {
		hand.RemoveTask(int(taskId))
		FindInfo.Status = "测评失败"

		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).Updates(&FindInfo).Error; err != nil {
			log.Error(err)

			return
		}
		return
	}

	pkgName := reponse["app_name"].(string)
	version := reponse["app_version"].(string)
	//fmt.Println("app", reponse)
	knum := reponse["item_knum"].(float64)
	if knum-num == 1 {
		var taskInfo model.CePingUserTask
		taskInfo.Status = "测评报告生成中"
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).
			Updates(&taskInfo).Error; err != nil {
			fmt.Println("修改失败", err)
			log.Error(err)

			return
		}
	}

	if knum == num {
		// 如果检测完毕 就获取检测信息
		GetIpaInfo(taskId, hand)
		//IosHandler.RemoveTask(int(taskId))
		hand.RemoveTask(int(taskId))
		return
	}

	var taskInfo model.CePingUserTask
	taskInfo.PkgName = pkgName
	taskInfo.Version = version
	if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).Updates(&taskInfo); err != nil {
		fmt.Println(err)
		log.Error(err)
		return
	}

}

// GetIpaInfo 获取ipa检测信息
func GetIpaInfo(taskId int, hand *IosHandle) {
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

	clientURL := IP + "/v4/ios/preview"

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
	//if infoNum.Msg != "" {
	//	fmt.Println("获取ios任务失败,错误为:", infoNum.Msg)
	//	return
	//}

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

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	paramMap := make(map[string]interface{})
	paramMap["token"] = app.Conf.CePing.Token
	paramMap["signature"] = app.Conf.CePing.Signature
	paramMap["taskid"] = req.TaskId

	value, err := jsoniter.Marshal(paramMap)
	if err != nil {
		response.FailWithMessage(errs.Error(), c)
		return
	}
	err = writer.WriteField("param", string(value))
	if err != nil {
		response.FailWithMessage(errs.Error(), c)
		return
	}
	err = writer.Close()
	if err != nil {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	clientURL := IP + "/v4/ios/search_one_detail"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	//注意别忘了设置header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	//Do方法发送请求
	resp, err := client.Do(request)
	if err != nil {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	post, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	reponse := make(map[string]interface{})
	err = jsoniter.Unmarshal(post, &reponse)
	if err != nil {
		log.Error(err)

		return
	}
	err = Check(reponse, post, IosSearchOneDetail, c)
	if err != nil {
		log.Error(err.Error())
		return
	}

	response.OkWithData(reponse, c)
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
