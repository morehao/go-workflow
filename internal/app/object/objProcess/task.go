package objProcess

type TaskBaseInfo struct {
	ActType       string `json:"actType" form:"actType"`             // 行为类型
	AgreeNum      int8   `json:"agreeNum" form:"agreeNum"`           // 同意数量
	Assignee      string `json:"assignee" form:"assignee"`           // 受派人
	ClaimTime     string `json:"claimTime" form:"claimTime"`         // 认领时间
	CreateTime    string `json:"createTime" form:"createTime"`       // 创建时间
	IsFinished    int8   `json:"isFinished" form:"isFinished"`       // 是否完成
	MemberCount   int8   `json:"memberCount" form:"memberCount"`     // 成员数量
	NodeID        string `json:"nodeId" form:"nodeId"`               // 节点ID
	ProcInstID    uint64 `json:"procInstId" form:"procInstId"`       // 流程实例ID
	Step          uint64 `json:"step" form:"step"`                   // 步骤
	UnCompleteNum int8   `json:"unCompleteNum" form:"unCompleteNum"` // 未完成数量
}
