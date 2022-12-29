package jiagu

import (
	"bytes"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/utils"
	"example.com/m/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 以部门为主体进行数据统计
func JiaGuDepartmentCount(c *gin.Context) {
	//model. JiaGuDepartmentCount{} // 当前接口返回的接口体所有类型数据
	// 需要返的数据
	var data []map[string]interface{}
	// 计数值
	var count int64
	// 开始和结束日期
	var startTime time.Time
	var endTime time.Time
	var err error

	// 获取前端传入的参数
	depName, depNameOK := c.GetPostForm("department_name")
	Time1, startTimeOK := c.GetPostForm("start_time")
	Time2, endTimeOk := c.GetPostForm("end_time")
	// 排序两参数
	sortBy, sortByOK := c.GetPostForm("sort_by")
	sortField, sortFieldOK := c.GetPostForm("sort_field")
	appType, appTypeOK := c.GetPostForm("app_type")
	// 必传加固应用类型
	if !appTypeOK {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"err":  model.ReqParameterMissing,
		})
		return
	}
	// 判断字符是否为空 如果为空则当该值不存在
	depNameOK = utils.IsStringEmpty(depName, depNameOK)
	startTimeOK = utils.IsStringEmpty(Time1, startTimeOK)
	endTimeOk = utils.IsStringEmpty(Time2, endTimeOk)
	sortByOK = utils.IsStringEmpty(sortBy, sortByOK)
	sortFieldOK = utils.IsStringEmpty(sortField, sortFieldOK)

	// 加固类型 赋值对应aoolication_type的id
	var appTypeNum int
	switch appType {
	case "android":
		appTypeNum = 1
		break
	case "h5":
		appTypeNum = 2
		break
	default:
		response.FailResult(http.StatusInternalServerError, "", model.ReqParameterMissing, c)
		return
	}

	// 将string转为datetime
	if startTimeOK && endTimeOk {
		startTime, err = time.Parse("2006-01-02 15:04", Time1)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
		endTime, err = time.Parse("2006-01-02 15:04", Time2)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
	}

	// 拼接sql
	var mainSqlString bytes.Buffer
	layout := "2006-01-02 15:04:05"
	// 获取主数据部分
	if startTimeOK && endTimeOk {
		mainSqlString.WriteString(" created_at BETWEEN ")
		mainSqlString.WriteString("'")
		mainSqlString.WriteString(startTime.Format(layout))
		mainSqlString.WriteString("'")
		mainSqlString.WriteString(" AND ")
		mainSqlString.WriteString("'")
		mainSqlString.WriteString(endTime.Format(layout))
		mainSqlString.WriteString("'")
	}
	// 部门名称
	if depNameOK { // 当有查询参数时
		if mainSqlString.String() != "" {
			mainSqlString.WriteString(" AND ")
		}
		mainSqlString.WriteString("departments.department_name like '")
		mainSqlString.WriteString("%" + depName + "%")
		mainSqlString.WriteString("'")
	}

	// 无查询条件时
	if mainSqlString.String() == "" {
		if err := app.DB.Model(model.Departments{}).Find(&data).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"err":  err.Error,
			})
			return
		}
	} else {
		if err := app.DB.Model(model.Departments{}).
			Where(mainSqlString.String()).
			Find(&data).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"err":  err.Error,
			})
			return
		}
	}

	// 当有日期限制时 遍历添加对应数据
	if startTimeOK && endTimeOk {
		for i := 0; i < len(data); i++ {
			// 统计用户
			if err := app.DB.Model(model.User{}).Where("department_id = ? AND created_at BETWEEN ? AND ?", data[i]["id"], startTime, endTime).Count(&count).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
					"err":  err.Error,
				})
				return
			}
			data[i]["user_count"] = count

			// 统计任务数
			if err := app.DB.Table("departments").
				Joins("INNER JOIN `user` ON `user`.department_id = departments.id").
				Joins("INNER JOIN jia_gu_task ON jia_gu_task.user_id = `user`.id").
				Where("departments.id = ? AND jia_gu_task.created_at BETWEEN ? AND ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], startTime, endTime, appTypeNum).
				Count(&count).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
					"err":  err.Error,
				})
				return
			}
			data[i]["task_count"] = count

			// 统计使用应用数
			if err := app.DB.Table("departments").
				Distinct("jia_gu_task.app_id").
				Joins("INNER JOIN `user` ON `user`.department_id = departments.id").
				Joins("INNER JOIN jia_gu_task ON jia_gu_task.user_id = `user`.id").
				Joins("INNER JOIN application_type ON application_type.id = jia_gu_task.app_type_id").
				Where("departments.id = ? AND jia_gu_task.created_at BETWEEN ? AND ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], startTime, endTime, appTypeNum).
				Count(&count).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
					"err":  err.Error,
				})
				return
			}
			data[i]["application_count"] = count

			// 使用标准加固策略任务数
			if err := app.DB.Table("departments").
				Distinct("jia_gu_task.id").
				Joins("INNER JOIN `user` ON `user`.department_id = departments.id").
				Joins("INNER JOIN jia_gu_task ON jia_gu_task.user_id = `user`.id").
				Where("departments.id = ? AND jia_gu_task.policy_reason = ? AND jia_gu_task.created_at BETWEEN ? AND ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], "", startTime, endTime, appTypeNum).
				Count(&count).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
					"err":  err.Error,
				})
				return
			}
			data[i]["use_recommend_policy_count"] = count

			// 未使用标准加固策略任务数
			if err := app.DB.Table("departments").
				Distinct("jia_gu_task.id").
				Joins("INNER JOIN `user` ON `user`.department_id = departments.id").
				Joins("INNER JOIN jia_gu_task ON jia_gu_task.user_id = `user`.id").
				Where("departments.id = ? AND jia_gu_task.policy_reason != ? AND jia_gu_task.created_at BETWEEN ? AND ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], "", startTime, endTime, appTypeNum).
				Count(&count).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
					"err":  err.Error,
				})
				return
			}
			data[i]["not_use_recommend_policy_count"] = count
			delete(data[i], "created_at")
			delete(data[i], "deleted_at")
			delete(data[i], "updated_at")
		}
		goto END
	}
	// 当无日期限制时 遍历添加对应数据
	for i := 0; i < len(data); i++ {
		// 统计用户
		if err := app.DB.Model(model.User{}).Where("department_id = ?", data[i]["id"]).Count(&count).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"err":  err.Error,
			})
			return
		}
		data[i]["user_count"] = count

		// 统计任务数
		if err := app.DB.Table("departments").
			Joins("INNER JOIN `user` ON `user`.department_id = departments.id").
			Joins("INNER JOIN jia_gu_task ON jia_gu_task.user_id = `user`.id").
			Where("departments.id = ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], appTypeNum).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"err":  err.Error,
			})
			return
		}
		data[i]["task_count"] = count

		// 统计使用应用数
		if err := app.DB.Table("departments").
			Distinct("jia_gu_task.app_id").
			Joins("INNER JOIN `user` ON `user`.department_id = departments.id").
			Joins("INNER JOIN jia_gu_task ON jia_gu_task.user_id = `user`.id").
			Joins("INNER JOIN application_type ON application_type.id = jia_gu_task.app_type_id").
			Where("departments.id = ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], appTypeNum).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"err":  err.Error,
			})
			return
		}
		data[i]["application_count"] = count

		// 使用推荐加固策略任务数
		if err := app.DB.Table("departments").
			Distinct("jia_gu_task.id").
			Joins("INNER JOIN `user` ON `user`.department_id = departments.id").
			Joins("INNER JOIN jia_gu_task ON jia_gu_task.user_id = `user`.id").
			Where("departments.id = ? AND jia_gu_task.policy_reason = ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], "", appTypeNum).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"err":  err.Error,
			})
			return
		}
		data[i]["use_recommend_policy_count"] = count

		// 未使用推荐加固策略任务数
		if err := app.DB.Table("departments").
			Distinct("jia_gu_task.id").
			Joins("INNER JOIN `user` ON `user`.department_id = departments.id").
			Joins("INNER JOIN jia_gu_task ON jia_gu_task.user_id = `user`.id").
			Where("departments.id = ? AND jia_gu_task.policy_reason != ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], "", appTypeNum).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"err":  err.Error,
			})
			return
		}
		data[i]["not_use_recommend_policy_count"] = count

		delete(data[i], "created_at")
		delete(data[i], "deleted_at")
		delete(data[i], "updated_at")
	}

