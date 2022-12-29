package service

import (
	"example.com/m/model"
	"example.com/m/pkg/app"
	"gorm.io/gorm"
	"math"
	"time"
)

type TemplateService struct {
	db *gorm.DB
}

var Template = &TemplateService{db: app.DB}

type QueryTemplateRequest struct {
	PageSize     int
	PageNumber   int
	TemplateType string //  安卓ad 苹果 ios 小程序 mp  SDK任务 sdk
	IsPage       bool   // 是否需要分页 1 是  2 不是
	IsSuper      bool   //是否为超管
	IsAdmin      bool   //是否为管理员
	DepartId     int    //部门id
	UserId       int    //用户id
}

// 模板返回的参数
type TemplateList struct {
	CreatedAt    model.FormatTime `json:"created_at"`
	CreatedId    int              `json:"created_id"`
	Id           int              `json:"id"`
	TemplateName string           `json:"template_name"`
}

func (t *TemplateService) ListTemplate(request QueryTemplateRequest) (total int64, tempList []TemplateList, err error) {
	tx := t.db.Model(&model.Template{}).Where("template_type = ?", request.TemplateType)
	if !request.IsSuper {
		if request.IsAdmin {
			tx = tx.Joins("inner join user on user.id = template.created_id").
				Where("user.department_id = ? ", request.DepartId)
		} else {
			tx = tx.Where("created_id = ? ", request.UserId)
		}
	}
	if !request.IsPage {
		// 获取数据
		if err := tx.Count(&total).Find(&tempList).Error; err != nil {
			return 0, nil, err
		}
	} else {
		// 获取条数
		if err := tx.Count(&total).Error; err != nil {
			return 0, nil, err
		}
		// 获取数据
		if err := tx.Scopes(Paginate(request.PageNumber, request.PageSize)).Find(&tempList).Error; err != nil {
			return 0, nil, err
		}
	}

	templateName := ""
	switch request.TemplateType {
	case "ad":
		templateName = "Android-全量模板"
	case "ios":
		templateName = "IOS-全量模板"
	case "sdk":
		templateName = "SDK-全量模板"
	}

	//"2022-10-11 16:08:07"
	total++
	lastPage := int(math.Ceil(float64(total) / float64(request.PageSize)))
	if !request.IsPage || request.PageNumber == lastPage { //拼接系统全量模板，全量模板id为0，不查数据库
		tempList = append(tempList, TemplateList{
			CreatedAt:    model.FormatTime(time.Date(2022, 10, 1, 0, 0, 0, 0, time.Local)),
			TemplateName: templateName,
		})
	}

	return total, tempList, nil
}

func Paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNum <= 0 {
			pageNum = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
