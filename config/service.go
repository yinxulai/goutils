package config

import (
	"flag"
	"fmt"
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

func New() *Service {
	service := new(Service)
	service.files = []string{}
	service.data = make(map[string]*string)
	service.standards = make(map[string]standard)
	return service
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

	// 如果没有读取到数据,且有默认值
	if *c.data[key] == "" && c.standards[key].Default != "" {
		return c.standards[key].Default, nil
	}

	// 如果没有读取到数据,且没有默认值
	if *c.data[key] == "" && c.standards[key].Default == "" {
		return "", fmt.Errorf("%s is nil", key)
	}

	return *c.data[key], nil
}

// check 检查加载到的数据
func (c *Service) load() (err error) {
	c.loadFlag() // 先加载命令行参数

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
			panic(fmt.Sprintf("%s is required, %s", standard.Key, standard.Description))
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

	flag.StringVar(&configPath, "config", "", "指定配置文件,将覆盖程序内设置")

	for _, standard := range c.standards {
		var value string
		cache[standard.Key] = &value
		flag.StringVar(&value, standard.Key, standard.Default, standard.Description)
	}

	flag.Parse()

	// 用户指定了,追加进去配置文件
	if configPath != "" {
		c.files = append(c.files, configPath)
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

// SetStandard 设置定义
func (c *Service) SetStandard(key string, deft string, required bool, description string) {

	c.RLock()
	defer c.RUnlock()
	c.checked = false
	stan := new(standard)

	stan.Key = key
	stan.Default = deft
	stan.Required = required
	stan.Description = description

	c.standards[stan.Key] = *stan
	c.data[stan.Key] = &stan.Default
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
