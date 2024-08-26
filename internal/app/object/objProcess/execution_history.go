package objProcess

type ExecutionHistoryBaseInfo struct {
	Rev         int32  `json:"rev" form:"rev"`                 // 修订版本号
	ProcInstID  uint64 `json:"procInstId" form:"procInstId"`   // 流程实例ID
	ProcDefID   uint64 `json:"procDefId" form:"procDefId"`     // 流程定义ID
	ProcDefName string `json:"procDefName" form:"procDefName"` // 流程定义名称
	NodeInfos   string `json:"nodeInfos" form:"nodeInfos"`     // 节点信息
	IsActive    int8   `json:"isActive" form:"isActive"`       // 是否活跃
	StartTime   string `json:"startTime" form:"startTime"`     // 开始时间
}
