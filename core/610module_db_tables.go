package core

import (
	"fmt"
	"strings"
)

// 判断表在当前数据库下是否存在
func (ds *DataSource) Exists(obj string) bool {
	var tabs []string
	if ds.driver == dbSqlite {
		tabs = ds.sqliteTables()
	} else if ds.driver == dbPostgres {
		tabs = ds.postgresTables()
	} else if ds.driver == dbSqlServer {
		tabs = ds.sqlServerTables()
	} else if ds.driver == dbOracle {
		tabs = ds.oracleTables()
	} else if ds.driver == dbMysql {
		tabs = ds.mysqlTables()
	} else {
		fmt.Printf("exists(): db %v is not supported!\n", ds.driver)
		return false
	}
	if len(tabs) < 1 {
		return false
	}
	obj = strings.TrimSpace(obj)
	for _, tab := range tabs {
		if tab == obj || tab == strings.ToUpper(obj) {
			return true
		}
	}
	return false
}

// 查询当前数据库下的所有表名称
func (ds *DataSource) Tables() []string {
	if ds.driver == dbSqlite {
		return ds.sqliteTables()
	} else if ds.driver == dbPostgres {
		return ds.postgresTables()
	} else if ds.driver == dbSqlServer {
		return ds.sqlServerTables()
	} else if ds.driver == dbOracle {
		return ds.oracleTables()
	} else if ds.driver == dbMysql {
		return ds.mysqlTables()
	} else {
		fmt.Printf("tables(): db %v is not supported!\n", ds.driver)
		return nil
	}
}

func (ds *DataSource) postgresTables() []string {
	sql := "select tablename from pg_tables where schemaname=$1"
	rows, err := ds.db.Query(sql, "public")
	assert(err != nil, "failed to get table name list:", err)
	defer cl(rows)

	var res []string
	for rows.Next() {
		var tabName string
		err = rows.Scan(&tabName)
		assert(err != nil, "failed to extract table name:", err)
		res = append(res, tabName)
	}
	return res
}

func (ds *DataSource) sqlServerTables() []string {
	sql := "select name from sys.tables"
	rows, err := ds.db.Query(sql)
	assert(err != nil, "failed to get table name list:", err)
	defer cl(rows)

	var res []string
	for rows.Next() {
		var tabName string
		err = rows.Scan(&tabName)
		assert(err != nil, "failed to extract table name:", err)
		res = append(res, tabName)
	}
	return res
}

func (ds *DataSource) oracleTables() []string {
	sql := "select table_name from user_tables where TABLE_NAME NOT LIKE 'MVIEW%'  AND TABLE_NAME NOT LIKE 'DEF$%' AND TABLE_NAME NOT LIKE 'REPCAT$%' AND TABLE_NAME NOT LIKE 'OL$%' AND TABLE_NAME NOT LIKE 'REDO_%' AND TABLE_NAME NOT LIKE 'ROLLING$%' AND TABLE_NAME NOT LIKE 'AQ$%' AND TABLE_NAME NOT LIKE 'LOGMNR%' AND TABLE_NAME NOT LIKE 'SCHEDULER%' AND TABLE_NAME NOT LIKE 'LOGSTDBY$%' AND TABLE_NAME!='SQLPLUS_PRODUCT_PROFILE' AND TABLE_NAME!='HELP'"
	rows, err := ds.db.Query(sql)
	assert(err != nil, "failed to get table name list:", err)
	defer cl(rows)

	var res []string
	for rows.Next() {
		var tabName string
		err = rows.Scan(&tabName)
		assert(err != nil, "failed to extract table name:", err)
		res = append(res, tabName)
	}
	return res
}

func (ds *DataSource) mysqlTables() []string {
	sql := "show tables"
	rows, err := ds.db.Query(sql)
	assert(err != nil, "failed to get table name list:", err)
	defer cl(rows)

	var res []string
	for rows.Next() {
		var tabName string
		err = rows.Scan(&tabName)
		assert(err != nil, "failed to extract table name:", err)
		res = append(res, tabName)
	}
	return res
}

func (ds *DataSource) sqliteTables() []string {
	sql := "select name from sqlite_schema where type = 'table'"
	rows, err := ds.db.Query(sql)
	assert(err != nil, "failed to get table name list:", err)
	defer cl(rows)

	var res []string
	for rows.Next() {
		var tabName string
		err = rows.Scan(&tabName)
		assert(err != nil, "failed to extract table name:", err)
		res = append(res, tabName)
	}
	return res
}
