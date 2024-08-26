package dtoProcDef

import (
	"go-workflow/internal/app/object/objCommon"
	"go-workflow/internal/app/object/objProcDef"
)

type ProcDefSaveReq struct {
	objProcDef.ProcDefBaseInfo
}

type ProcDefDetailReq struct {
	ID uint64 `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}

type ProcDefPageListReq struct {
	objCommon.PageQuery
}

type ProcDefDeleteReq struct {
	ID uint64 `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}
