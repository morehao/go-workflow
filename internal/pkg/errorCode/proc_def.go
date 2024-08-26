package errorCode

import "github.com/morehao/go-tools/gerror"

// 审批流程定义相关错误码，1002xx

var ProcDefSaveErr = gerror.Error{
	Code: 100200,
	Msg:  "保存审批流程定义失败",
}

var ProcDefDeleteErr = gerror.Error{
	Code: 100201,
	Msg:  "删除审批流程定义失败",
}

var ProcDefUpdateErr = gerror.Error{
	Code: 100202,
	Msg:  "修改审批流程定义失败",
}

var ProcDefGetDetailErr = gerror.Error{
	Code: 100203,
	Msg:  "查看审批流程定义失败",
}

var ProcDefGetPageListErr = gerror.Error{
	Code: 100204,
	Msg:  "查看审批流程定义列表失败",
}

var ProcDefNotExistErr = gerror.Error{
	Code: 100205,
	Msg:  "审批流程定义不存在",
}
