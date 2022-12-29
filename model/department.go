package model

import "gorm.io/gorm"

type Departments struct {
	gorm.Model
	DepartmentName string `json:"department_name"` // 部门名称
}
