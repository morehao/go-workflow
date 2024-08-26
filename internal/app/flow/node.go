package flow

import (
	"container/list"
	"errors"
	"fmt"
	"go-workflow/internal/app/object/objFlow"
	"go-workflow/internal/pkg/constants"
	"strconv"

	"github.com/morehao/go-tools/gutils"
)

// IsValidProcessConfig 检查流程配置是否有效
func IsValidProcessConfig(node *objFlow.Node) error {
	// 节点名称是否有效
	if len(node.NodeID) == 0 {
		return errors.New("节点的【nodeId】不能为空！！")
	}
	// 检查类型是否有效
	if len(node.Type) == 0 {
		return fmt.Errorf("节点【%s】的类型不能为空！！", node.NodeID)
	}
	var flag = false
	for _, nodeType := range constants.NodeTypeList {
		if nodeType == node.Type {
			flag = true
			break
		}
	}
	if !flag {
		return fmt.Errorf("节点【%s】的类型为【%s】，为无效类型, 有效类型为%s", node.NodeID, node.Type, gutils.ToJsonString(constants.NodeTypeList))
	}
	// 当前节点是否设置有审批人
	if node.Type == constants.NodeTypeApprover || node.Type == constants.NodeTypeNotifier {
		if node.Properties == nil || node.Properties.ActionerRules == nil {
			return fmt.Errorf("节点【%s】的Properties属性不能为空，如：`\"properties\": {\"actionerRules\": [{\"type\": \"target_label\",\"labelNames\": \"人事\",\"memberCount\": 1,\"actType\": \"and\"}],}`", node.NodeID)
		}
	}
	// 条件节点是否存在
	// TODO:为什么要判断条件节点长度是否为1，示例入参中最外层ConditionNodes长度为0
	if len(node.ConditionNodes) == 1 {
		return fmt.Errorf("节点【%s】条件节点下的节点数必须大于1", node.NodeID)
	}
	// 根据条件变量选择节点索引
	for _, conditionNode := range node.ConditionNodes {
		if conditionNode.Properties == nil {
			return fmt.Errorf("节点【%s】的Properties对象为空值！！", conditionNode.NodeID)
		}
		if len(conditionNode.Properties.Conditions) == 0 {
			return fmt.Errorf("节点【%s】的Conditions对象为空值！！", conditionNode.NodeID)
		}
		err := IsValidProcessConfig(conditionNode)
		if err != nil {
			return err
		}
	}

	// 子节点是否存在
	if node.ChildNode != nil {
		return IsValidProcessConfig(node.ChildNode)
	}
	return nil
}

// ParseProcessConfig 解析流程定义json数据
func ParseProcessConfig(node *objFlow.Node, variable map[string]string) (*list.List, error) {
	execNodeLinkedList := list.New()
	err := parseProcessConfig(node, variable, execNodeLinkedList)
	return execNodeLinkedList, err
}
func parseProcessConfig(node *objFlow.Node, variable map[string]string, execNodeLinkedList *list.List) (err error) {
	addNodeToExecutionList(node, execNodeLinkedList)
	// 存在条件节点
	if node.ConditionNodes != nil {
		// 如果条件节点只有一个或者条件只有一个，直接返回第一个
		if variable == nil || len(node.ConditionNodes) == 1 {
			err = parseProcessConfig(node.ConditionNodes[0].ChildNode, variable, execNodeLinkedList)
			if err != nil {
				return err
			}
		} else {
			// 根据条件变量选择节点索引
			condNode, err := GetConditionNode(node.ConditionNodes, variable)
			if err != nil {
				return err
			}
			if condNode == nil {
				return fmt.Errorf("节点【%s】找不到符合条件的子节点,检查变量【var】值是否匹配%s", node.NodeID, gutils.ToJsonString(variable))
				// panic(err)
			}
			err = parseProcessConfig(condNode, variable, execNodeLinkedList)
			if err != nil {
				return err
			}

		}
	}
	// 存在子节点
	if node.ChildNode != nil {
		err = parseProcessConfig(node.ChildNode, variable, execNodeLinkedList)
		if err != nil {
			return err
		}
	}
	return nil
}

