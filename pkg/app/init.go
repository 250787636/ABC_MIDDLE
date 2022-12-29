package app

import (
	"bytes"
	"crypto/md5"
	"example.com/m/pkg/lib/errors"
	"example.com/m/pkg/lib/minios"
	"example.com/m/pkg/log"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func init() {
	chmodToolFile()
	InitDB()
	//initMinio()
	//GetCpToken()
}

// 修改所有工具的可执行权限
func chmodToolFile() {
	toolPath := ToolDir
	executeFiles := []string{
		toolPath + "/aapt2",
		toolPath + "/phantomjs",
	}
	for _, file := range executeFiles {
		err := os.Chmod(file, os.ModePerm)
		if err != nil {
			log.Error(err)
		}
	}
}

// InitDB 初始化数据库
func InitDB() {
	conf := Conf.Mysql

	if err := CreateDatabaseIfNotExist(conf.Username, conf.Password,
		conf.Host, conf.Port, conf.Database); err != nil {
		log.Warn(errors.WithCaller(err, errors.M{"config": conf}, "创建数据库失败"))
	}
	var err error
	DB, err = NewGormDb(conf.Username, conf.Password, conf.Host,
		conf.Port, conf.Database)
	if err != nil {
		log.Fatal(errors.WithCaller(err, errors.M{"config": conf}, "数据库连接失败"))
	}

	if conf.ResetDb == true {
		fmt.Println("是否删除数据库中的所有表", conf.ResetDb)
		sql := "drop database " + conf.Database
		if err := DB.Debug().Exec(sql).Error; err != nil {
			panic(err)
		}

		if err := CreateDatabaseIfNotExist(conf.Username, conf.Password,
			conf.Host, conf.Port, conf.Database); err != nil {
			log.Fatal(errors.WithCaller(err, errors.M{"config": conf}, "创建数据库失败"))
		}
		var err error
		DB, err = NewGormDb(conf.Username, conf.Password, conf.Host,
			conf.Port, conf.Database)
		if err != nil {
			log.Fatal(errors.WithCaller(err, errors.M{"config": conf}, "数据库连接失败"))
		}

		defaultConf := RootDir + "/deploy.ini"
		path_name := defaultConf
		fmt.Println("111", path_name)
		cfg, err := ini.Load(path_name)
		if err != nil {
			fmt.Printf("Fail to read file: %v", err)
			os.Exit(1)
		}
		cfg.Section("mysql").Key("reset_db").SetValue("false")
		err = cfg.SaveTo("deploy.ini")
		if err != nil {
			panic(err)
		}
	}
}

// 初始化minio
func initMinio() {
	var err error
	MinioClient, err = minios.NewMinioClient(Conf.Minio.Endpoint, Conf.Minio.AccessKeyID,
		Conf.Minio.SecretAccessKey, Conf.Minio.UseSSL)
	if err != nil {
		log.Fatal(err)
	}
}

// 开启gorm连接
func NewGormDb(user, password, host, port, database string) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名,启用该选项后，`User`表将是`user`
		},
		Logger: logger.New(
			log.DefalutLogger,
			logger.Config{
				SlowThreshold: time.Second,   // 慢 SQL 阈值
				LogLevel:      logger.Silent, // Log level
			}),
	}
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		user, password, host, port, database)
	db, err := gorm.Open(mysql.Open(url), gormConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 数据库创建与否
func CreateDatabaseIfNotExist(user, password, host, port, database string) error {
	db, err := NewGormDb(user, password, host, port, "")
	if err != nil {
		return err
	}
	createDatabase := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", database)
	if err = db.Exec(createDatabase).Error; err != nil {
		return err
	}
	return nil
}

// GetCpToken 初始化的时候获取测评服务的token
func GetCpToken(userName, password, ip string) (token, signature string, err error) {

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	paramMap := make(map[string]interface{})
	paramMap["username"] = userName
	paramMap["password"] = password

	value, err := jsoniter.Marshal(paramMap)
	if err != nil {
		log.Error(err.Error())
		return "nil", "nil", err
	}

	err = writer.WriteField("param", string(value))
	if err != nil {
		log.Error(err.Error())
		return "nil", "nil", err
	}

	err = writer.Close()
	if err != nil {
		log.Error(err.Error())
		return "nil", "nil", err
	}

	clientURL := ip + "/v4/apply_auth"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		log.Error(err.Error())
		return "nil", "nil", err
	}

	//注意别忘了设置header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	//Do方法发送请求
	resp, err := client.Do(request)
	if err != nil {
		log.Error(err.Error())
		return "", "", err
	}
	defer resp.Body.Close()

	post, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("初始化的时候获取测评服务的token失败", err.Error())
		return "nil", "nil", err
	}

	var reStruct struct {
		AppKey    string `json:"appkey"`
		AppSecret string `json:"appsecret"`
		Msg       string `json:"msg"`
		State     int    `json:"state"`
	}

	err = jsoniter.Unmarshal(post, &reStruct)
	if err != nil {
		log.Error(string(post))
		log.Error(reStruct)
		log.Error(err.Error())
		return "", "", err
	}
	if reStruct.State != 200 {
		log.Error("调用测评平台获取sign接口失败")
		fmt.Println("---", reStruct.Msg)
		return "", "", errors.New(reStruct.Msg)
	}

	str := userName + reStruct.AppKey + reStruct.AppSecret
	st := fmt.Sprintf("%x", md5.Sum([]byte(str)))

	token, err = GetRealToken(userName, reStruct.AppKey, st, ip)
	if err != nil {
		log.Error(err.Error())
		return "nil", "nil", err
	}

	return token, st, nil
}

// GetRealToken 1.2.生成并获取AccessToken接口
func GetRealToken(userName, appKey, signature, ip string) (token string, err error) {
	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	paramMap := make(map[string]interface{})
	paramMap["username"] = userName
	paramMap["appkey"] = appKey
	paramMap["signature"] = signature
	//fmt.Println("username ", userName)
	//fmt.Println("appkey ", appKey)
	//fmt.Println("signature", signature)

	value, err := jsoniter.Marshal(paramMap)
	if err != nil {
		log.Error(err)
		return "nil", err
	}

	err = writer.WriteField("param", string(value))
	if err != nil {
		log.Error(err)
		return "nil", err
	}

	err = writer.Close()
	if err != nil {
		log.Error(err)
		return "nil", err
	}

	clientURL := ip + "/v4/apply_access_token"

	//生成post请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", clientURL, buff)
	if err != nil {
		log.Error(err)
		return "nil", err
	}

	//注意别忘了设置header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	//Do方法发送请求
	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return "nil", err
	}
	defer resp.Body.Close()

	post, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(string(post))
		log.Error(err)
		return "nil", err
	}

	var reStruct struct {
		Msg         string `json:"msg"`
		State       int    `json:"state"`
		AccessToken string `json:"accesstoken"`
	}

	err = jsoniter.Unmarshal(post, &reStruct)
	if err != nil {
		log.Error(err)
		log.Error(string(post))
		return "", err
	}

	if reStruct.State != 200 {
		log.Error("调用测评平台获取token接口失败")
		return "", errors.New(reStruct.Msg)
	}
	return reStruct.AccessToken, nil
}
