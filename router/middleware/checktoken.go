package middleware

import (
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckToken 验证 Header 上的 token
func CheckToken(r *gin.Context) {
	var user model.User
	var token string
	isRight := false

	token = r.Request.Header.Get("token")

	if token == "" {
		token = r.Query("token")
		if token == "" {
			goto END
		}
	}
	if err := app.DB.Where("token = ?", token).First(&user).Error; err != nil {
		log.Info(err)
	} else {
		isRight = true
		r.Set("userId", user.ID)
		r.Set("isAdmin", user.IsAdmin)
		r.Set("departmentId", user.DepartmentID)
		r.Set("superAdmin", user.AccountLevel == "超级管理员")
		log.Info("token验证成功")
	}
END:
	if !isRight {
		r.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusForbidden, // 没有访问权限
			"err":  "token不存在",
		})
		r.Abort()
	}
}
