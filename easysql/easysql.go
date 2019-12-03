package easysql

import (
	"database/sql"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	// 导入 mysql 驱动, 触发 init
	_ "github.com/go-sql-driver/mysql"
)

// TODO: SQL 注入问题

var globalDB *sql.DB
var config *Config

const beginStatus = 1

// SQL 结构化的 sql 语句
type SQL struct {
	db          *sql.DB // 数据库连接
	fields      string  // 字段
	orderBy     string  // 排序条件
	tableName   string  // 表名
	execString  string  // 执行sql语句
	limitNumber string  // 限制条数
	whereString string  // where语句

	tx           *sql.Tx // 原生事务
	commitSign   int8    // 提交标记，控制是否提交事务
	rollbackSign bool    // 回滚标记，控制是否回滚事务

	config *Config
}

// Config 配置
type Config struct {
	Source string
	Driver string
}

// Init 初始化连接池
func Init(driver, source string) {
	sqlStruct := new(SQL)
	sqlStruct.config = new(Config)
	sqlStruct.config.Driver = driver
	sqlStruct.config.Source = source
	newDB, err := sql.Open(sqlStruct.config.Driver, sqlStruct.config.Source)
	sqlStruct.handlerError(err)
	newDB.SetMaxOpenConns(2000)             // 最大链接
	newDB.SetMaxIdleConns(1000)             // 空闲连接，也就是连接池里面的数量
	newDB.SetConnMaxLifetime(7 * time.Hour) // 设置最大生成周期是7个小时
	globalDB = newDB
}

// GetConn GetConn
func GetConn() *SQL {
	sql := new(SQL)
	sql.db = globalDB
	return sql
}

// handlerError 检查错误
func (SQL *SQL) handlerError(err error) bool {
	if err != nil {
		log.Fatal("错误：", err)
		return true
	}
	return false
}

