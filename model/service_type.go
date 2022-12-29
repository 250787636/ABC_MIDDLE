package model

// 加固服务类型表
type ServiceType struct {
	ID          uint   `json:"id" grom:"primarykey"` // 服务id
	ServiceType string `json:"service_type"`         // 服务类型名
}
