package api

import (
	"archive/zip"
	"bytes"
	"errors"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/utils"
	"example.com/m/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 获取文件信息接口
func GetFile(c *gin.Context) {
	var data struct {
		FileName string `json:"file_name"`
		FileSize int64  `json:"file_size"`
		FilePath string `json:"file_path"`
	}
	useType, useTypeOk := c.GetPostForm("use_type")
	useTypeOk = utils.IsStringEmpty(useType, useTypeOk)
	if !useTypeOk {
		response.FailResult(http.StatusInternalServerError, "", model.ReqParameterMissing, c)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		response.FailResult(http.StatusInternalServerError, "", err.Error(), c)
		return
	}

	// 保存文件到本地
	fileUuid := uuid.New()
	// 拼接路径保存文件
	filePathString := fmt.Sprintf("media/%s/%s/%s", useType, fileUuid, file.Filename)
	err = os.MkdirAll(filepath.Dir(filePathString), os.ModePerm)
	if err != nil {
		response.FailResult(http.StatusInternalServerError, "", err.Error(), c)
		return
	}
	if err := c.SaveUploadedFile(file, filePathString); err != nil {
		response.FailResult(http.StatusInternalServerError, "", err.Error(), c)
		return
	}

	data.FileName = file.Filename
	data.FileSize = file.Size
	data.FilePath = filePathString

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"info": data,
		"msg":  model.ReqSuccess,
	})
}

// DownloadFile 批量下载源文件
func DownloadFile(c *gin.Context) {

	taskIdString := c.Query("file")
	fileType := c.Query("type")
	str1 := strings.ReplaceAll(taskIdString, "[", "")
	str2 := strings.ReplaceAll(str1, "]", "")
	idArray := strings.Split(str2, ",")

	var TaskInfo []model.CePingUserTask
	if err := app.DB.Debug().Model(&model.CePingUserTask{}).Where("task_id in (?)", idArray).Find(&TaskInfo).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 2.准备zip
	out, err := os.Create("test.zip")
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
		}
	}(out)

	writer := zip.NewWriter(out)

	//sufFIx := ".apk"
	//switch fileType {
	//case "apk":
	//	sufFIx = ".apk"
	//case "ios":
	//	sufFIx = ".ipa"
	//case "mp":
	//	sufFIx = ".jpg"
	//case "sdk":
	//	sufFIx = ".aar"
	//}

	var TaskMap = make(map[string]string)

	// 1.把应用名去后缀
	for _, task := range TaskInfo {
		if fileType == "sdk" {
			_, ok := TaskMap[task.PkgName]
			// 如果没有就放进去
			if !ok {
				TaskMap[task.PkgName] = task.FilePath
			}
		} else {
			_, ok := TaskMap[task.AppName]
			// 如果没有就放进去
			if !ok {
				TaskMap[task.AppName] = task.FilePath
			}
		}
	}

	for key, task := range TaskMap {
		//appName := strings.Split(key, ".")
		//fmt.Println("task", task.FilePath)
		fileWriter, err := writer.Create(key)
		if err != nil {
			if os.IsPermission(err) {
				log.Error("err:", err.Error())
				response.FailWithMessage(err.Error(), c)
				return
			}
			log.Error("Create file %s error: %s\n", key, err.Error())
			response.FailWithMessage(err.Error(), c)

			return
		}

		fileInfo, err := os.Open(task)
		if err != nil {
			log.Error("Open file %s error: %s\n", task, err.Error())
			response.FailWithMessage(err.Error(), c)

			return
		}
		fileBody, err := ioutil.ReadAll(fileInfo)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			log.Error("Read file %s error: %s\n", task, err.Error())
			return
		}

		_, err = fileWriter.Write(fileBody)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			log.Error("err:", err.Error())
			return
		}
	}

	if err := writer.Close(); err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)

		fmt.Println("Close error: ", err)
		return
	}
	fileContentDisposition := "inline;filename=测评源文件下载.zip"
	c.Header("Content-Type", "application/zip") // 这里是压缩文件类型 .zip
	c.Header("Content-Disposition", fileContentDisposition)

	c.File("test.zip")

}
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

