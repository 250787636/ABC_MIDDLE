package ceping

import (
	"errors"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/service"
	"example.com/m/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func Check(reponse map[string]interface{}, post []byte, handlerFunc gin.HandlerFunc, c *gin.Context) error {
	errMessage := ""
	if len(reponse) == 0 {
		errMessage = strings.Trim(string(post), `"`)
	} else if key, ok := reponse["state"].(float64); ok && key != 200 {
		errMessage = reponse["msg"].(string)
	}

	if errMessage == "签名验证失败" || errMessage == "token验证失败" {
		// 1.尝试是否可以获取到token
		_, _, err := app.GetCpToken(app.Conf.CePing.UserName, app.Conf.CePing.Password, app.Conf.CePing.Ip)
		if err != nil {
			// 如果获取不到就返回错误
			response.FailWithMessage("token获取失败，请检查配置", c)
			return errors.New("token获取失败，请检查配置")
		}
		// 2.获取到token便重新调用该方法
		app.Conf = app.LoadConfig()
		handlerFunc(c)
		return errors.New(errMessage)
	}
	if errMessage != "" {
		log.Error("err", errMessage)
		response.FailWithMessage("调用测评接口失败，错误信息:"+errMessage, c)
		return errors.New(errMessage)
	}
	return nil
}

type AddTemplateRequest struct {
	TemplateType   string `form:"template_type" binding:"required"` //模板类型 安卓ad 苹果 ios 小程序 mp
	TemplateName   string `form:"template_name" binding:"required"`
	ItemKeys       string `form:"item_keys" binding:"required"`
	IsOWASP        bool   `form:"is_owasp"`        // 是否是OWASP模板 true 是
	ReportLanguage string `form:"report_language"` //导出模板语言
}

// AddTemplate 新增模版
func AddTemplate(c *gin.Context) {

	req := AddTemplateRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}
	var tem model.Template
	if err := app.DB.Model(&model.Template{}).Where("template_name = ?", req.TemplateName).Where("template_type = ?", req.TemplateType).Find(&tem).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if tem.ID > 0 {
		response.FailWithMessage("该模板名称已存在", c)
		return
	}

	str1 := strings.ReplaceAll(req.ItemKeys, "[", "")
	str2 := strings.ReplaceAll(str1, "]", "")
	idArray := strings.Split(str2, ",")

	itemKeysArray := "["

	for key, value := range idArray {
		if key != len(idArray)-1 {
			itemKeysArray += "\"" + value + "\","

		} else {
			itemKeysArray += "\"" + value + "\""
		}
	}
	itemKeysArray += "]"
	fmt.Println("itemKeysArray", itemKeysArray)

	id, exist := c.Get("userId")
	if !exist {
		response.FailWithMessage("未获取到userid", c)
		return
	}
	userId, ok := id.(uint)
	if !ok {
		response.FailWithMessage("未获取到userid", c)
		return
	}
	//fmt.Println("userid", userId)
	var info model.Template
	info.CreatedID = int(userId)
	info.TemplateName = req.TemplateName
	info.TemplateType = req.TemplateType
	if req.IsOWASP == true {
		info.IsOwasp = 1
	} else {
		info.IsOwasp = 2
	}
	info.ReportLanguage = req.ReportLanguage

	info.Items = itemKeysArray

	if err := app.DB.Model(&model.Template{}).Create(&info).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("创建成功", c)
}

type GetTemplateRequest struct {
	PageSize     int    `form:"size"`
	PageNumber   int    `form:"page"`
	TemplateType string `form:"template_type" binding:"required"` //  安卓ad 苹果 ios 小程序 mp  SDK任务 sdk
	IsPage       int    `form:"is_page" binding:"required"`       // 是否需要分页 1 是  2 不是
}

// GetTemplate 获取模版列表
func GetTemplate(c *gin.Context) {
	req := GetTemplateRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	isSuper, _ := c.Get("superAdmin")
	isSuperAdmin, ok := isSuper.(bool)
	if !ok {
		log.Error("获取超级管理员标识错误")
		response.FailWithMessage("超级管理员标识错误", c)
		return
	}

	departId, _ := c.Get("departmentId")
	departmentId, ok := departId.(uint)
	if !ok {
		log.Error("获取该员工部门ID失败")
		response.FailWithMessage("获取该员工部门ID失败", c)
		return
	}

	isAdm, _ := c.Get("isAdmin")
	isAdmin, ok := isAdm.(bool)
	if !ok {
		log.Error("获取该员工是否为部门管理员失败")
		response.FailWithMessage("获取该员工是否为部门管理员失败", c)
		return
	}

	getUserId, _ := c.Get("userId")
	userId, ok := getUserId.(uint)
	if !ok {
		response.FailWithMessage("获取该员工ID失败", c)
		return
	}
	total, list, err := service.Template.ListTemplate(
		service.QueryTemplateRequest{
			PageSize:     req.PageSize,
			PageNumber:   req.PageNumber,
			TemplateType: req.TemplateType,
			IsPage:       req.IsPage == 1,
			IsSuper:      isSuperAdmin,
			IsAdmin:      isAdmin,
			DepartId:     int(departmentId),
			UserId:       int(userId),
		})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	response.OkWithList(list, int(total), req.PageNumber, req.PageSize, c)
}

