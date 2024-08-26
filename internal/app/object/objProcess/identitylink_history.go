package objProcess

type IdentitylinkHistoryBaseInfo struct {
	Comment    string `json:"comment" form:"comment"`       // 评论
	Company    string `json:"company" form:"company"`       // 公司
	Group      string `json:"group" form:"group"`           // 组
	ProcInstID uint64 `json:"procInstId" form:"procInstId"` // 流程实例ID
	Step       int32  `json:"step" form:"step"`             // 步骤
	TaskID     uint64 `json:"taskId" form:"taskId"`         // 任务ID
	Type       string `json:"type" form:"type"`             // 类型
	UserID     string `json:"userId" form:"userId"`         // 用户ID
	UserName   string `json:"userName" form:"userName"`     // 用户名称
}
