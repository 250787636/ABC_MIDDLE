package ceping

import (
	"archive/zip"
	"bytes"
	"errors"
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

// 测评内置账号使用的数据
var USERNAME, TOKEN, PASSWORD, SIGNATURE, IP = CePingLoading()

// 请求频繁，共用连接
var Client = http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: false,
	},
}

// CePingLoading 获取内置账号信息
func CePingLoading() (string, string, string, string, string) {
	user := app.Conf.CePing
	return user.UserName, user.Token, user.Password, user.Signature, user.Ip
}

type Handler struct {
	Locker  sync.Mutex
	TaskIds []int
}

var AndroidHandler = NewHandler()

func NewHandler() *Handler {
	hand := &Handler{
		TaskIds: make([]int, 0),
	}
	go func() {
		fmt.Println("开始检查任务状态")
		for {
			hand.Check(hand)
			time.Sleep(20 * time.Second)
		}
	}()
	return hand
}

func (h *Handler) Add(taskId int) {
	h.Locker.Lock()
	defer h.Locker.Unlock()
	h.TaskIds = append(h.TaskIds, taskId)
}
func (h *Handler) GetTaskIds() []int {
	h.Locker.Lock()
	defer h.Locker.Unlock()
	return h.TaskIds
}

func (h *Handler) Check(han *Handler) {
	if len(h.TaskIds) == 0 {
		return
	}
	for _, taskId := range h.GetTaskIds() {
		fmt.Println("检查任务状态", taskId)
		CheckState(taskId, han)
	}
}

func (h *Handler) RemoveTask(taskId int) {
	h.Locker.Lock()
	defer h.Locker.Unlock()
	for i, v := range h.TaskIds {
		if v == taskId {
			h.TaskIds = append(h.TaskIds[:i], h.TaskIds[i+1:]...)
			break
		}
	}
}

func CheckState(taskId int, han *Handler) {
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
		log.Error("err:", err.Error())
		return
	}
	err = writer.Close()
	if err != nil {
		log.Error("err:", err.Error())
		return
	}

	clientURL := IP + "/v4/search_one_detail"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		log.Error("err:", err.Error())
		return
	}

	//注意别忘了设置header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(request)
	if err != nil {
		log.Error("err:", err.Error())
		return
	}
	post, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err.Error())
		return
	}

	reponse := make(map[string]interface{})
	err = jsoniter.Unmarshal(post, &reponse)
	if err != nil {
		log.Error("err:", err.Error())
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
		CheckState(taskId, han)
		return
	}
	if errMessage != "" {
		log.Error("app.conf", app.Conf.CePing.Signature)
		log.Error("err", errMessage)
		return
	}

	knum, _ := reponse["item_knum"].(float64)
	num, _ := reponse["item_num"].(float64)
	score, _ := reponse["app_score"].(float64)

	fmt.Println("预览地址", reponse["view_url"])
	viewUrl, _ := reponse["view_url"].(string)

	var TaskInfo model.CePingUserTask

	TaskInfo.Score = int(score)
	TaskInfo.ViewUrl = viewUrl
	TaskInfo.FinishItem = int(num)
	if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).Updates(&TaskInfo).Error; err != nil {
		if err != nil {
			log.Error("err:", err.Error())
			return
		}
	}

	var FindInfo model.CePingUserTask
	if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).First(&FindInfo).Error; err != nil {
		if err != nil {
			log.Error("err:", err.Error())
			return
		}
	}
	//fmt.Println("该任务的创建时间为",FindInfo.CreatedAt)
	if FindInfo.CreatedAt.Add(1*time.Hour).Unix() < time.Now().Unix() {
		han.RemoveTask(taskId)
		FindInfo.Status = "测评失败"
		FindInfo.Score = int(score)
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).Updates(&FindInfo).Error; err != nil {
			log.Error("修改任务状态失败")
			log.Error(err.Error())
			return
		}
		return
	}

	fmt.Println("score", score)

	if knum-num == 1 {
		var taskInfo model.CePingUserTask
		taskInfo.Status = "测评报告生成中"
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).
			Updates(&taskInfo).Error; err != nil {
			log.Error(err.Error())
			return
		}
	}

	var taskInfo model.CePingUserTask
	taskInfo.Score = int(score)
	taskInfo.Status = "测评完成"
	//如果已完成测评项
	if knum == num {
		//才会写所有的数据
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).
			Updates(&taskInfo).Error; err != nil {
			log.Error(err.Error())
			return
		}
		var GetTask model.CePingUserTask
		GetTask.TaskID = strconv.Itoa(taskId)
		GetTaskInfo(GetTask)

		//fmt.Println("处理完成的任务ID", taskId)
		//AndroidHandler.RemoveTask(taskId)
		han.RemoveTask(taskId)
		//fmt.Println("还剩下的任务ID", han.TaskIds)
	}

}

