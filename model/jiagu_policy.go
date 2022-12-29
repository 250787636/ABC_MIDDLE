package model

// 加固7.5.2的策略结构体
// 安卓加固策略表
type JiaguPolicyAndroid struct {
	Id         int    `json:"id"`          //策略id
	Status     int    `json:"status"`      // 策略状态
	Name       string `json:"name"`        // 策略名称
	ConfigJson string `json:"config_json"` // 策略配置
	Type       string `json:"type"`        // 加固策略类型
	CustId     int    `json:"cust_id"`     // 客户id
}

// 加固7.5.2的策略结构体
// H5加固策略表
type JiaguPolicyH5 struct {
	Id         int    `json:"id"`          //策略id
	Status     int    `json:"status"`      // 策略状态
	Name       string `json:"name"`        // 策略名称
	ConfigJson string `json:"config_json"` // 策略配置
	Type       string `json:"type"`        // 加固策略类型
	CustId     int    `json:"cust_id"`     // 客户id
}
