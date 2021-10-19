package core

import (
	"strings"
)
import dbManager "database/sql"

// 获取表的字段列表。
func (ds *DataSource) Fields(obj string) []string {
	if ds.driver == dbSqlite {
		return ds.sqliteFields(obj)
	} else if ds.driver == dbPostgres {
		return ds.postgresFields(obj)
	} else if ds.driver == dbSqlServer {
		return ds.sqlServerFields(obj)
	} else if ds.driver == dbOracle {
		return ds.oracleFields(obj)
	} else if ds.driver == dbMysql {
		return ds.mysqlFields(obj)
	} else {
		return nil
	}
}
func (ds *DataSource) Fs(obj string) []string {
	return ds.Fields(obj)
}

func (ds *DataSource) postgresFields(obj string) []string {
	sql := "select column_name from information_schema.columns where table_schema='public' and table_name=$1"
	rows, err := ds.db.Query(sql, obj)
	assert(err != nil, "failed to query data:", err)
	defer cl(rows)

	var res []string
	for rows.Next() {
		var fieldName string
		err := rows.Scan(&fieldName)
		assert(err != nil, "failed to extract field name:", err)

		res = append(res, fieldName)
	}
	return res
}

func (ds *DataSource) sqlServerFields(obj string) []string {
	sql := "select name from syscolumns where id=(select max(id) from sysobjects where xtype='u' and name=@p1)"
	rows, err := ds.db.Query(sql, obj)
	assert(err != nil, "failed to query data:", err)
	defer cl(rows)

	var res []string
	for rows.Next() {
		var fieldName string
		err := rows.Scan(&fieldName)
		assert(err != nil, "failed to extract field name:", err)

		res = append(res, fieldName)
	}
	return res
}

func (ds *DataSource) oracleFields(obj string) []string {
	obj = strings.ToUpper(obj)
	sql := "select t.column_name from user_col_comments t where t.table_name = :1"
	rows, err := ds.db.Query(sql, obj)
	assert(err != nil, "failed to query data:", err)
	defer cl(rows)

	var res []string
	for rows.Next() {
		var fieldName string
		err := rows.Scan(&fieldName)
		assert(err != nil, "failed to extract field name:", err)

		res = append(res, fieldName)
	}
	return res
}

func (ds *DataSource) mysqlFields(obj string) []string {
	var rows *dbManager.Rows
	var err error
	if strings.Contains(obj, ".") {
		subs := strings.Split(obj, ".")
		dbName, tabName := subs[0], subs[1]
		sql := "select column_name from information_schema.columns where table_name=? and table_schema=?"
		rows, err = ds.db.Query(sql, tabName, dbName)

	} else {
		sql := "select column_name from information_schema.columns where table_name=?"
		rows, err = ds.db.Query(sql, obj)
	}
	assert(err != nil, "failed to query data:", err)
	defer cl(rows)

	var res []string
	for rows.Next() {
		var fieldName string
		err := rows.Scan(&fieldName)
		assert(err != nil, "failed to extract field name:", err)

		res = append(res, fieldName)
	}
	return res
}

func (ds *DataSource) sqliteFields(obj string) []string {
	sql := "select sql from sqlite_schema where name=?"

	rows, err := ds.db.Query(sql, obj)
	assert(err != nil, "failed to query data:", err)
	defer cl(rows)

	var tabDefinition string
	rowCount := 1
	for rows.Next() {
		if rowCount > 1 {
			break
		}

		err = rows.Scan(&tabDefinition)
		assert(err != nil, "failed to get table definition:", err)

		rowCount++
	}

	startIndex := strings.Index(tabDefinition, "(") + 1
	endIndex := strings.LastIndex(tabDefinition, ")")
	fstr := tabDefinition[startIndex:endIndex]
	cols := strings.Split(fstr, ",")
	var res []string
	for _, col := range cols {
		col = strings.TrimSpace(col)
		subStr := strings.Split(col, " ")
		res = append(res, subStr[0])
	}
	return res
}