type BinCheckApkRequest struct {
	CallbackUrl string `form:"callback_url"`
	AppName     string `form:"app_name"` // 文件名称
	TemplateId  uint   `form:"template_id"`
	FilePath    string `form:"file_path"`
}

// BinCheckApk 2.1 上传apk并发送检测接口  android
func BinCheckApk(c *gin.Context) {
	var request = BinCheckApkRequest{}
	if err := c.Bind(&request); err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	file, err := os.Open(request.FilePath)
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage("打开文件失败", c)
		return
	}
	var templateInfo model.Template
	if request.TemplateId == 0 {
		templateInfo.Items = GetAllItems("ad")
	} else {
		templateInfo.ID = request.TemplateId
		if err := app.DB.Model(&model.Template{}).First(&templateInfo).Error; err != nil {
			log.Error("err:", err.Error())
			response.FailWithMessage("未查询到当前模板", c)
			return
		}
	}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	part, err := writer.CreateFormFile("apk", path.Base(file.Name()))
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	paramMap := make(map[string]interface{})
	paramMap["token"] = app.Conf.CePing.Token
	paramMap["signature"] = app.Conf.CePing.Signature
	paramMap["items"] = templateInfo.Items
	paramMap["callback_url"] = request.CallbackUrl

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

	clientURL := IP + "/v4/bin_check_apk"
	// 发送一个POST请求
	req, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage("构建请求失败", c)
		return
	}
	// 设置你需要的Header（不要想当然的手动设置Content-Type）multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// 执行请求
	resp, err := Client.Do(req)
	if err != nil {
		log.Error("执行请求err:", err.Error())
		response.FailWithMessage("调用测评对外接口失败,err:"+err.Error(), c)
		return
	}

	// 3.读取返回内容
	post, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("读取返回内容err:", err.Error())
		response.FailWithMessage("读取测评对外接口内容失败,err:"+err.Error(), c)
		return
	}
	// 1.解析内容
	var adResponse struct {
		Msg      string `json:"msg"`
		State    int    `json:"state"`
		SecInfos string `json:"sec_infos"`
		ItemKnum int    `json:"item_knum"`
		TaskId   int    `json:"task_id"`
	}
	err = jsoniter.Unmarshal(post, &adResponse)
	if err != nil {
		log.Error("调用测评平台接口解析内容失败", err)
		response.FailWithMessage("调用测评平台接口解析内容失败,err:"+err.Error(), c)
		return
	}

	if adResponse.State != 200 {
		if adResponse.Msg == "签名验证失败" || adResponse.Msg == "token验证失败" {
			// 1.尝试是否可以获取到token
			_, _, err := app.GetCpToken(app.Conf.CePing.UserName, app.Conf.CePing.Password, app.Conf.CePing.Ip)
			if err != nil {
				// 如果获取不到就返回错误
				response.FailWithMessage("token获取失败，请检查配置", c)
				return
			}
			fmt.Println("------")
			// 2.获取到token便重新调用该方法
			app.Conf = app.LoadConfig()
			BinCheckApk(c)
			return
		}
		// 如果不是token的原因便返回原有的错误
		log.Error("调用测评平台上传apk接口失败信息", string(post))
		response.FailWithMessage("调用测评平台上传apk接口失败信息", c)
		return
	}

	// 2.解析secInfos
	var secInfos struct {
		AppName string `json:"app_name"`
		ApkName string `json:"apk_name"`
		ApkVer  string `json:"apk_ver"`
	}

	err = jsoniter.Unmarshal([]byte(adResponse.SecInfos), &secInfos)
	if err != nil {
		log.Error("调用测评平台接口解析内容失败", err)
		response.FailWithMessage("调用测评平台接口解析内容失败,err:"+err.Error(), c)
		return
	}
	//fmt.Println(secInfos)

	userId, _ := c.Get("userId")
	userID := userId.(uint)
	//fmt.Println("userId", userID)

	// 创建任务
	user := model.CePingUserTask{
		TaskID:     strconv.Itoa(adResponse.TaskId),
		UserID:     userID,
		TemplateID: request.TemplateId,
		TaskType:   1,
		FilePath:   request.FilePath,
		AppName:    request.AppName,  // 用户传的文件名
		PkgName:    secInfos.AppName, // 测评平台获取的appName 对应应用名称
		Version:    secInfos.ApkVer,
		ItemsNum:   adResponse.ItemKnum}

	if err := app.DB.Create(&user).Error; err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	AndroidHandler.Add(adResponse.TaskId)

	response.OkWithData(gin.H{"datalist": adResponse}, c)

}

