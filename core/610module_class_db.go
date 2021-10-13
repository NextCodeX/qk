package core

import (
	dbManager "database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb" // sql server
	_ "github.com/go-sql-driver/mysql"
	"io"
	"strings"

	_ "github.com/lib/pq"          // postgreSQL
	_ "github.com/sijms/go-ora/v2" // oracle
	_ "modernc.org/sqlite"
)

func (fns *InternalFunctionSet) Sqlserver(username, password, url string) Value {
	seperatorIndex := strings.Index(url, "/")
	var host, port, dbName string
	if seperatorIndex < 0 || seperatorIndex == len(url)-1 {
		runtimeExcption("error: database name is empty!")
	} else {
		netAddress := url[:seperatorIndex]
		colonIndex := strings.Index(netAddress, ":")
		if colonIndex < 0 || colonIndex == len(netAddress)-1 {
			runtimeExcption("host address format is error!")
		} else {
			host = netAddress[:colonIndex]
			port = netAddress[colonIndex+1:]
		}
		dbName = url[seperatorIndex+1:]
	}
	// e.g. "server=192.168.1.103;port=1433;database=STG;user id=SA;password=root@123"
	sourceName := fmt.Sprintf("server=%v;port=%v;database=%v;user id=%v;password=%v", host, port, dbName, username, password)
	return connDB("mysql", sourceName)
}

func (fns *InternalFunctionSet) Oracle(username, password, url string) Value {
	seperatorIndex := strings.Index(url, "/")
	var netAddress, uri string
	if seperatorIndex < 0 {
		netAddress = url
		uri = "?charset=utf8"
	} else {
		netAddress = url[:seperatorIndex]
		uri = url[seperatorIndex:]
	}
	// e.g. oracle://user:pass@server/service_name
	sourceName := fmt.Sprintf("oracle://%v:%v@%v/%v", username, password, netAddress, uri)
	return connDB("oracle", sourceName)
}

func (fns *InternalFunctionSet) Mysql(username, password, url string) Value {
	seperatorIndex := strings.Index(url, "/")
	var netAddress, uri string
	if seperatorIndex < 0 {
		netAddress = url
		uri = "?charset=utf8"
	} else {
		netAddress = url[:seperatorIndex]
		uri = url[seperatorIndex:]
	}
	// e.g. root:root@tcp(192.168.1.103:3306)/tx?charset=utf8
	sourceName := fmt.Sprintf("%v:%v@tcp(%v)%v", username, password, netAddress, uri)
	return connDB("mysql", sourceName)
}

func (fns *InternalFunctionSet) Sqlite(dbName string) Value {
	return connDB("sqlite", dbName)
}

func connDB(driverName, sourceName string) Value {
	db, err := dbManager.Open(driverName, sourceName)
	assert(err != nil, "failed to build db connection:", err)

	ds := &DataSource{driverName,
		sourceName,
		db}
	return newClass("db", &ds)
}

// 数据库连接对象
type DataSource struct {
	driver string
	source string
	db     *dbManager.DB
}

func (ds *DataSource) Exec(args []interface{}) int64 {
	sql, vals := ds.parseArgs("update", args)

	execResult, err := ds.db.Exec(sql, vals...)
	assert(err != nil, "failed to execute sql:", sql, err)

	affected, err := execResult.RowsAffected()
	assert(err != nil, "failed to update:", sql, err)
	return affected
}

func (ds *DataSource) Insert(args []interface{}) interface{} {
	sql, vals := ds.parseArgs("insert", args)

	stmt, err := ds.db.Prepare(sql)
	assert(err != nil, "failed to prepare:", err)
	defer cl(stmt)

	execResult, err := stmt.Exec(vals...)
	assert(err != nil, "failed to insert:", err)

	id, err := execResult.LastInsertId()
	assert(err != nil, "failed to insert:", err)
	return id
}

