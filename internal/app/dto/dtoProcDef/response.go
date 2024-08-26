package dtoProcDef

import (
	"go-workflow/internal/app/object/objCommon"
	"go-workflow/internal/app/object/objProcDef"
)

type ProcDefSaveResp struct {
	ID uint64 `json:"id"` // 数据自增id
}

type ProcDefDetailResp struct {
	ID uint64 `json:"id" validate:"required"` // 数据自增id
	objProcDef.ProcDefBaseInfo
	DeployTime string `json:"deployTime"` // 部署时间
	Version    uint64 `json:"version"`    // 版本号
	objCommon.OperatorBaseInfo
}

type ProcDefPageListItem struct {
	ID         uint64 `json:"id"`         // 数据自增id
	Name       string `json:"name"`       // 流程名称
	UserID     string `json:"userid"`     // 用户ID
	Username   string `json:"username"`   // 用户名称
	Company    string `json:"company"`    // 用户所在公司名称
	DeployTime string `json:"deployTime"` // 部署时间
}

type ProcDefPageListResp struct {
	List  []ProcDefPageListItem `json:"list"`  // 数据列表
	Total int64                 `json:"total"` // 数据总条数
}
