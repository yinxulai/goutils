package sqldb

import (
	"time"

	// 导入 mysql 驱动, 触发 init
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var initd = false
var globalDB *sqlx.DB
var globalStmtMap map[string]*sqlx.Stmt
var globalNamedStmtMap map[string]*sqlx.NamedStmt

// Init 初始化连接池
func Init(driver, source string) {
	db := sqlx.MustConnect(driver, source)
	globalStmtMap = make(map[string]*sqlx.Stmt)
	globalNamedStmtMap = make(map[string]*sqlx.NamedStmt)
	db.MapperFunc(func(name string) string { return name })
	db.SetConnMaxLifetime(30 * time.Minute) // 最大存活周期 30 分钟
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	globalDB = db

	initd = true
}

// GetDB GetDB
func GetDB() *sqlx.DB {
	return globalDB
}

// CreateNamedStmt CreateNamedStmt
func CreateNamedStmt(sql string) *sqlx.NamedStmt {
	stmt, err := globalDB.PrepareNamed(sql)
	if err != nil {
		panic(err)
	}
	globalNamedStmtMap[sql] = stmt
	return globalNamedStmtMap[sql]
}

// CreateStmt CreateStmt
func CreateStmt(sql string) *sqlx.Stmt {
	stmt, err := globalDB.Preparex(sql)
	if err != nil {
		panic(err)
	}

	globalStmtMap[sql] = stmt
	return globalStmtMap[sql]
}
