package dtoProcess

import (
	"go-workflow/internal/app/object/objCommon"
	"go-workflow/internal/app/object/objProcess"
)

type ProcInstStartReq struct {
	ProcDefName string            `json:"procDefName" validate:"required"` // 流程定义名称
	Title       string            `json:"title" validate:"required"`       // 流程标题
	Var         map[string]string `json:"var"`                             // 流程变量
}

type ProcInstUpdateReq struct {
	ID uint64 `json:"id" validate:"required" label:"数据自增id"` // 数据自增id
	objProcess.ProcInstBaseInfo
}

type ProcInstDetailReq struct {
	ID uint64 `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}

type ProcInstPageListReq struct {
	objCommon.PageQuery
}

type ProcInstDeleteReq struct {
	ID uint64 `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}

type CreatedPageListReq struct {
	objCommon.PageQuery
}

type TodoPageListReq struct {
	objCommon.PageQuery
	GroupList      []string `json:"groupList" form:"groupList"`           // 用户组列表
	DepartmentList []string `json:"departmentList" form:"departmentList"` // 部门列表
}

type NotifyPageListReq struct {
	objCommon.PageQuery
	GroupList []string `json:"groupList" form:"groupList"` // 用户组列表
}

type TaskCompleteReq struct {
	TaskID      uint64 `json:"taskId" validate:"required" label:"任务id"`      // 任务id
	OperateType int8   `json:"operateType" validate:"required" label:"操作类型"` // 操作类型，1：同意，2：拒绝
	Comment     string `json:"comment" label:"审批意见"`                         // 审批意见
	Candidate   string `json:"candidate" label:"候选人"`                        // 候选人
}
type TaskWithDrawReq struct {
	TaskID     uint64 `json:"taskId" validate:"required" label:"任务id"`       // 任务id
	ProcInstID uint64 `json:"procInstId" validate:"required" label:"流程实例id"` // 流程实例id
	Comment    string `json:"comment" label:"审批意见"`                          // 审批意见
}
