package objProcDef

import "go-workflow/internal/app/object/objFlow"

type ProcDefBaseInfo struct {
	Name     string        `json:"name" form:"name" validate:"required"`         // 流程名称
	Resource *objFlow.Node `json:"resource" form:"resource" validate:"required"` // 流程配置
}
