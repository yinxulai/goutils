package restful

import (
	"encoding/json"
)

// Package Restful 包结构
type Package struct {
	Data    interface{} // 数据
	State   uint16      // 状态码
	Message string      // 错误信息
}

// NewString 获取一个 Restful 结构
func NewString(state uint16, message string, data interface{}) string {
	return string(NewByte(state, message, data))
}

// NewByte 获取一个 Restful 结构
func NewByte(state uint16, message string, data interface{}) []byte {
	byteData, _ := json.Marshal(New(state, message, data))
	return byteData
}

// New 获取一个 Restful 结构
func New(state uint16, message string, data interface{}) *Package {
	entity := new(Package)
	entity.Data = data
	entity.State = state
	entity.Message = message

	// 使用 code 的含义填充默认 message
	if entity.Message == "" {
		entity.Message = StateMeaning[state]
	}

	return entity
}
