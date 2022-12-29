package ceping

import (
	"archive/zip"
	"bytes"
	"errors"
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
	"os"
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
		CheckAdInfo(taskId, han)
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

	// 1.解析内容
	var adResponse struct {
		utils.ResponseJson
		SecInfos string `json:"sec_infos"`
		ItemKnum int    `json:"item_knum"`
		TaskId   int    `json:"task_id"`
	}
	adClient := utils.NewClient(app.Conf.CePing.Ip, app.Conf.CePing.UserName, app.Conf.CePing.Password)
	adClient.Token = app.Conf.CePing.Token
	adClient.Signature = app.Conf.CePing.Signature
	err := adClient.BinCheckApkClient(&adResponse, templateInfo.Items, request.CallbackUrl, map[string]string{
		"apk": request.FilePath,
	})
	if err != nil {
		log.Error("调用测评平台接口解析内容失败", err)
		response.FailWithMessage("调用测评平台接口解析内容失败,err:"+err.Error(), c)
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

// CheckTask 检查正在测评的apk进度
func CheckAdInfo(taskId int, hand *Handler) {
	var responses struct {
		utils.ResponseJson
		FinishItem int    `json:"finish_item"`
		QueueNumer int    `json:"queue_numer"`
		Status     string `json:"status"`
		TotalItem  int    `json:"total_item"`
	}
	adClient := utils.NewClient(app.Conf.CePing.Ip, app.Conf.CePing.UserName, app.Conf.CePing.Password)
	adClient.Token = app.Conf.CePing.Token
	adClient.Signature = app.Conf.CePing.Signature
	if err := adClient.TaskProgress(&responses, "Android", taskId); err != nil {
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
		CheckAdInfo(taskId, hand)
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
		GetAdInfo(taskId)
		hand.RemoveTask(taskId)
	default:
		if err := app.DB.Model(&model.CePingUserTask{}).Where("task_id = ?", taskId).
			Updates(&taskInfo).Error; err != nil {
			log.Error(err.Error())
			return
		}
	}
}

// GetAdInfo 获取ad检测信息
func GetAdInfo(taskId int) {
	times := time.Now()
	var infoNum struct {
		utils.ResponseJson
		Data struct {
			RiskStatistic struct {
				Low    int `json:"low"`
				Medium int `json:"medium"`
				High   int `json:"high"`
			} `json:"risk_statistic"`
		} `json:"data"`
	}

	adClient := utils.NewClient(app.Conf.CePing.Ip, app.Conf.CePing.UserName, app.Conf.CePing.Password)
	adClient.Token = app.Conf.CePing.Token
	adClient.Signature = app.Conf.CePing.Signature
	if err := adClient.GetAdInfoClient(&infoNum, taskId); err != nil {
		log.Error(err)
	}

	if infoNum.Msg == "签名验证失败" || infoNum.Msg == "token验证失败" {
		app.Conf = app.LoadConfig()

		return
	}

	var info model.CePingUserTask
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

	var responses struct {
		utils.ResponseJson
		Status      string      `json:"status"`
		ResTxt      string      `json:"res_txt"`
		ItemKeys    string      `json:"item_keys"`
		ViewUrl     string      `json:"view_url"`
		AppPname    string      `json:"app_pname"`
		ItemNum     int         `json:"item_num"`
		DownUrl     string      `json:"down_url"`
		Id          int         `json:"id"`
		AppName     string      `json:"app_name"`
		UserId      int         `json:"user_id"`
		AppFilename string      `json:"app_filename"`
		CusBuilt    interface{} `json:"cus_built"`
		IsSample    int         `json:"isSample"`
		AppVersion  string      `json:"app_version"`
		ApkShield   string      `json:"apk_shield"`
		ApkMd5      string      `json:"apk_md5"`
		IsTarget    int         `json:"is_target"`
		PdfReport   string      `json:"pdf_report"`
		ItemKnum    int         `json:"item_knum"`
		TUpdate     string      `json:"t_update"`
		ResCode     int         `json:"res_code"`
		TCommit     string      `json:"t_commit"`
		SecInfos    string      `json:"sec_infos"`
		TaskId      int         `json:"task_id"`
		ApkUuid     string      `json:"apk_uuid"`
		Callback    string      `json:"callback"`
		AppSize     string      `json:"app_size"`
	}
	adClient := utils.NewClient(app.Conf.CePing.Ip, app.Conf.CePing.UserName, app.Conf.CePing.Password)
	adClient.Token = app.Conf.CePing.Token
	adClient.Signature = app.Conf.CePing.Signature
	if err := adClient.SearchOneDetailClient(&responses, req.TaskId); err != nil {
		log.Error(err)
	}
	response.OkWithData(responses, c)

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
