package model

type TemplateItem struct {
	Id         uint   `json:"id"`
	CategoryId int    `json:"category_id"`
	ItemKey    string `json:"item_key"`
	AuditName  string `json:"audit_name"`
	IsDynamic  bool   `json:"is_dynamic"`
	Status     int    `json:"status"`
	CePingType string `json:"ce_ping_type"` // ad ios mp sdk
}

// Category 模板的具体测评项表
type Category struct {
	Id           uint   `json:"id"`
	CategoryName string `json:"category_name"`
	CePingType   string `json:"ce_ping_type"` // ad ios mp sdk
}