type FixTemplateRequest struct {
	TemplateType   string `form:"template_type" binding:"required"` //  模板类型 1 android 2 ios 3 小程序
	TemplateId     int    `form:"template_id" binding:"required"`
	TemplateName   string `form:"template_name"`
	ItemKeys       string `form:"item_keys"` //测评项
	IsOWASP        bool   `form:"is_owasp"`
	ReportLanguage string `form:"report_language"` //导出模板语言
}

// FixTemplate 修改模版
func FixTemplate(c *gin.Context) {
	req := FixTemplateRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	str1 := strings.ReplaceAll(req.ItemKeys, "[", "")
	str2 := strings.ReplaceAll(str1, "]", "")
	idArray := strings.Split(str2, ",")

	itemKeysArray := "["

	for key, value := range idArray {
		if key != len(idArray)-1 {
			itemKeysArray += "\"" + value + "\","

		} else {
			itemKeysArray += "\"" + value + "\""
		}
	}
	itemKeysArray += "]"
	var info model.Template
	info.ID = uint(req.TemplateId)
	info.TemplateName = req.TemplateName
	info.ReportLanguage = req.ReportLanguage
	info.Items = itemKeysArray

	//fmt.Println("info", info)
	if req.IsOWASP == true {
		info.IsOwasp = 1
	} else {
		info.IsOwasp = 2
	}

	if err := app.DB.Model(&model.Template{}).Where("id = ?", req.TemplateId).Updates(&info).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("修改成功", c)
}

type DelTemplateRequest struct {
	TemplateId int `form:"template_id" binding:"required"`
}

// DeleteTemplate 删除模版
func DeleteTemplate(c *gin.Context) {
	req := DelTemplateRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	var info model.Template
	if err := app.DB.Model(&model.Template{}).Where("id = ?", req.TemplateId).Delete(&info).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)

}

type GetTemplateItemsRequest struct {
	TemplateType string `form:"template_type" binding:"required"`
	TemplateId   int    `form:"template_id" ` // 0 全部的
}

// GetTemplateItems 获取模版详细测评项
func GetTemplateItems(c *gin.Context) {
	req := GetTemplateItemsRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	type TemplateItem struct {
		ItemKey   string `json:"item_key"`
		AuditName string `json:"audit_name"`
		IsDynamic bool   `json:"is_dynamic"`
		Status    int    `json:"status"`
	}
	if req.TemplateId == 0 {
		toView, err := TemplateItemKeysToView(app.DB, req.TemplateType, nil)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		var newCataries []string
		for _, value := range toView.Categories {
			if value == "组件安全" && req.TemplateType == "ad" {
				continue
			}
			newCataries = append(newCataries, value)
		}
		tempName := ""
		switch req.TemplateType {
		case "ad":
			tempName = "Android全量检测"
		case "ios":
			tempName = "iOS全量检测"
		case "sdk":
			tempName = "SDK全量检测"
		}
		toView.CreatorAccount = "0"
		toView.CreateTime = "2022-10-01 00:00:00"
		toView.ReportLanguage = "zh_cn"
		toView.TemplateName = tempName
		toView.Categories = newCataries
		response.OkWithData(toView, c)
		return
	}

	var record model.Template
	record.TemplateType = req.TemplateType
	if err := app.DB.Model(&model.Template{}).Where("id = ?", req.TemplateId).Where("template_type = ?", req.TemplateType).First(&record).Error; err != nil {
		response.FailWithMessage("未找到该模板", c)
		return
	}
	str1 := strings.ReplaceAll(record.Items, "[", "")
	str2 := strings.ReplaceAll(str1, "]", "")
	idArray := strings.Split(str2, ",")
	keysToView, err := TemplateItemKeysToView(app.DB, req.TemplateType, idArray)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	keysToView.TemplateId = int(record.ID)
	keysToView.TemplateName = record.TemplateName
	keysToView.CreatorAccount = strconv.Itoa(record.CreatedID)
	keysToView.CreateTime = record.CreatedAt.Format("2006-01-02 15:04:05")
	if record.IsOwasp == 1 {
		keysToView.IsOWASP = true
	} else {
		keysToView.IsOWASP = false
	}
	keysToView.ReportLanguage = record.ReportLanguage
	//var itemKeys []string
	itemKeySet := make(map[string]int, 0)
	item1 := strings.ReplaceAll(record.Items, "[", "")
	item2 := strings.ReplaceAll(item1, "]", "")
	item3 := strings.Split(item2, ",")

	for _, value := range item3 {
		str3 := strings.ReplaceAll(value, "\"", "")
		itemKeySet[str3] = 1
	}
	for _, items := range keysToView.CategorizedItems {
		for i := range items {
			itemKey := items[i].ItemKey
			if _, ok := itemKeySet[itemKey]; ok {
				items[i].Status = 1
			}
		}
	}
	var newCataries []string
	for _, value := range keysToView.Categories {
		if value == "组件安全" && req.TemplateType == "ad" {
			continue
		}
		newCataries = append(newCataries, value)
	}
	keysToView.Categories = newCataries
	response.OkWithData(keysToView, c)
	return
}
