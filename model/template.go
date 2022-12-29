package model

import "gorm.io/gorm"

type Template struct {
	gorm.Model
	TemplateName   string `gorm:"column:template_name;type:varchar(100);NOT NULL;default:'';comment:任务ID" json:"template_name"`
	CreatedID      int    `gorm:"column:created_id;type:int(11);comment:创建者ID" json:"created_id"`
	TemplateType   string `gorm:"column:template_type;type:varchar(100);NOT NULL;default:'';comment:模板类型" json:"template_type"`
	Items          string `gorm:"column:items;type:longtext;NOT NULL;comment:模板内容" json:"items"`
	IsOwasp        int    `gorm:"column:is_owasp;type:tinyint(1);NOT NULL;default:0;comment:是否OWASP" json:"is_owasp"`
	ReportLanguage string `gorm:"column:report_language;type:varchar(100);NOT NULL;default:'';comment:报告语言" json:"report_language"` // "zh_cn": "中文简体", "zh_tw": "中文繁体", "ja_jp": "日文", "en_us": "英文"
}
