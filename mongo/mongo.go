package mongo

import (
	"sync"

	"github.com/yinxulai/goutils/logger"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

var uri string
var inited bool
var locker *sync.RWMutex
var defaultDatabase string
var defaultCollection string
var defaultLogger *logger.Logger
var globalClient *mongoDriver.Client

func init() {
	locker = new(sync.RWMutex)
}

// GetClient 获取一个数据库
func GetClient(url string) *mongoDriver.Client {
	return globalClient
}

// GetDatabase 获取一个数据库
func GetDatabase(database string) *mongoDriver.Database {
	return globalClient.Database(database)
}

// GetCollection 获取一个 Collection
func GetCollection(database, collection string) *mongoDriver.Collection {
	return GetDatabase(database).Collection(collection)
}

// SetDefaultDatabase 设置默认的 database
func SetDefaultDatabase(database string) {
	defaultDatabase = database
}

// SetDefaultCollection 设置默认的 collection
func SetDefaultCollection(collection string) {
	defaultCollection = collection
}

// SetLogger 设置日志插件
func SetLogger(logger *logger.Logger) {
	defaultLogger = logger
}

// SetURI 设置服务地址
// uri 格式：mongodb+srv://<username>:<password>@<cluster-address>/test?w=majority
func SetURI(uri string) {
	defaultOption := mongoOptions.Client().ApplyURI(uri)
	globalClient, _ = mongoDriver.NewClient(defaultOption)
	inited = true
}