// GetTaskInfo 获取任务详细信息
func GetTaskInfo(task model.CePingUserTask) {

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	paramMap := make(map[string]interface{})
	paramMap["token"] = app.Conf.CePing.Token
	paramMap["signature"] = app.Conf.CePing.Signature
	paramMap["taskid"] = task.TaskID

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

	clientURL := IP + "/v4/batch_statistics_result"

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
		log.Error("err", errMessage)
		return
	}
	if errMessage != "" {
		log.Error("app.conf", app.Conf.CePing.Signature)
		log.Error("err", errMessage)
	}

	info := make(map[string]interface{})
	infostring, _ := jsoniter.Marshal(reponse["vulnerability_statistic"])
	err = jsoniter.Unmarshal(infostring, &info)
	if err != nil {
		return
	}
	//fmt.Println("info", info)
	//fmt.Println("HIGH111", info["middle"])

	//fmt.Println("已完成测评开始获取任务详细信息")

	hignNum := info["high"].(float64)
	middleNum := info["middle"].(float64)
	lowNum := info["low"].(float64)

	fmt.Println("middle", middleNum)
	//fmt.Println("----", reflect.TypeOf(info["middle"]))

	task.Status = "测评完成"
	task.HighNum = int(hignNum)
	task.MiddleNum = int(middleNum)
	task.LowNum = int(lowNum)
	task.RiskNum = task.HighNum + task.LowNum + task.MiddleNum
	times := time.Now()
	task.FinishedTime = &times

	if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id", task.TaskID).Updates(&task).Error; err != nil {
		log.Error(err)
		return
	}

	return
}

// CheckTask 检查正在测评的apk进度
func CheckTask(taskId int, modelInfo model.CePingUserTask) {
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

	clientURL := IP + "/v4/search_one_detail"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		return
	}

	//注意别忘了设置header
	request.Header.Set("Content-Type", writer.FormDataContentType())
	for {
		time.Sleep(10 * time.Second)
		resp, err := client.Do(request)
		if err != nil {
			log.Error(err)
			return
		}
		post, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
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
			app.Conf = app.LoadConfig()

			break
		}
		if errMessage != "" {
			fmt.Println("app.conf", app.Conf.CePing.Signature)
			log.Error("err", errMessage)
		}

		knum, _ := reponse["item_knum"].(float64)
		num, _ := reponse["item_num"].(float64)
		score, _ := reponse["app_score"].(float64)
		resCode, _ := reponse["res_code"].(float64)

		var TaskInfo model.CePingUserTask
		TaskInfo.Status = "测评失败"
		TaskInfo.Score = int(score)
		fmt.Println("resCode", resCode)

		//fmt.Println("score", score)

		var taskInfo model.CePingUserTask
		taskInfo.Score = int(score)
		taskInfo.Status = "测评完成"
		//如果已完成测评项
		if knum == num {
			//才会写所有的数据
			if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", modelInfo.TaskID).
				Updates(&taskInfo).Error; err != nil {
				log.Error("修改失败", err)
				return
			}
			GetTaskInfo(modelInfo)
			break
		}
	}

}

func checkError(post []byte) (rep map[string]interface{}, err error) {
	reponse := make(map[string]interface{})
	jsoniter.Unmarshal(post, &reponse)
	errMessage := ""
	if len(reponse) == 0 {
		errMessage = strings.Trim(string(post), `"`)
	} else if key, ok := reponse["state"].(float64); ok && key != 200 {
		errMessage = reponse["msg"].(string)
	}
	fmt.Println("err", errMessage)
	if errMessage != "" {
		app.Conf = app.LoadConfig()
		fmt.Println("app.conf", app.Conf.CePing.Signature)

		return nil, errors.New("正在重新加载配置文件,请重试")
	}
	return reponse, nil
}

