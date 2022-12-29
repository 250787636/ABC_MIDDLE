package ceping

import (
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/utils/response"
	"github.com/gin-gonic/gin"
	"math"
	"time"
)

type DetailsRequest struct {
	DepartmentName string `form:"department_name"`
	//KeyWord        string `form:"key_word"`
	StartTime  string `form:"start_time"`
	EndTime    string `form:"end_time"`
	PageSize   int    `form:"size" binding:"required"`
	PageNumber int    `form:"page" binding:"required"`
}

type Info struct {
	DepartmentName string `json:"department_name"`
	UserTotal      int64  `json:"user_total"`
	TaskTotal      int64  `json:"task_total"`
	HighNum        int    `json:"high_num"`
	MiddleNum      int    `json:"middle_num"`
	LowNum         int    `json:"low_num"`
	RiskNum        int    `json:"risk_num"`
}

// GetDetails
func GetDetails(c *gin.Context) {
	req := DetailsRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	//dataInfo := make([]map[string]interface{}, 0)

	//info := make([]Info, 0)
	sql := app.DB.Debug().Model(&model.Departments{})

	if req.StartTime != "" {
		//time, err := StringToTimeYMD(req.StartTime)
		//if err != nil {
		//	fmt.Println(err)
		//}
		sql = sql.Where("departments.created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		//time, err := StringToTimeYMD(req.StartTime)
		//if err != nil {
		//	fmt.Println(err)
		//}
		sql = sql.Where("departments.created_at <= ?", req.EndTime)
	}
	if req.DepartmentName != "" {
		sql = sql.Where("departments.department_name like ?", "%"+req.DepartmentName+"%")
	}
	//if req.KeyWord != "" {
	//	sql = sql.Where("departments.department_name like ?", "%"+req.KeyWord+"%")
	//}

	// 1.获取部门id
	firstData := make([]map[string]interface{}, 0)
	sql.Find(&firstData)

	//reponseData := make([]map[string]interface{}, 0)

	var total int64
	for _, data := range firstData {
		userMap := make([]map[string]interface{}, 0)
		sql2 := app.DB.Model(&model.User{}).Where("department_id = ? ", data["id"])
		if req.StartTime != "" {
			sql2 = sql2.Where("created_at >= ?", req.StartTime)
		}
		if req.EndTime != "" {
			sql2 = sql2.Where("created_at <= ?", req.EndTime)
		}
		sql2.Count(&total)
		sql2.Find(&userMap)
		data["user_total"] = total

		// 3.获取任务
		sql3 := app.DB.Debug().Table("ce_ping_user_task").
			Joins("left join user on user.id = ce_ping_user_task.user_id").
			Joins("left join departments on departments.id = user.department_id").
			Where("departments.id = ?", data["id"]).
			//Where("ce_ping_user_task.task_type = 1"). // 这里只统计该部门下所有用户的安卓任务数
			Where("ce_ping_user_task.deleted_at IS NULL")
		if req.StartTime != "" {
			sql3 = sql3.Where("ce_ping_user_task.created_at >= ?", req.StartTime)
		}
		if req.EndTime != "" {
			sql3 = sql3.Where("ce_ping_user_task.created_at <= ?", req.EndTime)
		}
		sql3.Count(&total)
		data["task_total"] = total

		high := 0
		middle := 0
		low := 0
		risk := 0
		for _, value := range userMap {
			findMap := make(map[string]interface{}, 0)
			sql4 := app.DB.Model(&model.CePingUserTask{})
			sql4.Select("sum(ce_ping_user_task.high_num) as high_num,"+
				"sum(ce_ping_user_task.middle_num) as middle_num,"+
				"sum(ce_ping_user_task.low_num) as low_num").
				Joins("left join user on ce_ping_user_task.user_id = user.id").
				Where("ce_ping_user_task.user_id = ?", value["id"]).
				//Where("ce_ping_user_task.task_type = 1"). // 这里只统计该部门下所有用户的安卓风险数
				Find(&findMap)

			if _, ok := findMap["high_num"].(int); !ok {
				findMap["high_num"] = 0
			}

			if _, ok := findMap["middle_num"].(int); !ok {
				findMap["middle_num"] = 0
			}

			if _, ok := findMap["low_num"].(int); !ok {
				findMap["low_num"] = 0
			}

			high += findMap["high_num"].(int)
			middle += findMap["middle_num"].(int)
			low += findMap["low_num"].(int)
			risk = high + middle + low

		}

		data["high_num"] = high
		data["middle_num"] = middle
		data["low_num"] = low
		data["risk_num"] = risk

		delete(data, "created_at")
		delete(data, "deleted_at")
		delete(data, "updated_at")
		delete(data, "id")
	}

	total = int64(len(firstData))
	offNum := req.PageSize * (req.PageNumber - 1)

	max := req.PageSize * req.PageNumber
	if max > len(firstData) {
		max = len(firstData)
	}
	finData := make([]map[string]interface{}, 0)
	for i := offNum; i < max; i++ {
		finData = append(finData, firstData[i])
	}

	response.OkWithList(finData, int(total), req.PageNumber, req.PageSize, c)

}

func StringToTimeYMD(in string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04", in, time.Local)
}

func CutPage(total int, page int, limit int) (pager PageList) {
	pager.Count = total
	pager.TotalPage = int(math.Ceil(float64(total) / float64(limit)))
	pager.Page = page
	pager.Limit = limit

	return
}

type PageList struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Count     int `json:"count"`
	TotalPage int `json:"total_page"`
}
