package config

// defaultService Config
var defaultService *Service

func init() {
	defaultService = new(Service)
	defaultService.data = make(map[string]*string)
	defaultService.standards = make(map[string]standard)
}

// Set 获取一个配置
func Set(key string, value string) {
	defaultService.Set(key, value)
}

// Get 获取一个配置
func Get(key string) (value string, err error) {
	return defaultService.Get(key)
}

// SetStandard 设置定义
func SetStandard(key string, deft string, required bool, description string) {
	defaultService.SetStandard(key, deft, required, description)
}

// AutoLoad 自动加载
func AutoLoad() {
	defaultService.AutoLoad()
}

// LoadFlag 加载启动命令参数
func LoadFlag() {
	defaultService.LoadFlag()
}

// LoadJSONFile 加载文件
func LoadJSONFile(path string) error {
	return defaultService.LoadJSONFile(path)
}

// LoadJSONFiles 加载多个文件
func LoadJSONFiles(paths ...string) error {
	return defaultService.LoadJSONFiles(paths...)
}

// CreateJSONTemplate 写入 json 模版
func CreateJSONTemplate(path string) error {
	return defaultService.CreateJSONTemplate(path)
}
