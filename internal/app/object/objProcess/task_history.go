package objProcess

type TaskHistoryBaseInfo struct {
	NodeID        string `json:"nodeId" form:"nodeId"`               // 节点ID
	Step          int32  `json:"step" form:"step"`                   // 步骤
	ProcInstID    uint64 `json:"procInstId" form:"procInstId"`       // 流程实例ID
	Assignee      string `json:"assignee" form:"assignee"`           // 受派人
	CreateTime    string `json:"createTime" form:"createTime"`       // 创建时间
	ClaimTime     string `json:"claimTime" form:"claimTime"`         // 认领时间
	MemberCount   int8   `json:"memberCount" form:"memberCount"`     // 成员数量
	UnCompleteNum int8   `json:"unCompleteNum" form:"unCompleteNum"` // 未完成数量
	AgreeNum      int8   `json:"agreeNum" form:"agreeNum"`           // 同意数量
	ActType       string `json:"actType" form:"actType"`             // 行为类型
	IsFinished    int8   `json:"isFinished" form:"isFinished"`       // 是否完成
}
