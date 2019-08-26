package pattern

import "regexp"

// 常见正则
var (
	IPV4     = regexp.MustCompile(`^((0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])\.){3}(0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])$`) // 常见用户昵称 可以包含任何字符
	Email    = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)                                    // 常见用户昵称 可以包含任何字符
	Phone    = regexp.MustCompile(`^1\d{10}$`)                                                                       // 宽松的手机号 1 开头 11 位
	Nickname = regexp.MustCompile(`^.{2,128}$`)                                                                      // 常见用户昵称 可以包含任何字符
	Username = regexp.MustCompile(`^[a-zA-Z0-9]{6,128}$`)                                                            // 常见用户名 只能是 英文（大小写）+ 数字                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    // 安全级别中等 必须包含 大写、小写、数字其中两项
	Password = regexp.MustCompile(`^[a-zA-Z0-9!"#$%&'()*+,\-./:;<=>?@\[\\\]^_\{\|\}\~]{6,128}$`)                     // 安全级别高级 必须同时包含大写、小写、数字和特殊字符其中三项
)