END:
	if sortByOK && sortFieldOK {
		ArrSort(data, sortField, sortBy)
	}
	c.JSON(http.StatusOK, gin.H{
		"info": gin.H{
			"datalist": data,
		},
		"code": http.StatusOK,
		"msg":  model.ReqSuccess,
	})
}

// 以应用为主体进行数据统计
func JiaGuApplicationCount(c *gin.Context) {
	//model.JiaGuApplicationCount{} // 当前接口返回的接口体所有类型数据
	// 需要返的数据
	var data []map[string]interface{}
	// 计数值
	var count int64
	// 开始和结束日期
	var startTime time.Time
	var endTime time.Time
	var err error

	// 获取前端传入的参数
	appName, appNameOK := c.GetPostForm("application_name")
	Time1, startTimeOK := c.GetPostForm("start_time")
	Time2, endTimeOk := c.GetPostForm("end_time")
	// 排序两参数
	sortBy, sortByOK := c.GetPostForm("sort_by")
	sortField, sortFieldOK := c.GetPostForm("sort_field")
	appType, appTypeOK := c.GetPostForm("app_type")

	// 判断字符是否为空 如果为空则当该值不存在
	appNameOK = utils.IsStringEmpty(appName, appNameOK)
	startTimeOK = utils.IsStringEmpty(Time1, startTimeOK)
	endTimeOk = utils.IsStringEmpty(Time2, endTimeOk)
	sortByOK = utils.IsStringEmpty(sortBy, sortByOK)
	sortFieldOK = utils.IsStringEmpty(sortField, sortFieldOK)
	appTypeOK = utils.IsStringEmpty(appType, appTypeOK)

	if !appTypeOK {
		response.FailResult(http.StatusInternalServerError, "", model.ReqParameterMissing, c)
		return
	}

	// 加固类型 赋值对应aoolication_type的id
	var appTypeNum int
	switch appType {
	case "android":
		appTypeNum = 1
		break
	case "h5":
		appTypeNum = 2
		break
	default:
		response.FailResult(http.StatusInternalServerError, "", model.ReqParameterMissing, c)
		return
	}
	// 将string转为datetime
	if startTimeOK && endTimeOk {
		startTime, err = time.Parse("2006-01-02 15:04", Time1)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
		endTime, err = time.Parse("2006-01-02 15:04", Time2)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
	}

	// 拼接sql
	var mainSqlString bytes.Buffer
	layout := "2006-01-02 15:04:05"
	// 获取主数据部分
	mainSqlString.WriteString(" app_type_id = ")
	mainSqlString.WriteString(strconv.Itoa(appTypeNum))

	if startTimeOK && endTimeOk {
		mainSqlString.WriteString(" AND created_at BETWEEN ")
		mainSqlString.WriteString("'")
		mainSqlString.WriteString(startTime.Format(layout))
		mainSqlString.WriteString("'")
		mainSqlString.WriteString(" AND ")
		mainSqlString.WriteString("'")
		mainSqlString.WriteString(endTime.Format(layout))
		mainSqlString.WriteString("'")
	}
	if appNameOK { // 当有应用名
		mainSqlString.WriteString(" AND application.app_name like '")
		mainSqlString.WriteString("%" + appName + "%")
		mainSqlString.WriteString("'")
	}
	// 无查询条件时 {
	if err := app.DB.Model(model.Application{}).
		Where(mainSqlString.String()).
		Find(&data).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"err":  err.Error,
		})
		return
	}

	// 当有日期限制时 遍历添加对应数据
	if startTimeOK && endTimeOk {
		for i := 0; i < len(data); i++ {
			// 统计部门使用数
			if err := app.DB.Table("application").
				Distinct("departments.department_name").
				Joins("INNER JOIN jia_gu_task ON jia_gu_task.app_id = application.id").
				Joins("INNER JOIN `user` ON `user`.id = jia_gu_task.user_id").
				Joins("INNER JOIN departments ON departments.id = `user`.department_id").
				Where("application.id = ? AND jia_gu_task.created_at BETWEEN ? AND ? AND jia_gu_task.app_type_id = ?", data[i]["id"], startTime, endTime, appTypeNum).
				Count(&count).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
					"err":  err.Error,
				})
				return
			}
			data[i]["departments_count"] = count

			// 统计用户
			arr := data[i]["use_user"].(string)
			strArr := strings.Split(arr, ",")
			// 数组去重
			strArr = utils.RemoveDuplicatesAndEmpty(strArr)
			count = int64(len(strArr))
			data[i]["user_count"] = count

			// 统计任务数
			if err := app.DB.Table("application").
				Joins("INNER JOIN jia_gu_task ON jia_gu_task.app_id = application.id").
				Where("application.id = ? AND jia_gu_task.created_at BETWEEN ? AND ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], startTime, endTime, appTypeNum).
				Count(&count).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
					"err":  err.Error,
				})
				return
			}
			data[i]["task_count"] = count

			// 使用推荐加固策略任务数
			if err := app.DB.Table("application").
				Distinct("jia_gu_task.id").
				Joins("INNER JOIN jia_gu_task ON jia_gu_task.app_id = application.id").
				Where("application.id = ? AND jia_gu_task.policy_reason = ? AND jia_gu_task.created_at BETWEEN ? AND ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], "", startTime, endTime, appTypeNum).
				Count(&count).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
					"err":  err.Error,
				})
				return
			}
			data[i]["use_recommend_policy_count"] = count

			// 未使用推荐加固策略任务数
			if err := app.DB.Table("application").
				Distinct("jia_gu_task.id").
				Joins("INNER JOIN jia_gu_task ON jia_gu_task.app_id = application.id").
				Where("application.id = ? AND jia_gu_task.policy_reason != ? AND jia_gu_task.created_at BETWEEN ? AND ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], "", startTime, endTime, appTypeNum).
				Count(&count).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
					"err":  err.Error,
				})
				return
			}
			data[i]["not_use_recommend_policy_count"] = count
			// 加工应用名
			data[i]["app_name"] = data[i]["app_name"].(string) + "--" + data[i]["app_version"].(string)
			delete(data[i], "app_type_id")
			delete(data[i], "app_version")
			delete(data[i], "last_change_time")
			delete(data[i], "use_user")
			delete(data[i], "created_at")
			delete(data[i], "deleted_at")
			delete(data[i], "updated_at")
			delete(data[i], "app_cn_name")
			delete(data[i], "app_name")
			delete(data[i], "model_cn_name")
			delete(data[i], "model_name")
			delete(data[i], "the_model")
		}
		goto END
	}

	// 当无日期限制时 遍历添加对应数据
	for i := 0; i < len(data); i++ {
		// 统计部门使用数
		if err := app.DB.Table("application").
			Distinct("departments.department_name").
			Joins("INNER JOIN jia_gu_task ON jia_gu_task.app_id = application.id").
			Joins("INNER JOIN `user` ON `user`.id = jia_gu_task.user_id").
			Joins("INNER JOIN departments ON departments.id = `user`.department_id").
			Where("application.id = ? AND jia_gu_task.app_type_id = ?", data[i]["id"], appTypeNum).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"err":  err.Error,
			})
			return
		}
		data[i]["departments_count"] = count

		// 统计用户
		arr := data[i]["use_user"].(string)
		strArr := strings.Split(arr, ",")
		// 数组去重
		strArr = utils.RemoveDuplicatesAndEmpty(strArr)
		count = int64(len(strArr))
		data[i]["user_count"] = count

		// 统计任务数
		if err := app.DB.Table("application").
			Joins("INNER JOIN jia_gu_task ON jia_gu_task.app_id = application.id").
			Where("application.id = ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], appTypeNum).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"err":  err.Error,
			})
			return
		}
		data[i]["task_count"] = count

		// 使用推荐加固策略任务数
		if err := app.DB.Table("application").
			Distinct("jia_gu_task.id").
			Joins("INNER JOIN jia_gu_task ON jia_gu_task.app_id = application.id").
			Where("application.id = ? AND jia_gu_task.policy_reason = ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], "", appTypeNum).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"err":  err.Error,
			})
			return
		}
		data[i]["use_recommend_policy_count"] = count

		// 未使用推荐加固策略任务数
		if err := app.DB.Table("application").
			Distinct("jia_gu_task.id").
			Joins("INNER JOIN jia_gu_task ON jia_gu_task.app_id = application.id").
			Where("application.id = ? AND jia_gu_task.policy_reason != ? AND jia_gu_task.app_type_id = ? AND jia_gu_task.deleted_at is NULL", data[i]["id"], "", appTypeNum).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"err":  err.Error,
			})
			return
		}
		data[i]["not_use_recommend_policy_count"] = count

		// 加工应用名
		data[i]["app_name"] = data[i]["app_name"].(string) + "--" + data[i]["app_version"].(string)
		delete(data[i], "app_type_id")
		delete(data[i], "app_version")
		delete(data[i], "last_change_time")
		delete(data[i], "use_user")
		delete(data[i], "created_at")
		delete(data[i], "deleted_at")
		delete(data[i], "updated_at")
		delete(data[i], "app_cn_name")
		delete(data[i], "app_name")
		delete(data[i], "model_cn_name")
		delete(data[i], "model_name")
		delete(data[i], "the_model")
	}

