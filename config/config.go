package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	"github.com/yinxulai/goutils/file"
)

// Standard 结构声明
type Standard struct {
	Key         string
	Name        string
	Default     string
	Required    bool
	Description string
}

// Config 配置
type Config struct {
	checked   bool
	data      map[string]*string
	standards map[string]Standard
	sync.RWMutex
}

// Set 获取一个配置
func (c *Config) Set(key string, value string) {
	c.RLock()
	defer c.RUnlock()
	c.data[key] = &value
}

// Get 获取一个配置
func (c *Config) Get(key string) (value string, err error) {
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
func (c *Config) Check() (err error) {
	for _, standard := range c.standards {
		if standard.Required && c.data[standard.Key] == nil {
			return fmt.Errorf("%s is required, %s", standard.Key, standard.Description)
		}
	}

	c.checked = true
	return nil
}

// SetStandard 设置定义
func (c *Config) SetStandard(stans ...*Standard) {
	c.RLock()
	defer c.RUnlock()
	c.checked = false
	if stans != nil && len(stans) > 0 {
		for _, stan := range stans {
			c.standards[stan.Key] = *stan
			c.data[stan.Key] = &stan.Default
		}
	}
}

// AutoLoad 自动加载
func (c *Config) AutoLoad() {
	c.checked = false
}

// LoadFlag 加载启动命令参数
func (c *Config) LoadFlag() {
	c.RLock()
	defer c.RUnlock()
	c.checked = false
	for _, standard := range c.standards {
		value := c.data[standard.Key]
		flag.StringVar(value, standard.Key, standard.Default, standard.Description)
	}
}

// LoadJSON 加载文件
func (c *Config) LoadJSON(path string) error {
	c.RLock()
	var err error
	var data []byte
	defer c.RUnlock()
	c.checked = false

	data, err = file.ReadAll(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, c.data)
	if err != nil {
		return err
	}

	return nil
}

// LoadJSONs 加载多个文件
func (c *Config) LoadJSONs(paths ...string) error {
	var err error
	c.checked = false
	if len(paths) > 0 {
		for _, path := range paths {
			if err = c.LoadJSON(path); err != nil {
				return err
			}
		}
	}
	return nil
}

// CreateJSONTemplate 写入 json 模版
func (c *Config) CreateJSONTemplate(path string) error {
	var err error
	if !c.checked {
		if err = c.Check(); err != nil {
			return err
		}
	}

	return file.WriteJSON(path, false, c.data)
}
