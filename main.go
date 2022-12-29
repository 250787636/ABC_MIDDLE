package main

import (
	"example.com/m/model/automigrate"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/router"
	"example.com/m/router/api/ceping"
	"example.com/m/router/api/jiagu"
	"flag"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func main() {
	flag.Parse()
	// 将请求放入控制台与run.log中
	gin.DefaultWriter = io.MultiWriter(os.Stdout, log.FileWriter)
	r := gin.Default()
	// 进行建表操作
	automigrate.Run()
	// 进行服务注册操作
	router.Init(r)
	// 过5分钟进行加固策略更新
	go jiagu.GetPolicyListByType("android加固")
	go jiagu.GetPolicyListByType("h5加固")
	// 过1个半小时重新获取测评平台token
	go ceping.TimeToGetToken()
	log.Info(`                               
		  /$$$$$$  /$$$$$$$   /$$$$$$ 
		 /$$__  $$| $$__  $$ /$$__  $$
		| $$  \ $$| $$  \ $$| $$  \__/
		| $$$$$$$$| $$$$$$$ | $$      
		| $$__  $$| $$__  $$| $$      
		| $$  | $$| $$  \ $$| $$    $$
		| $$  | $$| $$$$$$$/|  $$$$$$/
		|__/  |__/|_______/  \______/
	`)
	log.Fatal(r.Run(":" + app.Conf.System.Port))
}
