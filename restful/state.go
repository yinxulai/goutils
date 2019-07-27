package restful

// 常见错误码
const (
	OK                  = uint16(200) // 正常
	CREATED             = uint16(201) // 创建完成
	ACCEPTED            = uint16(202) // 更新成功
	NOCONTENT           = uint16(204) // 表示资源有空
	MOVEDPERMANENTLY    = uint16(301) // 资源的URI已被更新
	NOTMODIFIED         = uint16(304) // 源未更改
	BADREQUEST          = uint16(400) // 坏请求 一般参数错误
	UNAUTHORIZED        = uint16(401) // 未授权
	FORBIDDEN           = uint16(403) // 被禁止访问 墙住
	NOTFOUND            = uint16(404) // 请求的资源不存在
	METHODNOTALLOWED    = uint16(405)
	NOTACCEPTABLE       = uint16(406)
	CONFLICT            = uint16(409)
	INTERNALSERVERERROR = uint16(500)
	SERVICEUNAVAILABLE  = uint16(503)
	DBERR               = uint16(4001)
	NODATA              = uint16(4002)
	DATAEXIST           = uint16(4003)
	DATAERR             = uint16(4004)
	SESSIONERR          = uint16(4101)
	LOGINERR            = uint16(4102)
	PARAMERR            = uint16(4103)
	USERERR             = uint16(4104)
	ROLEERR             = uint16(4105)
	PWDERR              = uint16(4106)
	REQERR              = uint16(4201)
	IPERR               = uint16(4202)
	THIRDERR            = uint16(4301)
	IOERR               = uint16(4302)
	SERVERERR           = uint16(4500)
	UNKOWNERR           = uint16(4501)
)
