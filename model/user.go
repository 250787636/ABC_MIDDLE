package model

import (
	"database/sql/driver"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type User struct {
	gorm.Model
	UserName      string `json:"name"`            // 用户名
	DepartmentID  uint   `json:"department_id"`   // 部门id
	AccountLevel  string `json:"account_level"`   // 账户等级
	JobTitle      string `json:"job_title"`       // 职位名称
	LastLoginTime string `json:"last_login_time"` // 最近一次登录

	Token       string `json:"token"`        // 生成验证token
	ClientToken string `json:"client_token"` // 客户传入的token
	IsAdmin     bool   `json:"is_admin"`     // 是否是管理员
}

type FormatTime time.Time

func (t FormatTime) MarshalJSON() ([]byte, error) {
	var timeStr string
	if !time.Time(t).IsZero() {
		timeStr = time.Time(t).Format("2006-01-02 15:04:05")
	}
	return []byte(strconv.Quote(timeStr)), nil
}
func (t FormatTime) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}
	return time.Time(t), nil
}
func (t *FormatTime) UnmarshalJSON(stringTime []byte) error {
	if t == nil {
		return nil
	}
	v, _ := strconv.Unquote(string(stringTime))
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
	*t = FormatTime(tt)
	return nil
}
