package config

// DefaultConfig Config
var DefaultConfig *Service

func init() {
	DefaultConfig = new(Service)
	DefaultConfig.data = make(map[string]*string)
	DefaultConfig.standards = make(map[string]standard)
}

// Set 获取一个配置
func Set(key string, value string) {
	DefaultConfig.Set(key, value)
}

// Get 获取一个配置
func Get(key string) (value string, err error) {
	return DefaultConfig.Get(key)
}

// SetStandard 设置定义
func SetStandard(key string, deft string, required bool, description string) {
	DefaultConfig.SetStandard(key, deft, required, description)
}

// AutoLoad 自动加载
func AutoLoad() {
	DefaultConfig.AutoLoad()
}

// LoadFlag 加载启动命令参数
func LoadFlag() {
	DefaultConfig.LoadFlag()
}

// LoadJSON 加载文件
func LoadJSON(path string) error {
	return DefaultConfig.LoadJSON(path)
}

// LoadJSONs 加载多个文件
func LoadJSONs(paths ...string) error {
	return DefaultConfig.LoadJSONs(paths...)
}

// CreateJSONTemplate 写入 json 模版
func CreateJSONTemplate(path string) error {
	return DefaultConfig.CreateJSONTemplate(path)
}
