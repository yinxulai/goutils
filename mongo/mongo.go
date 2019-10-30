package mongo

import (
	"context"
	"sync"
	"time"

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
var defaultTimeout time.Duration
var globalClient *mongoDriver.Client

// Init mongodb+srv://<username>:<password>@<cluster-address>/test?w=majority
func Init(uri string) {
	defaultOption := mongoOptions.Client().ApplyURI(uri)
	globalClient, _ = mongoDriver.NewClient(defaultOption)
	locker = new(sync.RWMutex)
	inited = true
}

// GetClient 获取一个数据库
func GetClient(url string) *mongoDriver.Client {
	return globalClient
}

// GetContext 延迟执行的 context
func GetContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	go func() {
		time.Sleep(defaultTimeout)
		cancel()
	}()

	return ctx
}

// GetDatabase 获取一个数据库
func GetDatabase(database string) *mongoDriver.Database {
	return globalClient.Database(database)
}

// GetCollectionByDatabase 从 database 获取一个 Collection
func GetCollectionByDatabase(database, collection string) *mongoDriver.Collection {
	return GetDatabase(database).Collection(collection)
}

// GetCollection 从默认 database 获取一个 Collection
func GetCollection(collection string) *mongoDriver.Collection {
	return GetDatabase(defaultDatabase).Collection(collection)
}

// SetDefaultTimeout 设置默认的 database
func SetDefaultTimeout(duartion time.Duration) {
	defaultTimeout = duartion
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
