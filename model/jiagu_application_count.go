package model

// 以应用开头统计数据
type JiaGuApplicationCount struct {
	ID                         uint   `json:"id" grom:"primarykey"`
	ApplicationName            string `json:"application_name"`
	DepartmentID               int    `json:"department_id"`
	UserCount                  int    `json:"user_count"`
	TaskCount                  int    `json:"task_count"`
	ApplicationCount           int    `json:"application_count"`
	UseRecommendPolicyCount    int    `json:"use_recommend_policy_count"`
	NotUseRecommendPolicyCount int    `json:"not_use_recommend_policy_count"`
	Reason                     int    `json:"reason"`
}
