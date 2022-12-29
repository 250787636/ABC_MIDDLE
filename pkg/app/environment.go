package app

import (
	"example.com/m/pkg/lib/minios"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	DB                *gorm.DB
	MinioClient       *minios.Client
	PasswordPublicKey string
)

var (
	// 本地环境变量 env
	Env         = env()
	RootDir     = rootDir()
	ResourceDir = RootDir + "/resources"
	ToolDir     = ResourceDir + "/tool"
	StaticDir   = ResourceDir + "/static"
)

// Env 的值为dev或空
// 1.当为空时加载deploy.ini和customer.ini配置
// 2.当为dev时加载develop.ini和customer.ini覆盖配置
func env() string {
	env := os.Getenv("Env")
	switch env {
	case "dev":
		return "dev"
	default:
		return ""
	}
}

// 1.如果为单元测试 获取文件的绝对路径
// 2.如果不是 则用 . 拼接
func rootDir() string {
	if isUnitTestEnv() {
		_, filePath, _, _ := runtime.Caller(0)
		for i := 0; i < 3; i++ { //循环次数根据当前文件距项目根目录的层级
			filePath = filepath.Dir(filePath)
		}
		return filePath
	}
	return "."
}

// 判断是否为单元测试
func isUnitTestEnv() bool {
	return len(os.Args) > 1 && strings.Contains(os.Args[1], "-test")
}
