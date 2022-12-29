package model

import "gorm.io/gorm"

type Application struct {
	gorm.Model
	AppName     string `json:"app_name"`      // 系统编号
	ModelName   string `json:"model_name"`    // 模板编号
	AppCnName   string `json:"app_cn_name"`   // 系统中文全称
	AppVersion  string `json:"app_version"`   // 系统英文简称
	ModelCnName string `json:"model_cn_name"` // 模板中文全称

	RecommendPolicy int    `json:"recommend_policy"` // 推荐策略
	AppTypeID       int    `json:"app_type_id"`      // 系统类型编号
	UseUser         string `json:"use_user"`         // 使用的用户
	LastChangeTime  string `json:"last_change_time"` // 录入时间

	TheApp   string `json:"the_app"`   // 系统编号 + 系统中文全称
	TheModel string `json:"the_model"` //  模板编号 + 模板中文全称
}
