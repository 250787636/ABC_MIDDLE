package model

type ApplicationType struct {
	ID      uint   `json:"id" grom:"primarykey"` // 应用id
	AppType string `json:"app_type"`             // 应用类型名
}
