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

// Service 配置
type Service struct {
	checked   bool
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
		if err = c.Check(); err != nil {
			return "", err
		}
	}

	if c.data[key] == nil {
		return "", fmt.Errorf("%s is nil", key)
	}

	return *c.data[key], nil
}

// Check 检查加载到的数据
func (c *Service) Check() (err error) {
	for _, standard := range c.standards {
		if standard.Required && c.data[standard.Key] == nil {
			panic(fmt.Sprintf("%s is required, %s", standard.Key, standard.Description))
		}
	}

	c.checked = true
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

// AutoLoad 自动加载
func (c *Service) AutoLoad() {
	c.checked = false
}

// LoadFlag 加载启动命令参数
func (c *Service) LoadFlag() {
	c.RLock()
	defer c.RUnlock()
	c.checked = false
	for _, standard := range c.standards {
		value := c.data[standard.Key]
		flag.StringVar(value, standard.Key, standard.Default, standard.Description)
	}
	flag.Parse()
}

// LoadJSONFile 加载文件
func (c *Service) LoadJSONFile(path string) error {
	c.RLock()
	var err error
	defer c.RUnlock()
	c.checked = false

	err = file.ReadJSON(path, &c.data)
	if err != nil {
		return err
	}

	return nil
}

// LoadJSONFiles 加载多个文件
func (c *Service) LoadJSONFiles(paths ...string) error {
	var err error
	c.checked = false
	if len(paths) > 0 {
		for _, path := range paths {
			if err = c.LoadJSONFile(path); err != nil {
				return err
			}
		}
	}
	return nil
}

// CreateJSONTemplate 写入 json 模版
func (c *Service) CreateJSONTemplate(path string) error {
	var err error
	if !c.checked {
		if err = c.Check(); err != nil {
			return err
		}
	}

	return file.WriteJSON(path, false, c.data)
}
