package model

import (
	"gorm.io/gorm"
	"time"
)

type JiaGuTask struct {
	gorm.Model
	UserID       uint       `json:"user_id"`       // 用户id
	TaskID       int        `json:"task_id"`       // 任务id
	AppID        int        `json:"app_id"`        // 应用appid
	AppTypeID    int        `json:"app_type_id"`   // 应用类型id
	PolicyID     int        `json:"policy_id"`     // 应用策略id
	PolicyReason string     `json:"policy_reason"` // 使用通用加固策略原因说明
	TaskStatus   string     `json:"task_status"`   // 加固任务状态
	FinishTime   *time.Time `json:"finish_time"`   // 完成时间
}
