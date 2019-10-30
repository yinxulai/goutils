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

var globalDB *sql.DB
var config *Config

// SQL 结构化的 sql 语句
type SQL struct {
	conn        *sql.DB // 数据库连接
	fields      string  // 字段
	orderBy     string  // 排序条件
	tableName   string  // 表名
	execString  string  // 执行sql语句
	limitNumber string  // 限制条数
	whereString string  // where语句
	config      *Config
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
	sql.conn = globalDB
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

// Close Close
func (SQL *SQL) Close() error {
	err := SQL.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// Count 统计
func (SQL *SQL) Count(tableName string) *SQL {
	SQL.fields = "select count(*) as count from " + tableName
	SQL.tableName = tableName
	return SQL
}

// Select 查询方法
func (SQL *SQL) Select(tableName string, field []string) *SQL {
	var allField string
	allField = strings.Join(field, ",")
	SQL.fields = "select " + allField + " from " + tableName
	SQL.tableName = tableName
	return SQL
}

// Where Where
func (SQL *SQL) Where(cond map[string]string) *SQL {
	var whereString = ""
	if len(cond) != 0 {
		whereString = " where "
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
	SQL.limitNumber = " limit " + strconv.Itoa(number)
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

// Update Update
func (SQL SQL) Update(tableName string, str map[string]string) (int64, error) {
	var tempStr = ""
	var allValue []interface{}
	for key, value := range str {
		tempStr += key + "=" + "?" + ","
		allValue = append(allValue, value)
	}
	tempStr = strings.TrimSuffix(tempStr, ",")
	SQL.execString = "update " + tableName + " set " + tempStr
	var allStr = SQL.execString + SQL.whereString
	stmt, err := SQL.conn.Prepare(allStr)
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
func (SQL SQL) Delete(tableName string) (int64, error) {
	var tempStr = ""
	tempStr = "delete from " + tableName + SQL.whereString
	stmt, err := SQL.conn.Prepare(tempStr)
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
func (SQL SQL) Insert(tableName string, data map[string]string) (int64, error) {
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
	var theStr = "insert into " + tableName + " " + allField + " values " + allValue
	stmt, err := SQL.conn.Prepare(theStr)
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
func (SQL SQL) Pagination(Page int, Limit int) (map[string]interface{}, error) {
	res, err := GetConn().Select(SQL.tableName, []string{"count(*) as count"}).QueryRow()
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
	queryString := SQL.fields + SQL.whereString + SQL.orderBy + " limit " + strconv.Itoa(setOff) + "," + strconv.Itoa(Limit)
	rows, err := SQL.conn.Query(queryString)
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
func (SQL SQL) QueryAll() ([]map[string]string, error) {
	var queryString = SQL.fields + SQL.whereString + SQL.orderBy + SQL.limitNumber
	rows, err := SQL.conn.Query(queryString)
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
func (SQL SQL) ExecSQL(queryString string) ([]map[string]string, error) {
	rows, err := SQL.conn.Query(queryString)
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
func (SQL SQL) QueryRow() (map[string]string, error) {
	var queryString = SQL.fields + SQL.whereString + SQL.orderBy + SQL.limitNumber
	result, err := SQL.conn.Query(queryString)
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

// 来自 https://github.com/tophubs/TopList/blob/master/Common/Db.go
// 感谢 tophubs
// Alain 整理
