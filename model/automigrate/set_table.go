package automigrate

import "example.com/m/model"

func setTable(migrate *AutoMigrate) {
	// 数据表多用于 c u d
	migrate.DataTables = []TableVersion{
		{model.Departments{}, ""},
		{model.JiaGuTask{}, ""},
		{model.CePingUserTask{}, ""},
		{model.Application{}, ""},
		{model.JiaguOperationManual{}, ""},
		{model.JiaguPolicyAndroid{}, ""},
		{model.JiaguPolicyH5{}, ""},
		{model.SdkUse{}, ""},
		// 模板分类表
		{model.Template{}, ""},
	}
	//工具表多用于 r
	migrate.ToolTables = []TableVersion{
		{model.ApplicationType{}, ""},
		{model.ServiceType{}, ""},
		// 需内置超级管理员
		{model.User{}, ""},
		{model.Category{}, ""},
		{model.TemplateItem{}, ""},
		{model.CepingAdAuditItem{}, ""},
		{model.CepingIosAuditItem{}, ""},
		{model.CepingAuditCategory{}, ""},
		{model.CepingSdkAuditItem{}, ""},
	}
}
