package constants

type NodeType string

const (
	NodeTypeStart     NodeType = "start"     // 开始，工作流的起点
	NodeTypeRoute     NodeType = "route"     // 路由，根据条件判断，决定下一步走向
	NodeTypeCondition NodeType = "condition" // 条件，用于判断，决定下一步走向
	NodeTypeApprover  NodeType = "approver"  // 审批人
	NodeTypeNotifier  NodeType = "notifier"  // 通知人
)

var NodeTypeList = []NodeType{NodeTypeStart, NodeTypeRoute, NodeTypeCondition, NodeTypeApprover, NodeTypeNotifier}

const (
	ActionRuleTypeManager = "target_management" // 目标管理
	ActionRuleTypeLabel   = "target_label"      // 目标标签
)

type ActionType string

const (
	ActionTypeAnd ActionType = "and" // 会签
	ActionTypeOr  ActionType = "or"  // 或签
)

type ActionConditionType string

const (
	ActionConditionTypeRange ActionConditionType = "dingtalk_actioner_range_condition" // 范围
	ActionConditionTypeValue ActionConditionType = "dingtalk_actioner_value_condition" // 值
)

var ActionConditionTypeList = []ActionConditionType{ActionConditionTypeRange, ActionConditionTypeValue}
