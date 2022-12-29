package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 管理员账户独享权限
func AdminOnly(c *gin.Context) {
	isAdmin, ok := c.Get("isAdmin")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  "当前用户缺失账户类型参数",
		})
		return
	}
	if !isAdmin.(bool) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  "当前用户为普通用户,没有访问此功能权限",
		})
		c.Abort()
	}
}