// NewExcel 导入excel
func NewExcel(c *gin.Context) {
	excelFile, err := c.FormFile("file")
	if err != nil {
		log.Error(err.Error())
		response.FailWithMessage("获取上传文件失败", c)
		return
	}

	// 保存文件到本地
	fileUuid := uuid.New()
	// 拼接路径保存文件
	filePathString := fmt.Sprintf("media/%s/%s", fileUuid, excelFile.Filename)
	err = os.MkdirAll(filepath.Dir(filePathString), os.ModePerm)
	if err != nil {
		log.Error(err.Error())
		response.FailWithMessage("保存文件失败", c)
		return
	}

	if err := c.SaveUploadedFile(excelFile, filePathString); err != nil {
		log.Error(err.Error())
		response.FailWithMessage("保存文件失败", c)
		return
	}

	f, err := excelize.OpenFile(filePathString)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Error(err.Error())
		}
	}()

	// 0.获取总数据
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Error(err.Error())
		return
	}

	//fmt.Println("len", len(rows))

	// 1.获取excel有效行数
	effectiveRow := 0
	for _, row := range rows {
		if len(row) != 9 {
			continue
		}
		//fmt.Printf("第 %d 行有%d个数据\n", rowIn, len(row))
		effectiveRow++
	}
	// 2.创建空二维数组
	temArray := make([][]string, effectiveRow)
	lieTotal := len(rows[0])
	for i := range temArray {
		temArray[i] = make([]string, lieTotal)
	}
	//fmt.Println("1111", temArray)
	temRow := 0
	// 3.将excel信息转存到二维数组
	for _, row := range rows {
		// 如果该行数据为1表示只有序号值，就不对该行进行获取
		if len(row) != 9 {
			continue
		}
		//fmt.Printf("第 %d 行有%d个数据\n", rowIn, len(row))

		for i := 0; i < len(row); i++ {
			//fmt.Printf("第%d个值为 %s\n", i, row[i])
			temArray[temRow][i] = row[i]
		}
		temRow++
	}
	//fmt.Println("获取到的数据", temArray)

	// 获取加固策略以及应用类型的map字典集
	_, _, policyVK, appTypeVK, idExist, useUserExist, err := getDictionaryList("import")
	if err != nil {
		log.Error(err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 导入数据条数 比 数据库中的数据条数多时  -1为减去标题数据
	var msg string
	// 4.获取二维数组信息
	for rowIndex, row := range temArray {
		//跳过第一行表头信息
		if rowIndex == 0 {
			continue
		}
		//遍历每一个单元
		application := model.Application{}
		intId, _ := strconv.Atoi(row[0])
		application.AppName = row[1]
		application.ModelName = row[2]
		application.AppCnName = row[3]
		application.AppVersion = row[4]
		application.ModelCnName = row[5]
		application.RecommendPolicy = policyVK[row[6]]
		application.AppTypeID = appTypeVK[row[7]]
		application.LastChangeTime = time.Now().Format("2006-01-02 15:04:05")
		application.TheApp = application.AppName + "-" + application.AppCnName
		application.TheModel = application.ModelName + "-" + application.ModelCnName
		idExist[intId] = true
		msg = "无操作"
		// 关联用户存在 id存在
		//useUserExist[intId] == true &&
		var testApp model.Application
		if idExist[intId] == true {
			err := app.DB.Model(model.Application{}).Where("id = ?", intId).First(&testApp).Error
			if err != nil {
				msg = fmt.Sprintf("存储数据编号:%s成功", row[0])
				// 进行create操作
				application.ID = uint(intId)
				if err := app.DB.Model(model.Application{}).Create(&application).Error; err != nil {
					log.Error(err.Error())
					msg = fmt.Sprintf("存储数据编号:%s失败", row[0])
				}
			} else {
				//进行updates
				msg = fmt.Sprintf("进行数据更新,编号:%s成功", row[0])
				if err := app.DB.Model(model.Application{}).Where("id = ?", uint(intId)).Updates(&application).Error; err != nil {
					log.Error(err.Error())
					msg = fmt.Sprintf("更新数据编号:%s失败", row[0])
				}
			}
		} else if useUserExist[intId] == true && idExist[intId] == false { // 关联用户存在 id不存在
			// 错误 该应用已关联任务无法删除
			msg = fmt.Sprintf("删除数据编号:%d删除失败,该应用已存在关联信息", row[0])
		}
		log.Info(msg)
	}
	// 导入数据条数比数据库中数据条数少时
	if len(temArray)-1 < len(idExist) {
		for k, v := range idExist {
			if !v && useUserExist[k] == true {
				msg = fmt.Sprintf("删除数据编号:%d删除失败,该应用已存在关联信息", k)
				log.Info(msg)
				continue
			}
			if !v {
				msg = fmt.Sprintf("删除数据编号:%d成功", k)
				// 进行删除
				if err := app.DB.Model(model.Application{}).Where("id = ?", uint(k)).Unscoped().Delete(&model.Application{
					Model: gorm.Model{ID: uint(k)},
				}).Error; err != nil {
					log.Error(err.Error())
					msg = fmt.Sprintf("删除数据编号:%d失败", k)
				}
				log.Info(msg)
			}
		}
	}
	response.OkWithMessage("导入excel文件成功", c)
}

//  记录导出功能
func ApplicationExporting(c *gin.Context) {
	var appList []model.Application
	var fileName string
	var err error

	// 获取所有应用
	if err := app.DB.Model(model.Application{}).Find(&appList).Error; err != nil {
		log.Error(err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 获取加固策略以及应用类型的map字典集
	policyKV, appTypeKV, _, _, _, _, err := getDictionaryList("export")
	if err != nil {
		log.Error(err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	//导出
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	row.AddCell().Value = "序号"
	row.AddCell().Value = "应用系统"
	row.AddCell().Value = "模块"
	row.AddCell().Value = "系统中文全称"
	row.AddCell().Value = "系统英文简称"
	row.AddCell().Value = "模块中文全称"
	row.AddCell().Value = "推荐加固策略"
	row.AddCell().Value = "应用类型"
	row.AddCell().Value = "录入时间"
	for _, v := range appList {
		row = sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(int(v.ID))
		row.AddCell().Value = v.TheApp
		row.AddCell().Value = v.TheModel
		row.AddCell().Value = v.AppCnName
		row.AddCell().Value = v.AppVersion
		row.AddCell().Value = v.ModelCnName
		row.AddCell().Value = policyKV[v.RecommendPolicy]
		row.AddCell().Value = appTypeKV[v.AppTypeID]
		row.AddCell().Value = v.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	buf := new(bytes.Buffer)
	err = file.Write(buf)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	fileName = "应用管理表" + time.Now().Format("2006-01-02 15:04:05") + ".xlsx"
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, url.QueryEscape(fileName)))
	c.Data(http.StatusOK, "text/xlsx", buf.Bytes())
}

// 获取所有加固策略，获取所有应用类型
func getDictionaryList(key string) (map[int]string, map[int]string, map[string]int, map[string]int, map[int]bool, map[int]bool, error) {
	// 安卓策略
	var policyListAndroid []model.JiaguPolicyAndroid
	// h5策略
	var policyListH5 []model.JiaguPolicyH5
	var appTypeList []model.ApplicationType
	var application []model.Application

	// 获取所有加固策略
	if err := app.DB.Model(model.JiaguPolicyAndroid{}).Find(&policyListAndroid).Error; err != nil {
		log.Error(err.Error())
		return nil, nil, nil, nil, nil, nil, err
	}
	if err := app.DB.Model(model.JiaguPolicyH5{}).Find(&policyListH5).Error; err != nil {
		log.Error(err.Error())
		return nil, nil, nil, nil, nil, nil, err
	}
	// 获取所有应用类型
	if err := app.DB.Model(model.ApplicationType{}).Find(&appTypeList).Error; err != nil {
		log.Error(err.Error())
		return nil, nil, nil, nil, nil, nil, err
	}
	if err := app.DB.Model(model.Application{}).Find(&application).Error; err != nil {
		log.Error(err.Error())
		return nil, nil, nil, nil, nil, nil, err
	}

	if key == "import" { // 导入接口
		// 生成加固策略字典集 vk用于 name 对 id
		policyVK := make(map[string]int, len(policyListAndroid)+len(policyListH5))
		for _, val := range policyListAndroid {
			policyVK[val.Name] = val.Id
		}
		for _, val := range policyListH5 {
			policyVK[val.Name] = val.Id
		}

		// 生成应用类型字典集 vk用于 name 对 id
		appTypeVK := make(map[string]int, len(appTypeList))
		for _, val := range appTypeList {
			appTypeVK[val.AppType] = int(val.ID)
		}

		// 获取数据库中的所有id
		idExist := make(map[int]bool, len(appTypeList))
		for _, val := range application {
			// 初始值为false
			idExist[int(val.ID)] = false
		}
		// 获取是否关联使用用户useUser
		useUserExist := make(map[int]bool, len(appTypeList))
		for _, val := range application {
			if val.UseUser != "" {
				useUserExist[int(val.ID)] = true
			}
		}
		return nil, nil, policyVK, appTypeVK, idExist, useUserExist, nil

	} else if key == "export" { // 导出接口
		// 生成加固策略字典集 kv 用于 id 对 name
		policyKV := make(map[int]string, len(policyListAndroid)+len(policyListH5))
		for _, val := range policyListAndroid {
			policyKV[val.Id] = val.Name
		}
		for _, val := range policyListH5 {
			policyKV[val.Id] = val.Name
		}

		// 生成应用类型字典集 kv 用于 id 对 name
		appTypeKV := make(map[int]string, len(appTypeList))
		for _, val := range appTypeList {
			appTypeKV[int(val.ID)] = val.AppType
		}

		return policyKV, appTypeKV, nil, nil, nil, nil, nil

	}
	return nil, nil, nil, nil, nil, nil, errors.New("参数错误")
}
