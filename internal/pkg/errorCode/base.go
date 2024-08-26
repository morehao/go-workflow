package errorCode

import "github.com/morehao/go-tools/gerror"

// 基础错误码，1001xx

var ParamInvalid = gerror.Error{
	Code: 100101,
	Msg:  "参数无效",
}

var UserAuthErr = gerror.Error{
	Code: 100102,
	Msg:  "用户认证失败",
}
