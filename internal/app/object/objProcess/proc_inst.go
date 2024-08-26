package objProcess

type ProcInstBaseInfo struct {
	Candidate     string `json:"candidate" form:"candidate"`         // 候选人
	Company       string `json:"company" form:"company"`             // 公司名称
	Department    string `json:"department" form:"department"`       // 部门名称
	Duration      uint64 `json:"duration" form:"duration"`           // 持续时间
	EndTime       string `json:"endTime" form:"endTime"`             // 结束时间
	IsFinished    int8   `json:"isFinished" form:"isFinished"`       // 是否完成
	NodeID        string `json:"nodeId" form:"nodeId"`               // 当前节点ID
	ProcDefID     uint64 `json:"procDefId" form:"procDefId"`         // 流程定义ID
	ProcDefName   string `json:"procDefName" form:"procDefName"`     // 流程定义名称
	StartTime     string `json:"startTime" form:"startTime"`         // 开始时间
	StartUserID   string `json:"startUserId" form:"startUserId"`     // 发起用户ID
	StartUserName string `json:"startUserName" form:"startUserName"` // 发起用户名称
	TaskID        uint64 `json:"taskId" form:"taskId"`               // 任务ID
	Title         string `json:"title" form:"title"`                 // 流程实例标题
}
