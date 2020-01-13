package config

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/yinxulai/goutils/file"
)

// Standard 结构声明
type standard struct {
	Key         string
	Default     string
	Required    bool
	Description string
}

// New New
func New() *Service {
	configService := new(Service)
	configService.files = []string{}
	configService.data = make(map[string]*string)
	configService.standards = make(map[string]standard)
	return configService
}

// Service 配置
type Service struct {
	loaded    bool     // 加载完成
	checked   bool     // 检查完成
	files     []string // 待加载的文件
	data      map[string]*string
	standards map[string]standard
	sync.RWMutex
}

// Set 获取一个配置
func (c *Service) Set(key string, value string) {
	c.RLock()
	defer c.RUnlock()
	c.data[key] = &value
}

// Get 获取一个配置
func (c *Service) Get(key string) (value string, err error) {
	if !c.checked {
		if err = c.check(); err != nil {
			return "", err
		}
	}

	if !c.loaded {
		if err = c.load(); err != nil {
			return "", err
		}
	}

	if *c.data[key] == "" {
		return "", fmt.Errorf("config: %s is nil", key)
	}

	return *c.data[key], nil
}

// MustGet 获取一个配置
func (c *Service) MustGet(key string) (value string) {
	value, err := c.Get(key)
	if err != nil {
		panic(err)
	}

	return value
}

// SetStandard 设置定义
func (c *Service) SetStandard(key string, deft string, required bool, description string) {
	c.RLock()
	defer c.RUnlock()

	// 恢复为未 check 状态
	c.checked = false

	// 检查 key 格式
	c.mustKeyCheck(key)

	// 记录
	stan := new(standard)
	stan.Key = key
	stan.Default = deft
	stan.Required = required
	stan.Description = description
	c.standards[stan.Key] = *stan

	// 注册 flag
	var value string
	value = stan.Default
	c.data[stan.Key] = &value
	flag.StringVar(&value, key, stan.Default, description)
}

// AddFile 加载文件
func (c *Service) AddFile(path string) {
	c.RLock()
	defer c.RUnlock()
	c.files = append(c.files, path)
}

// CreateJSONTemplate 写入 json 模版
func (c *Service) CreateJSONTemplate(path string) error {
	var err error
	if !c.checked {
		if err = c.check(); err != nil {
			return err
		}
	}

	return file.WriteJSON(path, false, c.data)
}

// check 检查加载到的数据
func (c *Service) load() (err error) {
	c.loadEnv()  // 先加载环境变量
	c.loadFlag() // 再加载命令行参数

	// 读取所有文件
	for _, filePath := range c.files {
		// 读取文件
		err = c.loadFile(filePath)
		if err != nil {
			return err
		}
	}

	c.loaded = true
	return nil
}

// check 检查加载到的数据
func (c *Service) check() (err error) {
	for _, standard := range c.standards {
		if standard.Required && c.data[standard.Key] == nil {
			panic(fmt.Sprintf("config: %s is required, %s", standard.Key, standard.Description))
		}
	}

	c.checked = true
	return nil
}

// loadFlag 加载启动命令参数
func (c *Service) loadFlag() {
	c.RLock()
	defer c.RUnlock()
	c.checked = false
	var configPath string
	cache := make(map[string]*string)

	// 如果用户没有覆盖默认的 config 行为
	if _, exists := c.standards["config"]; !exists {
		c.SetStandard("config", "", false, "指定配置文件,将覆盖程序内设置")
	}

	flag.Parse()

	// 用户指定了,追加进去配置文件
	if configPath != "" {
		c.files = append(c.files, *cache["config"])
	}

	for key, value := range cache {
		if value != nil && *value != "" {
			c.data[key] = value
		}
	}
}

// loadFile 加载文件
func (c *Service) loadFile(path string) error {
	c.RLock()
	var err error
	defer c.RUnlock()
	c.checked = false
	var data map[string]*string

	err = file.ReadJSON(path, &data)
	if err != nil {
		return err
	}

	for name, value := range data {
		if value != nil {
			c.data[name] = value
		}
	}

	return nil
}

// 加载环境变量
func (c *Service) loadEnv() {
	c.RLock()
	defer c.RUnlock()
	c.checked = false
	for key, standard := range c.standards {
		value := os.Getenv(strings.ToUpper(standard.Key))
		if value != "" {
			c.data[key] = &value
		}
	}
}

// MustKeyCheck 检查
func (c *Service) mustKeyCheck(key string) {
	if key == "" {
		panic(fmt.Errorf("config: 配置 key 名不允许为空"))
	}

	matched, err := regexp.MatchString("^[a-zA-Z0-9_]*$", key)

	if err != nil {
		panic(fmt.Errorf("config: 配置 key 名检查错误: %v", err))
	}

	if !matched {
		panic(fmt.Errorf("config: 配置 key 名仅允许大小写字母、数字、下划线"))
	}
}
