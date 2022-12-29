package model

// 以部门开头统计数据
type JiaGuDepartmentCount struct {
	ID                         uint   `json:"id" grom:"primarykey"`
	DepartmentName             string `json:"department_name"`
	ApplicationCount           int64  `json:"application_count"`
	UserCount                  int64  `json:"user_count"`
	TaskCount                  int64  `json:"task_count"`
	UseRecommendPolicyCount    int64  `json:"use_recommend_policy_count"`
	NotUseRecommendPolicyCount int64  `json:"not_use_recommend_policy_count"`
	Reason                     int64  `json:"reason"`
}
