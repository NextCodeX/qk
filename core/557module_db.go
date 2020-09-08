package core

import (
	dbManager "database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func (mr *ModuleRegister) DBModuleInit() {
	dsc := &DateSourceConstructor{}
	fs := collectFunctionInfo(&dsc)
	functionRegister("", fs)
}


// 数据库连接对象
type DataSource struct {
	driver string
	source string
}

type DateSourceConstructor struct{}

func (dsc *DateSourceConstructor) ConnDB(driverName, sourceName string) *ClassExecutor {
	ds := &DataSource{driverName, sourceName}
	return newClassExecutor("date", ds, &ds)
}

func (ds *DataSource) Exec(args []interface{}) int64 {
	sql, vals := ds.parseArgs("update", args)
	db, err := dbManager.Open(ds.driver, ds.source)
	assert(err != nil, "failed to build db connection:", err)
	defer db.Close()

	execResult, err := db.Exec(sql, vals)
	assert(err != nil, "failed to execute sql:", err)

	affected, err := execResult.RowsAffected()
	assert(err != nil, "failed to update:", err)
	return affected
}

func (ds *DataSource) Insert(args []interface{}) interface{} {
	sql, vals := ds.parseArgs("insert", args)
	db, err := dbManager.Open(ds.driver, ds.source)
	assert(err != nil, "failed to build db connection:", err)
	defer db.Close()

	stmt, err := db.Prepare(sql)
	assert(err != nil, "failed to build db connection:", err)
	defer stmt.Close()

	execResult, err := stmt.Exec(vals...)
	assert(err != nil, "failed to insert:", err)

	id, err := execResult.LastInsertId()
	assert(err != nil, "failed to insert:", err)
	return id
}

func (ds *DataSource) Update(args []interface{}) int64 {
	sql, vals := ds.parseArgs("update", args)
	db, err := dbManager.Open(ds.driver, ds.source)
	assert(err != nil, "failed to build db connection:", err)
	defer db.Close()

	stmt, err := db.Prepare(sql)
	assert(err != nil, "failed to build db connection:", err)
	defer stmt.Close()

	execResult, err := stmt.Exec(vals...)
	assert(err != nil, "failed to update:", err)

	affected, err := execResult.RowsAffected()
	assert(err != nil, "failed to update:", err)
	return affected
}

func (ds *DataSource) GetValue(args []interface{}) interface{} {
	sql, vals := ds.parseArgs("getValue", args)
	db, err := dbManager.Open(ds.driver, ds.source)
	assert(err != nil, "failed to build db connection:", err)
	defer db.Close()

	rows, err := db.Query(sql, vals...)
	assert(err != nil, "failed to query data:", err)
	defer rows.Close()


	colTypes, _ := rows.ColumnTypes()
	colCount := len(colTypes)
	assert(colCount>1, "sql query return over more field value")

	var value interface{}
	rowCount := 1
	for rows.Next() {
		assert(rowCount>1, "sql query return over more row data")

		err = rows.Scan(&value)
		assert(err != nil, "failed to extract field value:", err)


		value = ds.getFieldValue(value)

		rowCount++
	}
	return value
}

func (ds *DataSource) GetRow(args []interface{}) map[string]interface{} {
	sql, vals := ds.parseArgs("getRow", args)
	db, err := dbManager.Open(ds.driver, ds.source)
	assert(err != nil, "failed to build db connection:", err)
	defer db.Close()

	rows, err := db.Query(sql, vals...)
	assert(err != nil, "failed to query data:", err)
	defer rows.Close()

	colTypes, _ := rows.ColumnTypes()
	colCount := len(colTypes)
	var list []map[string]interface{}
	rowCount := 1
	for rows.Next() {
		assert(rowCount>1, "sql query return over more row data")

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
	if len(list) < 1 {
		return nil
	} else {
		return list[0]
	}
}

func (ds *DataSource) GetRows(args []interface{}) []interface{} {
	sql, vals := ds.parseArgs("getRows", args)
	db, err := dbManager.Open(ds.driver, ds.source)
	assert(err != nil, "failed to build db connection:", err)
	defer db.Close()

	rows, err := db.Query(sql, vals...)
	assert(err != nil, "failed to query data:", err)
	defer rows.Close()

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
	for i:=0; i<size; i++ {
		var container interface{}
		valueContainers = append(valueContainers, &container)
	}
	return valueContainers
}
