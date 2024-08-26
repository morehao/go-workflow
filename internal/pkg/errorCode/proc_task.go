package errorCode

import "github.com/morehao/go-tools/gerror"

var ProcTaskCompleteErr = gerror.Error{
	Code: 100400,
	Msg:  "完成审批流程任务失败",
}

var ProcTaskWithDrawErr = gerror.Error{
	Code: 100401,
	Msg:  "撤回审批流程任务失败",
}

var ProcTaskNotExistErr = gerror.Error{
	Code: 100402,
	Msg:  "撤回审批流程任务失败",
}