func (ds *DataSource) Update(args []interface{}) int64 {
	sql, vals := ds.parseArgs("update", args)

	stmt, err := ds.db.Prepare(sql)
	assert(err != nil, "failed to prepare:", err)
	defer cl(stmt)

	execResult, err := stmt.Exec(vals...)
	assert(err != nil, "failed to update:", err)

	affected, err := execResult.RowsAffected()
	assert(err != nil, "failed to update:", err)
	return affected
}

func (ds *DataSource) GetValue(args []interface{}) interface{} {
	sql, vals := ds.parseArgs("getValue", args)

	rows, err := ds.db.Query(sql, vals...)
	assert(err != nil, "failed to query data:", err)
	defer cl(rows)

	colTypes, _ := rows.ColumnTypes()
	colCount := len(colTypes)
	assert(colCount > 1, "sql query return over more field value")

	var value interface{}
	rowCount := 1
	for rows.Next() {
		assert(rowCount > 1, "sql query return over more row data")

		err = rows.Scan(&value)
		assert(err != nil, "failed to extract field value:", err)

		value = ds.getFieldValue(value)

		rowCount++
	}
	return value
}
func (ds *DataSource) Val(args []interface{}) interface{} {
	return ds.GetValue(args)
}

func (ds *DataSource) GetRow(args []interface{}) map[string]interface{} {
	sql, vals := ds.parseArgs("getRow", args)

	rows, err := ds.db.Query(sql, vals...)
	assert(err != nil, "failed to query data:", err)
	defer cl(rows)

	colTypes, _ := rows.ColumnTypes()
	colCount := len(colTypes)
	var list []map[string]interface{}
	rowCount := 1
	for rows.Next() {
		assert(rowCount > 1, "sql query return over more row data")

		valueContainers := getValueContainers(colCount)
		err = rows.Scan(valueContainers...)
		assert(err != nil, "failed to extract field value:", err)

		row := make(map[string]interface{})
		for i, colType := range colTypes {
			colName := colType.Name()
			tmp := *valueContainers[i].(*interface{})

			row[colName] = ds.getFieldValue(tmp)
		}
		list = append(list, row)
		rowCount++
	}
	if len(list) > 0 {
		return list[0]
	} else {
		return nil
	}
}
func (ds *DataSource) Row(args []interface{}) map[string]interface{} {
	return ds.GetRow(args)
}

func (ds *DataSource) GetRows(args []interface{}) []interface{} {
	sql, vals := ds.parseArgs("getRows", args)

	rows, err := ds.db.Query(sql, vals...)
	assert(err != nil, "failed to query data:", err)
	defer cl(rows)

	colTypes, _ := rows.ColumnTypes()
	colCount := len(colTypes)
	var list []interface{}
	for rows.Next() {
		valueContainers := getValueContainers(colCount)
		err = rows.Scan(valueContainers...)
		assert(err != nil, "failed to extract field value:", err)

		row := make(map[string]interface{})
		for i, colType := range colTypes {
			colName := colType.Name()
			tmp := *valueContainers[i].(*interface{})

			row[colName] = ds.getFieldValue(tmp)
		}
		list = append(list, row)
	}
	return list
}
func (ds *DataSource) Rows(args []interface{}) []interface{} {
	return ds.GetRows(args)
}

func (ds *DataSource) parseArgs(methodName string, args []interface{}) (sql string, values []interface{}) {
	assert(len(args) < 1, fmt.Sprintf("method db.%v() must has one parameters.", methodName))
	sql, ok := args[0].(string)
	assert(!ok, fmt.Sprintf("method db.%v() the first parameter must be string type.", methodName))
	return sql, args[1:]
}

func (ds *DataSource) getFieldValue(tmp interface{}) interface{} {
	bs, ok := tmp.([]byte)
	if ok {
		return string(bs)
	} else {
		return tmp
	}
}

func getValueContainers(size int) []interface{} {
	var valueContainers []interface{}
	for i := 0; i < size; i++ {
		var container interface{}
		valueContainers = append(valueContainers, &container)
	}
	return valueContainers
}

// 释放资源
func cl(obj io.Closer) {
	err := obj.Close()
	if err != nil {
		runtimeExcption(err)
	}
}
