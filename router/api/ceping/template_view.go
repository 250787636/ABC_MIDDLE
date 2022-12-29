package ceping

import (
	"example.com/m/model"
	"gorm.io/gorm"
)

type TemplateItemView struct {
	ItemKey   string `json:"item_key"`
	AuditName string `json:"audit_name"`
	IsDynamic bool   `json:"is_dynamic"`
	Status    int    `json:"status"`
}

type TemplateView struct {
	TemplateId       int                           `json:"template_id"`
	TemplateName     string                        `json:"template_name"`
	CreatorAccount   string                        `json:"creator_account"`
	CreateTime       string                        `json:"create_time"`
	IsOWASP          bool                          `json:"is_owasp"`
	ReportLanguage   string                        `json:"report_language"`
	Categories       []string                      `json:"categories"`
	CategorizedItems map[string][]TemplateItemView `json:"categorized_items"`
}

func TemplateItemKeysToView(db *gorm.DB, typ string, itemKeys []string) (view TemplateView, err error) {
	var auditCategories []model.CepingAuditCategory
	err = db.Where("type = ?", typ).Order("category_sort ASC").Find(&auditCategories).Error
	if err != nil {
		return view, err
	}
	categoryIdMap := make(map[int]string, len(auditCategories))
	for _, category := range auditCategories {
		view.Categories = append(view.Categories, category.CategoryName)
		categoryIdMap[category.Id] = category.CategoryName
	}

	itemKeySet := StringArrayToSet(itemKeys)
	view.CategorizedItems = make(map[string][]TemplateItemView, len(itemKeys))
	addItemView := func(itemKey string, name string, categoryId int) {
		var value int
		if _, ok := itemKeySet[itemKey]; ok {
			value = 1
		}
		if itemKeys == nil {
			value = 1
		}
		category := categoryIdMap[categoryId]
		view.CategorizedItems[category] = append(view.CategorizedItems[category], TemplateItemView{
			ItemKey:   itemKey,
			AuditName: name,
			Status:    value,
		})
	}
	switch typ {
	case "ad":
		var auditItems []model.CepingAdAuditItem
		if err = db.Where("status = 1").Find(&auditItems).Error; err != nil {
			return view, err
		}
		for _, auditItem := range auditItems {
			itemKey := auditItem.ItemKey
			if itemKey == "sec_infos" {
				continue
			}
			var value int
			if _, ok := itemKeySet[itemKey]; ok {
				value = 1
			}
			if itemKeys == nil {
				value = 1
			}
			category := categoryIdMap[auditItem.CategoryId]
			view.CategorizedItems[category] = append(view.CategorizedItems[category], TemplateItemView{
				ItemKey:   itemKey,
				AuditName: auditItem.Name,
				IsDynamic: auditItem.IsDynamic == 1,
				Status:    value,
			})
		}
	case "ios":
		var auditItems []model.CepingIosAuditItem
		if err = db.Where("status = 1").Find(&auditItems).Error; err != nil {
			return view, err
		}
		for _, auditItem := range auditItems {
			itemKey := auditItem.ItemKey
			if itemKey == "ios_sec_infos" {
				continue
			}
			addItemView(itemKey, auditItem.Name, auditItem.CategoryId)
		}
	case "sdk":
		var auditItems []model.CepingSdkAuditItem
		if err = db.Where("status = 1").Find(&auditItems).Error; err != nil {
			return view, err
		}
		for _, auditItem := range auditItems {
			itemKey := auditItem.ItemKey
			if itemKey == "sdk_sec_infos" {
				continue
			}
			addItemView(itemKey, auditItem.Name, auditItem.CategoryId)
		}
	case "mp":
		var auditItems []model.CepingMpAuditItem
		if err = db.Where("status = 1").Find(&auditItems).Error; err != nil {
			return view, err
		}
		for _, auditItem := range auditItems {
			itemKey := auditItem.ItemKey
			if itemKey == "mp_sec_infos" {
				continue
			}
			addItemView(itemKey, auditItem.Name, auditItem.CategoryId)
		}
	}
	return view, nil
}
func StringArrayToSet(arr []string) map[string]struct{} {
	result := make(map[string]struct{}, len(arr))
	for _, v := range arr {
		result[v] = struct{}{}
	}
	return result
}
