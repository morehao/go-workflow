package errorCode

import "github.com/morehao/go-tools/gerror"

var ProcInstCreateErr = gerror.Error{
	Code: 100300,
	Msg:  "创建审批流程实例失败",
}

var ProcInstDeleteErr = gerror.Error{
	Code: 100301,
	Msg:  "删除审批流程实例失败",
}

var ProcInstUpdateErr = gerror.Error{
	Code: 100302,
	Msg:  "修改审批流程实例失败",
}

var ProcInstGetDetailErr = gerror.Error{
	Code: 100303,
	Msg:  "查看审批流程实例失败",
}

var ProcInstGetPageListErr = gerror.Error{
	Code: 100304,
	Msg:  "查看审批流程实例列表失败",
}

var ProcInstNotExistErr = gerror.Error{
	Code: 100305,
	Msg:  "审批流程实例不存在",
}
var ProcInstCreatedPageListErr = gerror.Error{
	Code: 100306,
	Msg:  "我创建的流程实例分页列表失败",
}
var ProcInstTodoPageListErr = gerror.Error{
	Code: 100307,
	Msg:  "待我审批的流程实例分页列表失败",
}
var ProcInstNotifyPageListErr = gerror.Error{
	Code: 100308,
	Msg:  "抄送我的的流程实例分页列表失败",
}
