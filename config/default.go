package config

// defaultService Config
var defaultService *Service

func init() {
	defaultService = New()
}

// Set 获取一个配置
func Set(key string, value string) {
	defaultService.Set(key, value)
}

// Get 获取一个配置
func Get(key string) (value string, err error) {
	return defaultService.Get(key)
}

// MustGet 获取一个配置
func MustGet(key string) (value string) {
	return defaultService.MustGet(key)
}

// AddFile 获取一个配置
func AddFile(path string) {
	defaultService.AddFile(path)
}

// SetStandard 设置定义
func SetStandard(key string, deft string, required bool, description string) {
	defaultService.SetStandard(key, deft, required, description)
}

// CreateJSONTemplate 写入 json 模版
func CreateJSONTemplate(path string) error {
	return defaultService.CreateJSONTemplate(path)
}
