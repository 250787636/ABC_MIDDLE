package jiagu

import (
	"bytes"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/utils"
	"example.com/m/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type JiaGuData struct {
	ID           uint      `json:"id"`
	AppName      string    `json:"app_name"`
	AppType      string    `json:"app_type"`
	AppVersion   string    `json:"app_version"`
	CreatedAt    time.Time `json:"created_at"`
	FinishTime   time.Time `json:"finish_time"`
	PolicyId     int       `json:"policy_id"`
	PolicyReason string    `json:"policy_reason"`
	TaskId       int       `json:"task_id"`
	TaskStatus   string    `json:"task_status"`
	UserName     string    `json:"user_name"`
}

// 加固记录导出功能
func Exporting(c *gin.Context) {
	var list []JiaGuData
	var fileName string
	num, typeIdOK := c.GetQuery("type_id")
	typeIdOK = utils.IsStringEmpty(num, typeIdOK)
	if !typeIdOK {
		response.FailResult(http.StatusInternalServerError, "", model.ReqParameterMissing, c)
		return
	}
	typeId, _ := strconv.Atoi(num)
	// 文件名赋值
	if typeId == 1 {
		fileName = "android加固记录.xlsx"
	} else if typeId == 2 {
		fileName = "h5加固记录.xlsx"
	}

	departmentId, _ := c.Get("departmentId")
	var sqlString bytes.Buffer
	sqlString.WriteString("application_type.id = ")
	sqlString.WriteString(strconv.Itoa(typeId))

	superAdmin, _ := c.Get("superAdmin")
	if !superAdmin.(bool) {
		sqlString.WriteString(" AND user.department_id = ")
		sqlString.WriteString(strconv.Itoa(int(departmentId.(uint))))
	}
	//"application_type.id = ? AND user.department_id = ?", typeId, departmentId
	if err := app.DB.Table("jia_gu_task").
		Select(" jia_gu_task.id,user.user_name,application.app_name,application.app_version,jia_gu_task.task_id,jia_gu_task.policy_id,jia_gu_task.policy_reason,jia_gu_task.created_at,jia_gu_task.task_status,jia_gu_task.finish_time,application_type.app_type").
		Joins("INNER JOIN user ON user.id = jia_gu_task.user_id").
		Joins("INNER JOIN application_type ON application_type.id = jia_gu_task.app_type_id").
		Joins("INNER JOIN application ON application.id = jia_gu_task.app_id").
		Where(sqlString.String()).
		Scan(&list).
		Error; err != nil {
		response.FailResult(http.StatusInternalServerError, "", err.Error(), c)
		return
	}
	//导出
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	row.AddCell().Value = "记录ID"
	row.AddCell().Value = "用户名"
	row.AddCell().Value = "系统"
	row.AddCell().Value = "应用类型"
	row.AddCell().Value = "系统英文简称"
	row.AddCell().Value = "策略id"
	row.AddCell().Value = "不使用推荐策略理由"
	row.AddCell().Value = "任务id"
	row.AddCell().Value = "任务状态"
	row.AddCell().Value = "提交时间"
	row.AddCell().Value = "完成时间"
	// 前端列表为序列号 导出记录需为序列号 关联onces BLk3MRv8GAlE7yKR
	theNum := 1
	layout := "2006-01-02 15:04:05"
	for _, v := range list {
		row = sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(theNum)
		row.AddCell().Value = v.UserName
		row.AddCell().Value = v.AppName
		row.AddCell().Value = v.AppType
		row.AddCell().Value = v.AppVersion
		row.AddCell().Value = strconv.Itoa(v.PolicyId)
		row.AddCell().Value = v.PolicyReason
		row.AddCell().Value = strconv.Itoa(v.TaskId)
		row.AddCell().Value = v.TaskStatus
		row.AddCell().Value = v.CreatedAt.Format(layout)
		row.AddCell().Value = v.FinishTime.Format(layout)
		theNum++
	}
	buf := new(bytes.Buffer)
	err = file.Write(buf)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, url.QueryEscape(fileName)))
	c.Data(http.StatusOK, "text/xlsx", buf.Bytes())
}
