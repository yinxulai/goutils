package config

// Standard 结构声明
type Standard struct {
	Key         string
	Name        string
	Required    bool
	Description string
}

// Config 配置
type Config struct {
	data      map[string]string
	standards map[string]Standard
}

// Set 获取一个配置
func (c *Config) Set(key string) {

}

// Get 获取一个配置
func (c *Config) Get(key string) (value string, err error) {
	return "", err
}

// SetStandard 设置定义
func (c *Config) SetStandard(stans ...*Standard) {
	if stans != nil && len(stans) > 0 {
		for _, stan := range stans {
			c.standards[stan.Key] = *stan
		}
	}
}

// AutoLoad 自动加载
func (c *Config) AutoLoad(merge bool) {

}

// LoadFile 加载文件
func (c *Config) LoadFile() {

}

// LoadEnv 加载环境变量
func (c *Config) LoadEnv() {

}
