package model

import (
	"gorm.io/gorm"
	"time"
)

// CePingUserTask  测评用户 对应 其任务
type CePingUserTask struct {
	gorm.Model
	UserID       uint       `json:"user_id"`
	TaskID       string     `gorm:"column:task_id;type:varchar(255);NOT NULL;comment:任务ID" json:"task_id"`
	Status       string     `gorm:"column:status;type:varchar(20);NOT NULL;default:'测评中';comment:测评状态" json:"status"`
	FinishedTime *time.Time `gorm:"column:finished_at;comment:完成时间" json:"finished_time"`
	HighNum      int        `gorm:"column:high_num;type:int(11);NOT NULL;comment:高危数量" json:"high_num"`
	MiddleNum    int        `gorm:"column:middle_num;type:int(11);NOT NULL;comment:中危数量" json:"middle_num"`
	LowNum       int        `gorm:"column:low_num;type:int(11);NOT NULL;comment:低危数量" json:"low_num"`
	RiskNum      int        `gorm:"column:risk_num;type:int(11);NOT NULL;comment:风险数量" json:"risk_num"`
	PkgName      string     `gorm:"column:pkg_name;type:varchar(255);NOT NULL;comment:应用名" json:"pkg_name"`
	AppName      string     `gorm:"column:app_name;type:varchar(255);NOT NULL;comment:文件名" json:"app_name"` // 文件名称
	Version      string     `gorm:"column:version;type:varchar(255);NOT NULL;comment:版本" json:"version"`
	TemplateID   uint       `gorm:"column:template_id;type:int(11);comment:模版ID" json:"template_id"`
	TaskType     int        `gorm:"column:task_type;type:int(11);NOT NULL;default:1;comment:任务类型" json:"task_type"` // 1 android 2 ios 3 小程序
	FilePath     string     `gorm:"column:file_path;type:varchar(255);NOT NULL;default:'';comment:文件路径" json:"file_path"`
	Score        int        `gorm:"column:score;type:int(11);NOT NULL;default:0;comment:总分" json:"score"`
	ItemsNum     int        `gorm:"column:items_num;type:int(11);NOT NULL;default:0;comment:总项数" json:"items_num"`
	ViewUrl      string     `gorm:"column:view_url;type:varchar(255);NOT NULL;default:'';comment:查看地址" json:"view_url"`
	FinishItem   int        `gorm:"column:finish_item;type:int(11);NOT NULL;default:0;comment:完成项数" json:"finish_item"`
}

type SdkUse struct {
	TaskID  string `json:"task_id"`  // 任务id
	TaskTag string `json:"task_tag"` // 分类标签
}
