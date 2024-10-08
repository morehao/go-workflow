package errorCode

import "github.com/morehao/go-tools/gerror"

// 内部错误码，1000xx

var ErrorDbInsert = gerror.Error{
	Code: 100000,
	Msg:  "db insert error",
}

var ErrorDbDelete = gerror.Error{
	Code: 100001,
	Msg:  "db delete error",
}

var ErrorDbUpdate = gerror.Error{
	Code: 100002,
	Msg:  "db update error",
}

var ErrorDbFind = gerror.Error{
	Code: 100003,
	Msg:  "db find error",
}
