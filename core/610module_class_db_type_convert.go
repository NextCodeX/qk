package core

import (
	dbManager "database/sql"
	"fmt"
	"strconv"
	"time"
)

func (ds *DataSource) getFieldValue(tmp interface{}, colType *dbManager.ColumnType) interface{} {
	if ds.driver == dbPostgres {
		return ds.pgTypeMatch(tmp, colType)
	}

	//fmt.Printf("column %v type: %v, %v\n", colType.Name(), colType.ScanType(), colType.DatabaseTypeName())
	if bs, ok := tmp.([]byte); ok && isFloatType(colType.DatabaseTypeName()) {
		//fmt.Println("bytes ->", string(bs), len(bs))
		res, err := strconv.ParseFloat(string(bs), 64)
		if err != nil {
			fmt.Printf("column %v type: %v, %v\n -> failed to float", colType.Name(), colType.ScanType(), colType.DatabaseTypeName())
			return nil
		}
		return res
	} else if tt, ok := tmp.(time.Time); ok {
		return tt.Format(CommonDatetimeFormat)
	} else {
		return tmp
	}
}

func (ds *DataSource) pgTypeMatch(tmp interface{}, colType *dbManager.ColumnType) interface{} {
	//fmt.Printf("column %v type: %v, %v\n", colType.Name(), colType.ScanType(), colType.DatabaseTypeName())
	dbtype := colType.DatabaseTypeName()
	if bs, ok := tmp.([]byte); ok {
		if dbtype == "NAME" {
			return string(bs)
		} else if dbtype == "NUMERIC" {
			res, err := strconv.ParseFloat(string(bs), 64)
			if err != nil {
				fmt.Printf("column %v type: %v, %v\n -> failed to float", colType.Name(), colType.ScanType(), colType.DatabaseTypeName())
				return nil
			}
			return res
		}
	}
	return tmp
}

func isFloatType(dbType string) bool {
	switch dbType {
	case "DECIMAL":
		return true
	}
	return false
}