// Close 关闭数据库连接（不是释放返回连接池）
func (SQL *SQL) Close() error {
	if SQL.tx != nil {
		// 如果有事务就回滚
		err := SQL.Rollback()
		if err != nil {
			return err
		}
	}

	if SQL.db != nil {
		err := SQL.db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// Select 查询方法
func (SQL *SQL) Select(tableName string, field []string) *SQL {
	var allField string
	if field == nil || len(field) <= 0 {
		allField = "*"
	} else {
		allField = strings.Join(field, ",")
	}

	SQL.fields = "SELECT " + allField + " FROM `" + tableName + "`"
	SQL.tableName = tableName
	return SQL
}

// Where Where
func (SQL *SQL) Where(cond map[string]string) *SQL {
	var whereString = ""
	if len(cond) != 0 {
		whereString = " WHERE "
		for key, value := range cond {
			if !strings.Contains(key, "=") && !strings.Contains(key, ">") && !strings.Contains(key, "<") {
				key += "="
			}
			whereString += key + "'" + value + "'" + " AND "
		}
	}
	// 删除所有字段最后一个
	whereString = strings.TrimSuffix(whereString, "AND ")
	SQL.whereString = whereString
	return SQL
}

// Limit Limit
func (SQL *SQL) Limit(number int) *SQL {
	SQL.limitNumber = " LIMIT " + strconv.Itoa(number)
	return SQL
}

// OrderByString OrderByString
func (SQL *SQL) OrderByString(orderString ...string) *SQL {
	if len(orderString) > 2 || len(orderString) <= 0 {
		log.Fatal("传入参数错误")
	} else if len(orderString) == 1 {
		SQL.orderBy = " ORDER BY " + orderString[0] + " ASC"
	} else {
		SQL.orderBy = " ORDER BY " + orderString[0] + " " + orderString[1]
	}
	return SQL
}

// ToString 获取目前的 sql 语句
func (SQL *SQL) ToString() string {
	return SQL.fields + SQL.whereString + SQL.orderBy + SQL.limitNumber
}

// Update Update
func (SQL *SQL) Update(tableName string, str map[string]string) (int64, error) {
	var tempStr = ""
	var allValue []interface{}
	for key, value := range str {
		tempStr += key + "=" + "?" + ","
		allValue = append(allValue, value)
	}
	tempStr = strings.TrimSuffix(tempStr, ",")
	SQL.execString = "UPDATE `" + tableName + "` SET " + tempStr
	var allStr = SQL.execString + SQL.whereString

	var err error
	var stmt *sql.Stmt
	if SQL.tx != nil {
		stmt, err = SQL.tx.Prepare(allStr)
	} else {
		stmt, err = SQL.db.Prepare(allStr)
	}

	if SQL.handlerError(err) {
		return 0, err
	}

	res, err := stmt.Exec(allValue...)
	if SQL.handlerError(err) {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if SQL.handlerError(err) {
		return 0, err
	}

	return rows, nil
}

// Delete 删除方法
func (SQL *SQL) Delete(tableName string) (int64, error) {
	var tempStr = ""
	tempStr = "DELETE FROM `" + tableName + "`" + SQL.whereString

	var err error
	var stmt *sql.Stmt
	if SQL.tx != nil {
		stmt, err = SQL.tx.Prepare(tempStr)
	} else {
		stmt, err = SQL.db.Prepare(tempStr)
	}

	if SQL.handlerError(err) {
		return 0, err
	}
	res, err := stmt.Exec()
	if SQL.handlerError(err) {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if SQL.handlerError(err) {
		return 0, err
	}

	return rows, nil
}

// Insert 插入方法
func (SQL *SQL) Insert(tableName string, data map[string]string) (int64, error) {
	var allField = ""
	var allValue = ""
	var allTrueValue []interface{}
	if len(data) != 0 {
		for key, value := range data {
			allField += key + ","
			allValue += "?" + ","
			allTrueValue = append(allTrueValue, value)
		}
	}
	allValue = strings.TrimSuffix(allValue, ",")
	allField = strings.TrimSuffix(allField, ",")
	allValue = "(" + allValue + ")"
	allField = "(" + allField + ")"
	var theStr = "INSERT INTO `" + tableName + "` " + allField + " VALUES " + allValue

	var err error
	var stmt *sql.Stmt
	if SQL.tx != nil {
		stmt, err = SQL.tx.Prepare(theStr)
	} else {
		stmt, err = SQL.db.Prepare(theStr)
	}

	if SQL.handlerError(err) {
		return 0, err
	}
	res, err := stmt.Exec(allTrueValue...)
	if SQL.handlerError(err) {
		return 0, err
	}
	SQL.handlerError(err)
	id, err := res.LastInsertId()
	if SQL.handlerError(err) {
		return 0, err
	}
	return id, nil
}

// Pagination 分页查询
func (SQL *SQL) Pagination(Page int, Limit int) (map[string]interface{}, error) {
	res, err := GetConn().Select(SQL.tableName, []string{"COUNT(*) as count"}).QueryRow()
	if SQL.handlerError(err) {
		return nil, err
	}
	count, err := strconv.Atoi(res["count"])
	if SQL.handlerError(err) {
		return nil, err
	}
	// 计算总页码数
	totalPage := int(math.Ceil(float64(count) / float64(Limit)))
	if Page > totalPage {
		Page = totalPage
	}
	if Page <= 0 {
		Page = 1
	}
	// 计算偏移量
	setOff := (Page - 1) * Limit
	queryString := SQL.fields + SQL.whereString + SQL.orderBy + " LIMIT " + strconv.Itoa(setOff) + "," + strconv.Itoa(Limit)
	rows, err := SQL.db.Query(queryString)
	defer rows.Close()
	if SQL.handlerError(err) {
		return nil, err
	}
	Column, err := rows.Columns()
	if SQL.handlerError(err) {
		return nil, err
	}
	// 创建一个查询字段类型的slice
	values := make([]sql.RawBytes, len(Column))
	// 创建一个任意字段类型的slice
	scanArgs := make([]interface{}, len(values))
	// 创建一个slice保存所以的字段
	var allRows []interface{}
	for i := range values {
		// 把values每个参数的地址存入scanArgs
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		// 把存放字段的元素批量放进去
		err = rows.Scan(scanArgs...)
		if SQL.handlerError(err) {
			return nil, err
		}
		tempRow := make(map[string]string, len(Column))
		for i, col := range values {
			var key = Column[i]
			tempRow[key] = string(col)
		}
		allRows = append(allRows, tempRow)
	}
	returnData := make(map[string]interface{})
	returnData["totalPage"] = totalPage
	returnData["currentPage"] = Page
	returnData["rows"] = allRows
	return returnData, nil
}

// QueryAll QueryAll
func (SQL *SQL) QueryAll() ([]map[string]string, error) {
	var queryString = SQL.ToString()
	rows, err := SQL.db.Query(queryString)
	defer rows.Close()
	if SQL.handlerError(err) {
		return nil, err
	}
	Column, err := rows.Columns()
	if SQL.handlerError(err) {
		return nil, err
	}
	// 创建一个查询字段类型的slice
	values := make([]sql.RawBytes, len(Column))
	// 创建一个任意字段类型的slice
	scanArgs := make([]interface{}, len(values))
	// 创建一个slice保存所以的字段
	var allRows []map[string]string
	for i := range values {
		// 把values每个参数的地址存入scanArgs
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		// 把存放字段的元素批量放进去
		err = rows.Scan(scanArgs...)
		if SQL.handlerError(err) {
			return nil, err
		}
		tempRow := make(map[string]string, len(Column))
		for i, col := range values {
			var key = Column[i]
			tempRow[key] = string(col)
		}
		allRows = append(allRows, tempRow)
	}
	return allRows, nil
}

// ExecSQL ExecSQL
func (SQL *SQL) ExecSQL(queryString string) ([]map[string]string, error) {
	rows, err := SQL.db.Query(queryString)
	defer rows.Close()
	if SQL.handlerError(err) {
		return nil, err
	}
	Column, err := rows.Columns()
	if SQL.handlerError(err) {
		return nil, err
	}
	// 创建一个查询字段类型的slice
	values := make([]sql.RawBytes, len(Column))
	// 创建一个任意字段类型的slice
	scanArgs := make([]interface{}, len(values))
	// 创建一个slice保存所以的字段
	var allRows []map[string]string
	for i := range values {
		// 把values每个参数的地址存入scanArgs
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		// 把存放字段的元素批量放进去
		err = rows.Scan(scanArgs...)
		SQL.handlerError(err)
		tempRow := make(map[string]string, len(Column))
		for i, col := range values {
			var key = Column[i]
			tempRow[key] = string(col)
		}
		allRows = append(allRows, tempRow)
	}
	return allRows, nil
}

// QueryRow 查询单行
func (SQL *SQL) QueryRow() (map[string]string, error) {
	var queryString = SQL.ToString()
	result, err := SQL.db.Query(queryString)
	defer result.Close()
	if SQL.handlerError(err) {
		return nil, err
	}
	Column, err := result.Columns()
	// 创建一个查询字段类型的slice的键值对
	values := make([]sql.RawBytes, len(Column))
	// 创建一个任意字段类型的slice的键值对
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		// 把values每个参数的地址存入scanArgs
		scanArgs[i] = &values[i]
	}

	for result.Next() {
		err = result.Scan(scanArgs...)
		if SQL.handlerError(err) {
			return nil, err
		}
	}
	tempRow := make(map[string]string, len(Column))
	for i, col := range values {
		var key = Column[i]
		tempRow[key] = string(col)
	}
	return tempRow, nil
}

// RawTX 获取原始事务对象
func (SQL *SQL) RawTX() *sql.Tx {
	if SQL.tx != nil {
		return SQL.tx
	}
	return nil
}

// RawDB 获取原始 DB 对象
func (SQL *SQL) RawDB() *sql.DB {
	if SQL.db != nil {
		return SQL.db
	}
	return nil
}

// Begin 开启事务
// 默认只有写操作会进入事务、读取还是使用默认的
func (SQL *SQL) Begin() error {
	SQL.rollbackSign = true
	if SQL.tx == nil {
		tx, err := SQL.db.Begin()
		if err != nil {
			return err
		}
		SQL.tx = tx
		SQL.commitSign = beginStatus
		return nil
	}
	SQL.commitSign++
	return nil
}

// Rollback 回滚事务
// 会回滚本次的全部事务、范围不是最近的一次 Begin
// 而是本 SQL 第一次 Begin 后所有的操作
func (SQL *SQL) Rollback() error {
	if SQL.tx != nil && SQL.rollbackSign == true {
		err := SQL.tx.Rollback()
		if err != nil {
			return err
		}
		SQL.tx = nil
		return nil
	}
	return nil
}

// Commit 提交事务
// Begin Commit 必须一一对应
// 嵌套模型如下
// {开启 commitSign = 1
// 	{ 开启 commitSign = 2
// 		{ 开启 commitSign = 3
// 	 } commitSign = 3 提交
// 	} commitSign = 2 提交
// } commitSign = 1 提交
func (SQL *SQL) Commit() error {
	SQL.rollbackSign = false
	if SQL.tx != nil {
		if SQL.commitSign == beginStatus {
			err := SQL.tx.Commit()
			if err != nil {
				return err
			}
			SQL.tx = nil
			return nil
		} else {
			SQL.commitSign--
		}
		return nil
	}
	return nil
}

// username:password@tcp(localhost:3306)/database?charset=utf8mb4
// 事务支持参考 https://github.com/alberliu/session，https://www.jianshu.com/p/2a144332c3db
// 来自 https://github.com/tophubs/TopList/blob/master/Common/Db.go
// 感谢 tophubs
// Alain 整理