func addNodeToExecutionList(node *objFlow.Node, execNodeLinkedList *list.List) {
	switch node.Type {
	case constants.NodeTypeApprover, constants.NodeTypeNotifier:
		var approver string
		if node.Properties.ActionerRules[0].Type == constants.ActionRuleTypeManager {
			approver = "主管"
		} else {
			approver = node.Properties.ActionerRules[0].LabelNames
		}
		execNodeLinkedList.PushBack(objFlow.ExecNode{
			NodeID:       node.NodeID,
			Type:         node.Properties.ActionerRules[0].Type,
			Approver:     approver,
			ApproverType: string(node.Type),
			MemberCount:  node.Properties.ActionerRules[0].MemberCount,
			ActType:      node.Properties.ActionerRules[0].ActType,
		})
		break
	default:
	}
}

// GetConditionNode 获取条件节点
func GetConditionNode(nodes []*objFlow.Node, paramMap map[string]string) (result *objFlow.Node, err error) {
	for _, node := range nodes {
		var flag int
		for _, v := range node.Properties.Conditions[0] {
			paramValue := paramMap[v.ParamKey]
			if len(paramValue) == 0 {
				return nil, fmt.Errorf("流程启动变量【var】的key【%s】的值不能为空", v.ParamKey)
			}
			yes, err := checkConditions(v, paramValue)
			if err != nil {
				return nil, err
			}
			if yes {
				flag++
			}
		}
		// 满足所有条件
		if flag == len(node.Properties.Conditions[0]) {
			result = node
		}
	}
	return result, nil
}

func checkConditions(cond *objFlow.NodeCondition, value string) (bool, error) {
	// 判断类型
	switch cond.Type {
	case constants.ActionConditionTypeRange:
		val, err := strconv.Atoi(value)
		if err != nil {
			return false, err
		}
		if len(cond.LowerBound) == 0 && len(cond.UpperBound) == 0 && len(cond.LowerBoundEqual) == 0 && len(cond.UpperBoundEqual) == 0 && len(cond.BoundEqual) == 0 {
			return false, fmt.Errorf("条件节点【%s】的 【upperBound】或者【lowerBound】或者【upperBoundEqual】或者【lowerBoundEqual】或者【boundEqual】不能为空，值如：'upperBound:1000'", cond.Type)
		}
		// 判断下限，lowerBound
		if len(cond.LowerBound) > 0 {
			low, err := strconv.Atoi(cond.LowerBound)
			if err != nil {
				return false, err
			}
			if val <= low {
				// fmt.Printf("val:%d小于lowerBound:%d\n", val, low)
				return false, nil
			}
		}
		if len(cond.LowerBoundEqual) > 0 {
			le, err := strconv.Atoi(cond.LowerBoundEqual)
			if err != nil {
				return false, err
			}
			if val < le {
				// fmt.Printf("val:%d小于lowerBound:%d\n", val, low)
				return false, nil
			}
		}
		// 判断上限,upperBound包含等于
		if len(cond.UpperBound) > 0 {
			upper, err := strconv.Atoi(cond.UpperBound)
			if err != nil {
				return false, err
			}
			if val >= upper {
				return false, nil
			}
		}
		if len(cond.UpperBoundEqual) > 0 {
			ge, err := strconv.Atoi(cond.UpperBoundEqual)
			if err != nil {
				return false, err
			}
			if val > ge {
				return false, nil
			}
		}
		if len(cond.BoundEqual) > 0 {
			equal, err := strconv.Atoi(cond.BoundEqual)
			if err != nil {
				return false, err
			}
			if val != equal {
				return false, nil
			}
		}
		return true, nil
	case constants.ActionConditionTypeValue:
		if len(cond.ParamValues) == 0 {
			return false, fmt.Errorf("条件节点【%s】的 【paramValues】数组不能为空，值如：'paramValues:['调休','年假']", cond.Type)
		}
		for _, val := range cond.ParamValues {
			if value == val {
				return true, nil
			}
		}
		// log.Printf("key:" + cond.ParamKey + "找不到对应的值")
		return false, nil
	default:
		return false, fmt.Errorf("未知的NodeCondition类型【%s】,正确类型应为以下中的一个:%s", cond.Type, gutils.ToJsonString(constants.ActionConditionTypeList))
	}
}
