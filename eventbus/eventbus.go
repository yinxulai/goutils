package eventbus

import (
	"errors"
	"fmt"
	"sync"
)

// EventBus 事件发布订阅
var EventBus *Service
var eventBusIsInitded bool

// Message 消息
type Message struct {
	sync.RWMutex
	values map[string]interface{}
}

// Get 获取
func (message *Message) Get(key string) (value interface{}) {
	message.RLock()
	value, _ = message.values[key]
	message.RUnlock()
	return value
}

// Set 设置
func (message *Message) Set(key string, value interface{}) {
	message.Lock()
	if message.values == nil {
		message.values = make(map[string]interface{})
	}
	message.values[key] = value
	message.Unlock()
}

// Has 是否存在
func (message *Message) Has(key string) bool {
	message.RLock()
	_, being := message.values[key]
	message.RUnlock()
	return being
}

// Keys 所有的 key
func (message *Message) Keys(key string) (keys []string) {
	keys = []string{}
	message.RLock()
	if message.values != nil {
		for key := range message.values {
			keys = append(keys, key)
		}
	}
	message.RUnlock()
	return keys
}

// Values 所有的 value
func (message *Message) Values(key string) (values []interface{}) {
	values = []interface{}{}
	message.RLock()
	if message.values != nil {
		for _, value := range message.values {
			values = append(values, value)
		}
	}
	message.RUnlock()
	return values
}

// Len 长度
func (message *Message) Len(key string) (length int) {
	length = 0
	message.RLock()
	if message.values != nil {
		length = len(message.values)
	}
	message.RUnlock()
	return length
}

// Client 回调
type Client struct {
	ID       int                // 所在频道的编号
	Callback func(msg *Message) // 回调方法
}

// Channel 频道
type Channel struct {
	sync.RWMutex
	Name         string          // 频道名称
	clients      map[int]*Client // 订阅方法列表
	messageCount uint64          // 消息记数
}

// AddClient 添加
func (c *Channel) AddClient(client *Client) {
	c.Lock()
	if c.clients == nil {
		c.clients = make(map[int]*Client)
	}
	c.Unlock()

	c.RLock()
	_, being := c.clients[client.ID]
	c.RUnlock()

	c.Lock()
	if being {
		c.clients[client.ID] = client
	} else {
		client.ID = len(c.clients)
		c.clients[client.ID] = client
	}
	c.Unlock()
}

// PutMessage 向 clients 推送消息
func (c *Channel) PutMessage(message *Message) {
	c.Lock()
	if c.clients == nil {
		c.clients = make(map[int]*Client)
	}
	c.messageCount++ // 消息统计
	for _, client := range c.clients {
		go client.Callback(message)
	}
	c.Unlock()
}

// Service 发布订阅服务
type Service struct {
	sync.RWMutex
	channels map[string]*Channel
}

// Emit 触发事件
func (ebs *Service) Emit(key string, msg *Message) {
	ebs.Lock()
	if ebs.channels == nil {
		ebs.channels = make(map[string]*Channel)
	}
	channel, being := ebs.channels[key]

	if being {
		channel.PutMessage(msg)
	}
	ebs.Unlock()
}

// On 订阅事件
func (ebs *Service) On(key string, client *Client) {
	ebs.Lock()
	if ebs.channels == nil {
		ebs.channels = make(map[string]*Channel)
	}
	ebs.Unlock()

	ebs.RLock()
	channel, being := ebs.channels[key]
	ebs.RUnlock()

	ebs.Lock()
	if !being {
		ebs.channels[key] = &Channel{Name: key}
	}
	channel.AddClient(client)
	ebs.Unlock()
}

// Len 长度
func (ebs *Service) Len() int {
	return len(ebs.channels)
}

// ChannelNames Channel 的 name
func (ebs *Service) ChannelNames() (names []string) {
	names = []string{}
	ebs.Lock()
	if ebs.channels != nil {
		for name := range ebs.channels {
			names = append(names, name)
		}
	}
	ebs.Unlock()
	return names
}

// GetChannel 获取 Channel 对象
func (ebs *Service) GetChannel(name string) (channel *Channel) {
	ebs.Lock()
	if ebs.channels != nil {
		channel = ebs.channels[name]
	}
	ebs.Unlock()
	return channel
}

// CreateMessage 创建一个 Message
func (ebs *Service) CreateMessage() *Message {
	return &Message{}
}

// CreateClient 创建一个 Client
func (ebs *Service) CreateClient(callback func(msg *Message)) *Client {
	return &Client{
		Callback: callback,
	}
}

// init 初始化
func init() {
	if eventBusIsInitded {
		fmt.Print(errors.New("EventBusService is inited! "))
		return
	}
	EventBus = &Service{}
	eventBusIsInitded = true
}
