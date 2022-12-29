package model

import "gorm.io/gorm"

// 加固手册
type JiaguOperationManual struct {
	gorm.Model
	ServiceId int    `json:"service_id"` // 服务名id
	FileName  string `json:"file_name"`  // 文件名称
	FilePath  string `json:"file_path"`  // 文件下载地址
}
