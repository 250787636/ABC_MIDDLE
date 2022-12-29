package api

import (
	"crypto/md5"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/utils"
	"example.com/m/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// 验证用户的用户名,查找或者创建
func CheckAccount(c *gin.Context) {
	// 数据绑定所需要的变量
	var user model.User
	// 如果存在进行获取数据并更新
	var department model.Departments
	// 验证存储信息是否有效
	allow_1 := false
	allow_2 := false
	// 获取当前时间
	forTime := time.Now().Format("2006-01-02 15:04:05")

	// 数据绑定
	id, idOK := c.GetPostForm("user_id")
	name, nameOK := c.GetPostForm("name")
	departmentName, departmentNameOK := c.GetPostForm("department_name")
	jobTitle, jobTitleOK := c.GetPostForm("job_title")
	accountLevel, accountLevelOK := c.GetPostForm("account_level")
	clentToken, _ := c.GetPostForm("token")

	// 当为内置超级管理员时进行跳过验证
	if name == "root" && departmentName == "0" && jobTitle == "0" && accountLevel == "超级管理员" {
		if err := app.DB.Model(model.User{}).Where("user_name = ?", name).First(&user).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"info": gin.H{},
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		} else {
			if err := app.DB.Model(model.User{}).Where("user_name = ?", name).Update("last_login_time", forTime).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"info": gin.H{},
					"code": http.StatusInternalServerError,
					"err":  err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"info": gin.H{
					"id":             user.ID,
					"name":           user.UserName,
					"token":          user.Token,
					"is_admin":       user.IsAdmin,
					"is_super_admin": user.AccountLevel == "超级管理员",
				},
				"code": http.StatusOK,
				"msg":  model.ReqSuccess,
			})
			return
		}
	}
	//  普通用户及管理员用户进行判空
	nameOK = utils.IsStringEmpty(name, nameOK)
	departmentNameOK = utils.IsStringEmpty(departmentName, departmentNameOK)
	jobTitleOK = utils.IsStringEmpty(jobTitle, jobTitleOK)
	accountLevelOK = utils.IsStringEmpty(accountLevel, accountLevelOK)
	idOK = utils.IsStringEmpty(id, idOK)
	if !(nameOK && departmentNameOK && jobTitleOK && accountLevelOK && idOK) {
		response.FailResult(http.StatusInternalServerError, "", model.ReqParameterMissing, c)
		return
	}

	//部门部分
	department.DepartmentName = departmentName
	// 如果不存在部门则创建
	if err := app.DB.Model(model.Departments{}).Where("department_name = ?", department.DepartmentName).FirstOrCreate(&department).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.CreateDepartmentFail,
		})
		return
	} else {
		log.Info("此部门已存在,验证通过")
		allow_1 = true
	}

	// 用户账号部分
	user.UserName = id + "-" + name
	user.JobTitle = jobTitle
	user.AccountLevel = accountLevel
	user.LastLoginTime = forTime
	user.DepartmentID = department.ID
	user.IsAdmin = false
	if user.AccountLevel == "部门管理员" || user.AccountLevel == "超级管理员" {
		user.IsAdmin = true
	}

	//生成token
	uui, _ := uuid.NewUUID()
	orgToken := []byte(user.UserName + uui.String() + forTime)
	m5Token := md5.Sum(orgToken)
	token := fmt.Sprintf("%x", m5Token)
	user.Token = token
	user.ClientToken = clentToken

	//未存在此账号 则创建
	if err := app.DB.Model(model.User{}).Where("user_name = ? AND department_id = ?", user.UserName, user.DepartmentID).FirstOrCreate(&user).Error; err != nil {
		log.Info(err)
	} else {
		log.Info("此账号已存在,验证通过")
		// 当发现数据库与生成的token不同时 进行token更新
		if user.Token != token {
			if err := app.DB.Model(model.User{}).
				Where("id = ?", user.ID).
				Update("token", token).
				First(&user).
				Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusInternalServerError,
					"err":  model.CreateAccountFail,
				})
				return
			}
		}
		allow_2 = true
	}
	if allow_1 && allow_2 {
		c.JSON(http.StatusOK, gin.H{
			"info": gin.H{
				"id":             user.ID,
				"name":           user.UserName,
				"token":          user.Token,
				"is_admin":       user.IsAdmin,
				"is_super_admin": user.AccountLevel == "超级管理员",
			},
			"code": http.StatusOK,
			"msg":  model.ReqSuccess,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"info": gin.H{},
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
	}

}

// 1.3修改用户信息接口
func ModifyAccount(c *gin.Context) {
	// 需要使用的结构体
	var user model.User
	var department model.Departments
	// 判断是否为管理员 或 超级管理员
	var isAdminNum int
	// 获取参数
	name, nameOk := c.GetPostForm("name")
	oldDepartmentName, oldDepartmentNameOK := c.GetPostForm("old_department_name")
	newJobTitle, newJobTitleOk := c.GetPostForm("new_job_title")
	newAccountLevel, newAccountLevelOK := c.GetPostForm("new_account_level")
	newDepartmentName, newDepartmentOk := c.GetPostForm("new_department_name")

	if newAccountLevel == "超级管理员" {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  "无法升级为超级管理员",
		})
		return
	}
	//判空
	nameOk = utils.IsStringEmpty(name, nameOk)
	oldDepartmentNameOK = utils.IsStringEmpty(oldDepartmentName, oldDepartmentNameOK)
	newJobTitleOk = utils.IsStringEmpty(newJobTitle, newJobTitleOk)
	newAccountLevelOK = utils.IsStringEmpty(newAccountLevel, newAccountLevelOK)
	newDepartmentOk = utils.IsStringEmpty(newDepartmentName, newDepartmentOk)

	if !(nameOk && oldDepartmentNameOK) {
		response.FailResult(http.StatusInternalServerError, "", model.ReqParameterMissing, c)
		return
	}

	// 通过部门名获取部门id
	if err := app.DB.Where("department_name = ?", oldDepartmentName).
		First(&department).Error; err != nil {
		response.FailResult(http.StatusInternalServerError, "", err.Error(), c)
		return
	}

	// 通过用户名与部门id获取当前用户
	if err := app.DB.
		Where("user_name = ? AND department_id = ?", name, department.ID).
		First(&user).Error; err != nil {
		response.FailResult(http.StatusInternalServerError, "", err.Error(), c)
		return
	}

	// 若存在非必传参数,则添加到原数据中并进行更新
	if newAccountLevelOK {
		isAdminNum = 0
		user.IsAdmin = false
		user.AccountLevel = newAccountLevel
		if newAccountLevel == "部门管理员" || newAccountLevel == "超级管理员" {
			isAdminNum = 1
			user.IsAdmin = true
		}
	}
	if newJobTitleOk {
		user.JobTitle = newJobTitle
	}
	if newDepartmentOk {
		var newDep model.Departments
		newDep.DepartmentName = newDepartmentName
		if err := app.DB.Where("department_name = ?", newDepartmentName).FirstOrCreate(&newDep).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
		user.DepartmentID = newDep.ID
	}

	if !(newAccountLevelOK || newJobTitleOk || newDepartmentOk) {
		c.JSON(http.StatusOK, gin.H{
			"info": gin.H{
				"is_admin": user.IsAdmin,
				"name":     user.UserName,
				"token":    user.Token,
			},
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterModFail,
		})
		return
	}

	// 修改传入的用户
	if err := app.DB.Where("id = ?", user.ID).
		Updates(&user).
		Update("is_admin", isAdminNum).
		Error; err != nil {
		response.FailResult(http.StatusInternalServerError, "", err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"info": gin.H{
			"is_admin": user.IsAdmin,
			"name":     user.UserName,
			"token":    user.Token,
		},
		"code": http.StatusOK,
		"msg":  model.ModSuccess,
	})
}

// 1.4获取部门接口
func GetDepartment(c *gin.Context) {
	type departments struct {
		ID             int    `json:"id"`
		DepartmentName string `json:"department_name"`
	}

	var departmentList []departments
	if err := app.DB.Find(&departmentList).Error; err != nil {
		response.FailResult(http.StatusInternalServerError, "", err.Error(), c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"info": gin.H{
			"datalist": departmentList,
		},
		"code": http.StatusOK,
		"msg":  model.ReqSuccess,
	})

}
