package dtoUser

import (
	"go-workflow/internal/app/object/objCommon"
	"go-workflow/internal/app/object/objUser"
)

type UserCreateReq struct {
	objUser.UserBaseInfo
}

type UserUpdateReq struct {
	ID uint64 `json:"id" validate:"required" label:"数据自增id"` // 数据自增id
	objUser.UserBaseInfo
}

type UserDetailReq struct {
	ID uint64 `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}

type UserPageListReq struct {
	objCommon.PageQuery
}

type UserDeleteReq struct {
	ID uint64 `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}
