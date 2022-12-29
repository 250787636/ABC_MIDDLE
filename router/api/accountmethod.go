package api

import (
	"bytes"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/utils"
	"example.com/m/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 查询用户信息
func GetListAccount(c *gin.Context) {
	var accountList = new([]map[string]interface{})
	// 判断非必传值是否存在
	num, department_ok := c.GetPostForm("department_id")
	username, username_ok := c.GetPostForm("username")

	page, size, offNum := utils.GetFromDataPageSizeOffNum(c)
	if page+size+offNum == -3 {
		return
	}

	department_ok = utils.IsStringEmpty(num, department_ok)
	username_ok = utils.IsStringEmpty(username, username_ok)

	departmentID, _ := strconv.Atoi(num)

	var total int64
	var sqlString bytes.Buffer
	// 拼接sql语句
	sqlString.WriteString("user.deleted_at is NULL")

	// 获取当前账户信息
	myDepartmentId, _ := c.Get("departmentId")
	isSuperAdmin, _ := c.Get("superAdmin")

	// 不为超级管理员则会默认加入当前用户的部门编号进行数据隔离
	if isSuperAdmin.(bool) {
		if department_ok {
			sqlString.WriteString(" AND user.department_id = ")
			sqlString.WriteString(strconv.Itoa(departmentID))
		}
	} else {
		sqlString.WriteString(" AND user.department_id = ")
		sqlString.WriteString(strconv.Itoa(int(myDepartmentId.(uint))))
	}

	if username_ok {
		sqlString.WriteString(" AND user.user_name like '")
		sqlString.WriteString("%" + username + "%")
		sqlString.WriteString("'")
	}

	if err := app.DB.Table("user").
		Select("user.id,user.user_name,departments.department_name,user.account_level,user.job_title,user.last_login_time").
		Joins("LEFT JOIN departments  ON user.department_id = departments.id").
		Where(sqlString.String()).
		Count(&total).
		Offset(offNum).
		Limit(size).
		Order("user.created_at desc").
		Scan(accountList).Error; err != nil {
		response.FailResult(http.StatusInternalServerError, "", err.Error(), c)
		return
	}

	//for _, val := range *accountList {
	//	arr := strings.Split(val["user_name"].(string), "-")
	//	if len(arr) > 1 {
	//		val["user_name"] = arr[1]
	//	} else {
	//		val["user_name"] = arr[0]
	//	}
	//
	//}

	c.JSON(http.StatusOK, gin.H{
		"info": gin.H{
			"datalist": accountList,
			"total":    total,
		},
		"code": http.StatusOK,
		"msg":  model.ReqSuccess,
	})
}