type SearchOneRequest struct {
	TaskId int `form:"task_id" binding:"required"`
}

// SearchOneProgress 2.4查询某个正在测评的apk进度接口
func SearchOneProgress(c *gin.Context) {
	req := SearchOneRequest{}

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

	clientURL := IP + "/v4/search_one_progress"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
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
	//err = Check(reponse, post, BinCheckApk, c)
	//if err != nil {
	//	log.Error(err.Error())
	//
	//	return
	//}

	response.OkWithData(reponse, c)
}

type SearchOneDetailRequest struct {
	TaskId int `form:"task_id" binding:"required"`
}

//SearchOneDetail 2.5.查询某个测评apk的结果接口
func SearchOneDetail(c *gin.Context) {
	req := SearchOneDetailRequest{}
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

	clientURL := IP + "/v4/search_one_detail"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
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

	fmt.Println(string(post))
	reponse := make(map[string]interface{})
	err = jsoniter.Unmarshal(post, &reponse)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = Check(reponse, post, SearchOneDetail, c)
	if err != nil {
		log.Error(err.Error(), c)
		return
	}
	response.OkWithData(reponse, c)

}

// ApkReport 2.6.下载测评apk的word或pdf报告接口
func ApkReport(c *gin.Context) {

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

	clientURL := IP + "/v4/apk_report"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
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

	if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		reponse, _ := ioutil.ReadAll(resp.Body)
		responseMap := make(map[string]interface{})
		err := jsoniter.Unmarshal(reponse, &responseMap)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		err = Check(responseMap, reponse, ApkReport, c)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	// 控制用户请求所得的内容存为一个文件的时候提供一个默认的文件名
	c.Writer.Header().Set("Content-Disposition", contentDisposition)
	_, _ = io.Copy(c.Writer, resp.Body)

}

// BatchDownload 批量下载报告
func BatchDownload(c *gin.Context) {
	// 1.获取参数
	taskIdString := c.Query("task_id")
	downloadType := c.Query("download_type")
	fileType := c.Query("file_type")
	str1 := strings.ReplaceAll(taskIdString, "[", "")
	str2 := strings.ReplaceAll(str1, "]", "")
	idArray := strings.Split(str2, ",")
	fmt.Println("idArray", idArray)

	var TaskInfo []model.CePingUserTask
	if err := app.DB.Debug().Model(&model.CePingUserTask{}).Where("task_id in (?)", idArray).Find(&TaskInfo).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	out, err := os.Create("test.zip")
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Error("err:", err.Error())
		}
	}(out)

	writerZip := zip.NewWriter(out)

	clientURL := ""
	switch fileType {
	case "ad":
		clientURL = IP + "/v4/apk_report"
	case "ios":
		clientURL = IP + "/v4/ios/ipa_report"
	case "mp":
		clientURL = IP + "/v4/mp/mini_report"
	case "sdk":
		clientURL = IP + "/v4/sdk/sdk_report"
	default:
		clientURL = IP + "/v4/apk_report"
	}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	paramMap := make(map[string]interface{})
	for _, taskId := range idArray {
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

		//生成post请求
		client := &http.Client{}
		request, err := http.NewRequest("POST", clientURL, buff)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
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

		if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
			reponse, _ := ioutil.ReadAll(resp.Body)
			responseMap := make(map[string]interface{})
			err := jsoniter.Unmarshal(reponse, &responseMap)
			if err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}
			errMessage := ""
			if len(reponse) == 0 {
				errMessage = strings.Trim(string(reponse), `"`)
			} else if key, ok := responseMap["state"].(float64); ok && key != 200 {
				errMessage = responseMap["msg"].(string)
			}
			fmt.Println("err", errMessage)
			if errMessage == "签名验证失败" || errMessage == "token验证失败" {
				app.Conf = app.LoadConfig()
				response.FailWithMessage("token获取失败或者失效，请重试", c)
				return
			}
			if errMessage == "任务查询失败" {
				log.Error(err.Error())
				response.FailWithMessage(err.Error(), c)
				return
			}
			if errMessage != "" {
				log.Error(err.Error())
				response.FailWithMessage(err.Error(), c)
				return
			}
		}

		var taskInfo model.CePingUserTask
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).Find(&taskInfo).Error; err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		suffix := ".docx"
		if downloadType == "word" {
			suffix = ".docx"
		} else {
			suffix = ".pdf"
		}

		appName := taskInfo.AppName
		if taskInfo.TaskType == 4 {
			myTime := time.Now().Format("2006-01-02-15-04-05")
			appName = taskInfo.AppName + "-SDK安全测评报告(" + myTime + ")"
		}

		//appName := strings.Split(taskInfo.AppName, ".")
		fileWriter, err := writerZip.Create(appName + suffix)
		if err != nil {
			if os.IsPermission(err) {
				response.FailWithMessage(err.Error(), c)
				return
			}
			log.Error("Create file %s error: %s\n", taskInfo.AppName, err.Error())
			response.FailWithMessage(err.Error(), c)

			return
		}

		fileBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		_, err = fileWriter.Write(fileBody)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			log.Error("Write file error: ", err)
			return
		}
	}

	if err := writerZip.Close(); err != nil {
		response.FailWithMessage(err.Error(), c)
		log.Error("Close error: ", err)
		return
	}
	c.Header("Content-Type", "application/zip") // 这里是压缩文件类型 .zip
	c.Header("Content-Disposition", "inline;filename=测评报告下载.zip")
	c.File("test.zip")

}