END:
	if sortByOK && sortFieldOK {
		ArrSort(data, sortField, sortBy)
	}
	c.JSON(http.StatusOK, gin.H{
		"info": gin.H{
			"datalist": data,
		},
		"code": http.StatusOK,
		"msg":  model.ReqSuccess,
	})
}

// 冒泡排序 condition只能为数字的 key
func ArrSort(arr []map[string]interface{}, condition string, order string) []map[string]interface{} {
	// 记录最后一次交换位置
	lastExchangeIndex := 0
	// 无序数列边界值
	sortBorder := len(arr) - 1
	for i := 0; i < len(arr); i++ {
		isSorted := true
		for j := 0; j < sortBorder; j++ {
			temp := map[string]interface{}{}
			left := arr[j][condition].(int64)
			right := arr[j+1][condition].(int64)

			if order == "desc" { // 降序
				if left < right {
					temp = arr[j]
					arr[j] = arr[j+1]
					arr[j+1] = temp
					// 已有元素交换，所以不是有序
					isSorted = false
					// 下标更新为最后一次交换元素的位置
					lastExchangeIndex = j
				}
			} else if order == "asc" { // 升序
				if left > right {
					temp = arr[j]
					arr[j] = arr[j+1]
					arr[j+1] = temp
					// 已有元素交换，所以不是有序
					isSorted = false
					// 下标更新为最后一次交换元素的位置
					lastExchangeIndex = j
				}
			}
		}
		sortBorder = lastExchangeIndex
		if isSorted {
			break
		}
	}
	return arr
}
