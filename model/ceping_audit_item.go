package model

type CepingAuditCategory struct {
	Type         string `gorm:"type:varchar(20) NOT NULL;primarykey"`
	Id           int    `gorm:"type:int(11) NOT NULL;primarykey"`
	CategoryKey  string `gorm:"type:varchar(100)"`
	CategoryName string `gorm:"type:varchar(100)"`
	CategorySort int    `gorm:"type:int(11)"`
}

type CepingAdAuditItem struct {
	Id         int    `gorm:"primarykey"`
	CategoryId int    `gorm:"type:int(11)"`
	ItemKey    string `gorm:"type:varchar(255)"`
	Name       string `gorm:"type:varchar(255)"`
	Level      string `gorm:"type:varchar(10)"`
	Score      int    `gorm:"type:int(11)"`
	IsDynamic  int    `gorm:"type:tinyint(4)"`
	Sort       int    `gorm:"type:int(11)"`
	Status     int    `gorm:"tinyint(4)"`
	Solution   string `gorm:"type:longtext"`
}

type CepingIosAuditItem struct {
	Id         int    `gorm:"primarykey"`
	CategoryId int    `gorm:"type:int(11)"`
	ItemKey    string `gorm:"type:varchar(255)"`
	Name       string `gorm:"type:varchar(255)"`
	Level      string `gorm:"type:varchar(10)"`
	Score      int    `gorm:"type:int(11)"`
	Sort       int    `gorm:"type:int(11)"`
	Status     int    `gorm:"type:tinyint(4)"`
	Solution   string `gorm:"type:longtext"`
}

type CepingSdkAuditItem struct {
	Id         int    `gorm:"primarykey"`
	CategoryId int    `gorm:"type:int(11)"`
	ItemKey    string `gorm:"type:varchar(255)"`
	Name       string `gorm:"type:varchar(255)"`
	Level      string `gorm:"type:varchar(10)"`
	Sort       int    `gorm:"type:int(11)"`
	Status     int    `gorm:"type:tinyint(4)"`
	Score      int    `gorm:"type:int(11)"`
	Solution   string `gorm:"type:longtext"`
}

type CepingMpAuditItem struct {
	Id         int    `gorm:"primarykey"`
	CategoryId int    `gorm:"type:int(11)"`
	ItemKey    string `gorm:"type:varchar(255)"`
	Name       string `gorm:"type:varchar(255)"`
	Level      string `gorm:"type:varchar(10)"`
	Solution   string `gorm:"type:longtext"`
	Sort       int    `gorm:"type:int(11)"`
	Status     int    `gorm:"type:tinyint(4)"`
}
