package core

import (
	dbManager "database/sql"
	"strconv"
	"time"
)

type DBValueConvertor func(tmp interface{}, colType *dbManager.ColumnType) (interface{}, bool)

var dbValConvertorList []DBValueConvertor

func (ds *DataSource) initValConvertors() {
	if dbValConvertorList != nil {
		return
	}
	dbValConvertorList = append(dbValConvertorList, ds.pgTypeMatch)
	dbValConvertorList = append(dbValConvertorList, ds.toFloat)
}

func (ds *DataSource) getFieldValue(tmp interface{}, colType *dbManager.ColumnType) interface{} {
	ds.initValConvertors()

	var res interface{}
	var ok bool
	for _, convertor := range dbValConvertorList {
		if res, ok = convertor(tmp, colType); ok {
			return res
		}
	}

	if tt, ok := tmp.(time.Time); ok {
		return tt.Format(CommonDatetimeFormat)
	}
	return tmp
}

func (ds *DataSource) pgTypeMatch(tmp interface{}, colType *dbManager.ColumnType) (interface{}, bool) {
	if ds.driver != dbPostgres {
		return nil, false
	}

	dbtype := colType.DatabaseTypeName()
	if bs, ok := tmp.([]byte); ok {
		switch dbtype {
		case "NAME":
			return string(bs), true
		case "NUMERIC":
			if res, err := strconv.ParseFloat(string(bs), 64); err == nil {
				return res, true
			}
		}
	}
	return nil, false
}

func (ds *DataSource) toFloat(tmp interface{}, colType *dbManager.ColumnType) (interface{}, bool) {
	dbtype := colType.DatabaseTypeName()
	if bs, ok := tmp.([]byte); ok && dbtype == "DECIMAL" {
		if res, err := strconv.ParseFloat(string(bs), 64); err == nil {
			return res, true
		}
	}
	return nil, false
}