type BatchStatisticsRequest struct {
	TaskID string `form:"taskid" binding:"required"`
}

// BatchStatisticsResult 2.7.查询测评apk统计及源结果接口
func BatchStatisticsResult(c *gin.Context) {
	var req = BatchStatisticsRequest{}
	if err := c.Bind(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	paramMap := make(map[string]interface{})
	paramMap["token"] = app.Conf.CePing.Token
	paramMap["signature"] = app.Conf.CePing.Signature
	paramMap["taskid"] = req.TaskID

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

	clientURL := IP + "/v4/batch_statistics_result"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
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
			// 如果获取不到就返回错误
			response.FailWithMessage("token获取失败，请检查配置", c)
			return
		}
		// 2.获取到token便重新调用该方法
		app.Conf = app.LoadConfig()
		BatchStatisticsResult(c)
		return
	}
	if errMessage != "" {
		fmt.Println("app.conf", app.Conf.CePing.Signature)
		log.Error(err.Error())
		response.FailWithMessage("调用测评平台接口失败"+err.Error(), c)
		return
	}

	marshal, err := jsoniter.Marshal(reponse["vulnerability_statistic"])
	if err != nil {
		return
	}
	var miNum map[string]interface{}
	err = jsoniter.Unmarshal(marshal, &miNum)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("333", miNum["middle"])
	midel, ok := miNum["middle"].(float64)
	if !ok {
		fmt.Println("不成功")
	}
	fmt.Println("成功", midel)

	response.OkWithData(reponse, c)

}

type BatchStatisticsOriginalRequest struct {
	TaskID string `form:"taskid" binding:"required"`
}

type GetItemsRequest struct {
	UserId int `form:"userid"`
}
type BatchFileDeleteRequest struct {
	TaskIds string `form:"task_id" binding:"required"`
}

// BatchFileDelete 2.12.批量删除apk物理文件接口
func BatchFileDelete(c *gin.Context) {
	var req = BatchFileDeleteRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	idArrayStr := strings.Split(req.TaskIds, ",")
	fmt.Println("------------", idArrayStr)
	for _, id := range idArrayStr {
		intId, _ := strconv.Atoi(id)
		if err := app.DB.Where("task_id = ?", id).Delete(&model.CePingUserTask{TaskID: strconv.Itoa(intId)}).Error; err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}
	res := gin.H{
		"code": 200,
		"info": gin.H{
			"msg":   "删除成功",
			"state": 200,
		},
		"msg": "请求成功",
	}
	response.OkWithData(res, c)
}

type GetAllInfoRequest struct {
	PkgName     string `form:"pkg_name"`
	CreatedName string `form:"user_name"`
	StartTime   string `form:"start_time"`
	EndTime     string `form:"end_time"`
	PageSize    int    `form:"size" binding:"required"`
	PageNumber  int    `form:"page" binding:"required"`
	TaskType    int    `form:"task_type" binding:"required"` // 1 android 2 ios 3 小程序  4 sdk
}

