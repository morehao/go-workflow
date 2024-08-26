package objFlow

import (
	"database/sql/driver"
	"go-workflow/internal/pkg/constants"

	jsoniter "github.com/json-iterator/go"
)

// Node 表示工作流中的一个节点。
type Node struct {
	Name           string             `json:"name"`           // 节点名称
	Type           constants.NodeType `json:"type"`           // 节点类型
	NodeID         string             `json:"nodeId"`         // 节点ID
	PrevID         string             `json:"prevId"`         // 前一个节点ID
	ChildNode      *Node              `json:"childNode"`      // 子节点
	ConditionNodes []*Node            `json:"conditionNodes"` // 条件节点列表
	Properties     *NodeProperties    `json:"properties"`     // 节点属性
}

func (node Node) Value() (driver.Value, error) {
	res, err := jsoniter.Marshal(node)
	return string(res), err
}

func (node *Node) Scan(input interface{}) error {
	bytes := input.([]byte)
	if len(bytes) == 0 {
		return nil
	}
	return jsoniter.Unmarshal(input.([]byte), node)
}

// NodeProperties 定义了节点的特定属性。
type NodeProperties struct {
	ActivateType       string             `json:"activateType"`       // 激活类型，例如"ONE_BY_ONE"代表依次审批
	AgreeAll           bool               `json:"agreeAll"`           // 是否需要全部同意
	Conditions         [][]*NodeCondition `json:"conditions"`         // 条件列表
	ActionerRules      []*ActionerRule    `json:"actionerRules"`      // 行动者规则列表
	NoneActionerAction string             `json:"noneActionerAction"` // 无行动者时的操作
}

// NodeCondition 定义了节点的条件。
type NodeCondition struct {
	Type            constants.ActionConditionType `json:"type"`            // 条件类型
	ParamKey        string                        `json:"paramKey"`        // 参数键
	ParamLabel      string                        `json:"paramLabel"`      // 参数标签
	IsEmpty         bool                          `json:"isEmpty"`         // 参数是否为空
	LowerBound      string                        `json:"lowerBound"`      // 下界
	LowerBoundEqual string                        `json:"lowerBoundEqual"` // 下界是否包含等于
	UpperBoundEqual string                        `json:"upperBoundEqual"` // 上界是否包含等于
	UpperBound      string                        `json:"upperBound"`      // 上界
	BoundEqual      string                        `json:"boundEqual"`      // 界限是否相等
	Unit            string                        `json:"unit"`            // 单位
	ParamValues     []string                      `json:"paramValues"`     // 参数值列表
	OriValue        []string                      `json:"oriValue"`        // 原始值列表
	Conds           []*NodeCond                   `json:"conds"`           // 更细分的条件列表
}

// NodeCond 表示 NodeCondition 中的单个条件。
type NodeCond struct {
	Type  string    `json:"type"`  // 条件类型
	Value string    `json:"value"` // 条件值
	Attrs *NodeUser `json:"attrs"` // 用户属性
}

// NodeUser 表示节点条件上下文中的用户。
type NodeUser struct {
	Name   string `json:"name"`   // 用户名
	Avatar string `json:"avatar"` // 用户头像
}

// ActionerRule 定义了节点中行动者的规则。
type ActionerRule struct {
	Type        string               `json:"type"`        // 规则类型
	LabelNames  string               `json:"labelNames"`  // 标签名称
	Labels      int                  `json:"labels"`      // 标签
	IsEmpty     bool                 `json:"isEmpty"`     // 是否为空
	MemberCount int64                `json:"memberCount"` // 需要通过的人数，如果是会签
	ActType     constants.ActionType `json:"actType"`     // 行动类型，"and" 表示会签，"or" 表示或签，默认为或签
	Level       int8                 `json:"level"`       // 级别
	AutoUp      bool                 `json:"autoUp"`      // 是否自动升级
}

type ExecNode struct {
	NodeID       string               `json:"nodeId"`       // 节点ID
	Type         string               `json:"type"`         // 节点类型
	Approver     string               `json:"approver"`     // 审批人
	ApproverType string               `json:"approverType"` // 审批人类型
	MemberCount  int64                `json:"memberCount"`  // 需要通过的人数，如果是会签
	Level        int64                `json:"level"`        // 级别
	ActType      constants.ActionType `json:"actType"`      // 行动类型，"and" 表示会签，"or" 表示或签，默认为或签
}

type ExecNodeList []ExecNode

func (node ExecNodeList) Value() (driver.Value, error) {
	res, err := jsoniter.Marshal(node)
	return string(res), err
}

func (node *ExecNodeList) Scan(input interface{}) error {
	bytes := input.([]byte)
	if len(bytes) == 0 {
		return nil
	}
	return jsoniter.Unmarshal(input.([]byte), node)
}
