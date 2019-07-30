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
	service.data = make(map[string]*string)
	service.standards = make(map[string]standard)
	service.SetStandard("config", "./config.json", false, "自定义配置文件")
	return service
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

// AutoLoad TODO: 自动加载
// 自动查找 .json 文件并尝试解析
// 全部处理完后检查是否必要参数都已齐全
func (c *Service) AutoLoad() {
	c.LoadFlag()
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
	var data map[string]*string

	// 如果环境变量指定了文件 则加载指定文件
	// 程序内的设定则无效 变量等级最高
	if c.data["config"] != nil {
		path = *c.data["config"]
	}

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