// GetAllInfo 获取 测评列表数据
func GetAllInfo(c *gin.Context) {
	req := GetAllInfoRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}
	var total int64

	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	isSuper, _ := c.Get("superAdmin")
	isSuperAdmin, ok := isSuper.(bool)
	if !ok {
		response.FailWithMessage("超级管理员标识错误", c)
		return
	}

	departId, _ := c.Get("departmentId")
	departmentId, ok := departId.(uint)
	if !ok {
		response.FailWithMessage("获取该员工部门ID失败", c)
		return
	}

	isAdm, _ := c.Get("isAdmin")
	isAdmin, ok := isAdm.(bool)
	if !ok {
		response.FailWithMessage("获取该员工是否为部门管理员失败", c)
		return
	}

	getUserId, _ := c.Get("userId")
	userId, ok := getUserId.(uint)
	if !ok {
		response.FailWithMessage("获取该员工ID失败", c)
		return
	}
	var responseAll []map[string]interface{}
	sql := app.DB.Debug().Model(&model.CePingUserTask{}).Where("ce_ping_user_task.task_type = ?", req.TaskType)

	if req.StartTime != "" {

		sql = sql.Where("ce_ping_user_task.created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {

		sql = sql.Where("ce_ping_user_task.created_at <= ?", req.EndTime)
	}
	if req.CreatedName != "" {
		sql = sql.Where("user.user_name like ?", "%"+req.CreatedName+"%")
	}
	if req.PkgName != "" {
		sql = sql.Where("ce_ping_user_task.pkg_name like ?", "%"+req.PkgName+"%")
	}
	// 1.如果是超级管理员可以查看所有的数据
	if isSuperAdmin {
		// 当为sdk的时候需要连sdkUse表
		if req.TaskType == 4 {
			sql = sql.Joins("inner join user on ce_ping_user_task.user_id = user.id").
				Joins("inner join sdk_use on sdk_use.task_id = ce_ping_user_task.task_id").
				Select("ce_ping_user_task.pkg_name," +
					"ce_ping_user_task.app_name," +
					"ce_ping_user_task.version," +
					"ce_ping_user_task.items_num," +
					"ce_ping_user_task.finish_item," +
					"ce_ping_user_task.score," +
					"user.user_name," +
					"ce_ping_user_task.created_at," +
					"ce_ping_user_task.finished_at," +
					"ce_ping_user_task.view_url," +
					"ce_ping_user_task.task_id," +
					"ce_ping_user_task.task_type," +
					"ce_ping_user_task.template_id," +
					"ce_ping_user_task.status," +
					"ce_ping_user_task.file_path," +
					"sdk_use.task_tag")
		} else {
			sql = sql.Joins("inner join user on ce_ping_user_task.user_id = user.id").
				Select("ce_ping_user_task.pkg_name," +
					"ce_ping_user_task.app_name," +
					"ce_ping_user_task.version," +
					"ce_ping_user_task.items_num," +
					"ce_ping_user_task.finish_item," +
					"ce_ping_user_task.score," +
					"user.user_name," +
					"ce_ping_user_task.created_at," +
					"ce_ping_user_task.finished_at," +
					"ce_ping_user_task.view_url," +
					"ce_ping_user_task.task_id," +
					"ce_ping_user_task.task_type," +
					"ce_ping_user_task.template_id," +
					"ce_ping_user_task.status," +
					"ce_ping_user_task.file_path")
		}
		sql.Count(&total).
			Offset((req.PageNumber - 1) * req.PageSize).
			Limit(req.PageSize).
			Order("ce_ping_user_task.created_at desc").
			Scan(&responseAll)

	} else if isAdmin {
		// 2.如果是部门管理员 就可以查看该部门下的所有数据
		sql.Where(" AND user.department_id = ?", departmentId)
		sql.Joins("inner join user on ce_ping_user_task.user_id = user.id").
			Select("ce_ping_user_task.pkg_name," +
				"ce_ping_user_task.app_name," +
				"ce_ping_user_task.version," +
				"ce_ping_user_task.items_num," +
				"ce_ping_user_task.finish_item," +
				"ce_ping_user_task.score," +
				"user.user_name," +
				"ce_ping_user_task.created_at," +
				"ce_ping_user_task.finished_at," +
				"ce_ping_user_task.view_url," +
				"ce_ping_user_task.task_id," +
				"ce_ping_user_task.task_type," +
				"ce_ping_user_task.template_id," +
				"ce_ping_user_task.status," +
				"ce_ping_user_task.file_path").
			Count(&total).
			Offset((req.PageNumber - 1) * req.PageSize).
			Limit(req.PageSize).
			Order("ce_ping_user_task.created_at desc").
			Scan(&responseAll)

	} else {
		// 3.如果是普通用户就获取自己的测评信息
		sql.Where(" AND ce_ping_user_task.user_id = ? ", userId)
		// 当为sdk的时候需要连sdkUse表
		if req.TaskType == 4 {
			sql = sql.Joins("inner join user on ce_ping_user_task.user_id = user.id").
				Joins("inner join sdk_use on sdk_use.task_id = ce_ping_user_task.task_id").
				Select("ce_ping_user_task.pkg_name," +
					"ce_ping_user_task.app_name," +
					"ce_ping_user_task.version," +
					"ce_ping_user_task.items_num," +
					"ce_ping_user_task.score," +
					"user.user_name," +
					"ce_ping_user_task.finish_item," +
					"ce_ping_user_task.created_at," +
					"ce_ping_user_task.finished_at," +
					"ce_ping_user_task.view_url," +
					"ce_ping_user_task.task_id," +
					"ce_ping_user_task.task_type," +
					"ce_ping_user_task.template_id," +
					"ce_ping_user_task.status," +
					"ce_ping_user_task.file_path" +
					"sdk_use.task_tag")
		} else {
			sql = sql.Joins("inner join user on ce_ping_user_task.user_id = user.id").
				Select("ce_ping_user_task.pkg_name," +
					"ce_ping_user_task.app_name," +
					"ce_ping_user_task.version," +
					"ce_ping_user_task.items_num," +
					"ce_ping_user_task.score," +
					"user.user_name," +
					"ce_ping_user_task.finish_item," +
					"ce_ping_user_task.created_at," +
					"ce_ping_user_task.finished_at," +
					"ce_ping_user_task.view_url," +
					"ce_ping_user_task.task_id," +
					"ce_ping_user_task.task_type," +
					"ce_ping_user_task.template_id," +
					"ce_ping_user_task.status," +
					"ce_ping_user_task.file_path")
		}
		sql.Count(&total).
			Offset((req.PageNumber - 1) * req.PageSize).
			Limit(req.PageSize).
			Order("ce_ping_user_task.created_at desc").
			Scan(&responseAll)
	}
	for _, v := range responseAll {
		v["created_at"] = v["created_at"].(time.Time).Format("2006-01-02 15:04:05")
		if v["finished_at"] != nil {
			v["finished_at"] = v["finished_at"].(time.Time).Format("2006-01-02 15:04:05")
		}
		if int(v["template_id"].(int64)) == 0 {
			switch int(v["task_type"].(int32)) {
			case 1:
				v["template_name"] = "Android-全量模板"
			case 2:
				v["template_name"] = "IOS-全量模板"
			case 3:
				v["template_name"] = "小程序-全量模板"
			case 4:
				v["template_name"] = "SDK-全量模板"
			}
		} else {
			var temp model.Template
			if err := app.DB.Model(model.Template{}).Select("template_name").Where("id = ?", v["template_id"]).Scan(&temp).Error; err != nil {
				log.Error(err.Error())
			}
			v["template_name"] = temp.TemplateName
		}

	}
	response.OkWithList(responseAll, int(total), req.PageNumber, req.PageSize, c)
	return
}

func TimeToGetToken() {
	time.Sleep(1 * time.Hour)
	app.LoadConfig()
}

//  types 取值 ad ios mp sdk
func GetAllItems(types string) string {
	toView, _ := TemplateItemKeysToView(app.DB, types, nil)
	itemkeyArr := make([]string, 0, 0)
	for _, name := range toView.CategorizedItems {
		for _, values := range name {
			if values.ItemKey != "" {
				itemkeyArr = append(itemkeyArr, values.ItemKey)
			}
		}
	}
	itemKeysArray := "["
	for key, value := range itemkeyArr {
		if key != len(itemkeyArr)-1 {
			itemKeysArray += "\"" + value + "\","

		} else {
			itemKeysArray += "\"" + value + "\""
		}
	}
	itemKeysArray += "]"

	return itemKeysArray
}
