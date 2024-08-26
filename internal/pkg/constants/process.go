package constants

const (
	NodeIDTextStart = "开始"
	NodeIDTextEnd   = "结束"
)

const (
	ProcessTaskStatusUnfinished = 1 // 任务未完成
	ProcessTaskStatusFinished   = 2 // 任务完成
)

const (
	NodeInfoTypeStarter = "starter" // 开始节点
)

type IdentityType string

const (
	IdentityTypeCandidate   IdentityType = "candidate"   // 候选人
	IdentityTypeParticipant IdentityType = "participant" // 参与者
	IdentityTypeManager     IdentityType = "manager"     // 主管
	IdentityTypeNotifier    IdentityType = "notifier"    // 通知人
)

var IdentityTypeList = []IdentityType{IdentityTypeCandidate, IdentityTypeParticipant, IdentityTypeManager, IdentityTypeNotifier}

const (
	TaskOperateTypeApprove = 1 // 批准
	TaskOperateTypeReject  = 2 // 拒绝
)
