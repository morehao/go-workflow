package dtoProcess

import (
	"go-workflow/internal/app/object/objCommon"
	"go-workflow/internal/app/object/objProcess"
)

type ProcInstStartResp struct {
	ID uint64 `json:"id"` // 数据自增id
}

type ProcInstDetailResp struct {
	ID uint64 `json:"id" validate:"required"` // 数据自增id
	objProcess.ProcInstBaseInfo
	objCommon.OperatorBaseInfo
}

type ProcInstPageListItem struct {
	ID uint64 `json:"id" validate:"required"` // 数据自增id
	objProcess.ProcInstBaseInfo
}

type ProcInstPageListResp struct {
	List  []ProcInstPageListItem `json:"list"`  // 数据列表
	Total int64                  `json:"total"` // 数据总条数
}

type CreatedPageListResp struct {
	List  []ProcInstPageListItem `json:"list"`  // 数据列表
	Total int64                  `json:"total"` // 数据总条数
}

type TodoPageListResp struct {
	List  []ProcInstPageListItem `json:"list"`  // 数据列表
	Total int64                  `json:"total"` // 数据总条数
}

type NotifyPageListResp struct {
	List  []ProcInstPageListItem `json:"list"`  // 数据列表
	Total int64                  `json:"total"` // 数据总条数
}

type TaskCompleteResp struct {
	ID uint64 `json:"id"` // 数据自增id
}
type TaskWithDrawResp struct {
}
