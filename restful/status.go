package restful

// 常见错误码
const (
	OK                  = uint16(200)  // 正常
	CREATED             = uint16(201)  // 创建完成
	ACCEPTED            = uint16(202)  // 更新成功
	NOCONTENT           = uint16(204)  // 表示资源有空
	MOVEDPERMANENTLY    = uint16(301)  // 资源的URI已被更新
	NOTMODIFIED         = uint16(304)  // 源未更改
	BADREQUEST          = uint16(400)  // 坏请求 一般参数错误
	UNAUTHORIZED        = uint16(401)  // 未授权
	FORBIDDEN           = uint16(403)  // 被禁止访问 墙住
	NOTFOUND            = uint16(404)  // 请求的资源不存在
	METHODNOTALLOWED    = uint16(405)  // 不正确的请求方式
	NOTACCEPTABLE       = uint16(406)  // 无法解析的服务端响应
	CONFLICT            = uint16(409)  // 发生冲突错误 比如修改不存在的东西、创建已存在的东西
	INTERNALSERVERERROR = uint16(500)  // 服务器内部错误
	SERVICEUNAVAILABLE  = uint16(503)  // 服务器不可用 维护、暂停、维护等
	DBERR               = uint16(4001) // 数据库错误
	NODATA              = uint16(4002) // 无数据
	DATAEXIST           = uint16(4003) // 数据已存在
	DATAERR             = uint16(4004) // 数据错误
	SESSIONERR          = uint16(4101) // session 会话错误
	LOGINERR            = uint16(4102) // 登陆错误
	PARAMERR            = uint16(4103) // 参数错误
	USERERR             = uint16(4104) // 用户错误
	ROLEERR             = uint16(4105) // 角色错误
	PWDERR              = uint16(4106) // 密码错误
	REQERR              = uint16(4201) // 请求错误
	IPERR               = uint16(4202) // IP 错误
	THIRDERR            = uint16(4301) // 第三方服务错误
	IOERR               = uint16(4302) // IO 错误
	SERVERERR           = uint16(4500) // 服务错误
	UNKOWNERR           = uint16(4501) // 未知错误
)

// StateMeaning State 的含义
var StateMeaning = map[uint16]string{
	OK:                  "完成",
	CREATED:             "创建成功",
	ACCEPTED:            "更新成功",
	NOCONTENT:           "内容为空",
	MOVEDPERMANENTLY:    "URI 已变更",
	NOTMODIFIED:         "源未更改",
	BADREQUEST:          "无效请求",
	UNAUTHORIZED:        "未授权",
	FORBIDDEN:           "被禁止访问",
	NOTFOUND:            "资源不存在",
	METHODNOTALLOWED:    "错误的请求方式",
	NOTACCEPTABLE:       "无法解析的服务端响应",
	CONFLICT:            "冲突错误",
	INTERNALSERVERERROR: "服务器内部错误",
	SERVICEUNAVAILABLE:  "服务暂时不可用，可能在维护",
	DBERR:               "数据库错误",
	NODATA:              "无数据",
	DATAEXIST:           "数据已存在",
	DATAERR:             "数据错误",
	SESSIONERR:          "session 会话错误",
	LOGINERR:            "登陆错误",
	PARAMERR:            "参数错误",
	USERERR:             "用户错误",
	ROLEERR:             "角色错误",
	PWDERR:              "密码错误",
	REQERR:              "请求错误",
	IPERR:               "IP 错误",
	THIRDERR:            "第三方服务错误",
	IOERR:               "IO 错误",
	SERVERERR:           "服务错误",
	UNKOWNERR:           "未知错误",
}
