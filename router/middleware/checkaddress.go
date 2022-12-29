package middleware

import (
	"example.com/m/pkg/app"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
)

// 中间件检测用户地址ip
func CheckAddress(r *gin.Context) {
	var ipList []string
	ip := ClientIP(r.Request)
	//172.16.102.55
	//172.*
	//172.16.*
	//172.16.102.*
	ipOnce, ipTwice, ipThrice := GetIpType(ip)
	// 权限判断符 当没有权限时进行拦截 不让其进行进一步操作
	isAllow := true
	if app.Conf.Platforms.IP == "*" {
		isAllow = false
		goto END
	}
	ipList = strings.Split(app.Conf.Platforms.IP, ",")

	for _, v := range ipList {
		// 如果保护 全部符号
		if strings.Contains(v, "*") {
			if v == "*" {
				isAllow = false
				goto END
			}
			// 以点拆分配置文件中的ip 看存在几位数
			configIpList := strings.Split(v, ".")
			switch len(configIpList) {
			case 2:
				if v == ipOnce {
					isAllow = false
					break
				}
			case 3:
				if v == ipTwice {
					isAllow = false
					break
				}

			case 4:
				if v == ipThrice {
					isAllow = false
					break
				}
			}
		} else {
			if v == ip {
				isAllow = false
			}
		}
	}
END:
	if isAllow {
		r.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusLocked, // 没有权限
			"err":  "当前电脑ip没有访问权限",
		})
		r.Abort()
	}
}

// 获取客户端IP
// 解析X-Real-Ip和X-Forwarded-For
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// 获取 包含*的类种
func GetIpType(ip string) (ipOnce, ipTwice, ipThrice string) {
	// 以点分割 找出当前ip存在*的情况
	computerIPList := strings.Split(ip, ".")
	for i, _ := range computerIPList {
		switch i {
		case 0:
			ipOnce += computerIPList[i] + ".*"
		case 1:
			ipTwice += computerIPList[i-1] + "."
			ipTwice += computerIPList[i] + ".*"
		case 2:
			ipThrice += computerIPList[i-2] + "."
			ipThrice += computerIPList[i-1] + "."
			ipThrice += computerIPList[i] + ".*"
		}
	}
	return ipOnce, ipTwice, ipThrice
}
