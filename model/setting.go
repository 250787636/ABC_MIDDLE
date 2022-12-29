package model

import "gorm.io/gorm"

type GlobalVariable struct {
	gorm.Model
	Type  string `gorm:"type:varchar(255)"`
	Name  string `gorm:"type:varchar(255)"`
	Value string `gorm:"type:varchar(255)"`
}
