package config

// globalService Config
var globalService *configService

func init() {
	globalService = New()
}

// Set 获取一个配置
func Set(key string, value string) {
	globalService.Set(key, value)
}

// Get 获取一个配置
func Get(key string) (value string, err error) {
	return globalService.Get(key)
}

// MustGet 获取一个配置
func MustGet(key string) (value string) {
	return globalService.MustGet(key)
}

// AddFile 获取一个配置
func AddFile(path string) {
	globalService.AddFile(path)
}

// LoadFlag 加载启动命令参数
func LoadFlag() {
	globalService.loadFlag()
}

// SetStandard 设置定义
func SetStandard(key string, deft string, required bool, description string) {
	globalService.SetStandard(key, deft, required, description)
}

// CreateJSONTemplate 写入 json 模版
func CreateJSONTemplate(path string) error {
	return globalService.CreateJSONTemplate(path)
}
