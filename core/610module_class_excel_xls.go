package core

import (
	"github.com/shakinm/xlsReader/xls"
	"github.com/shakinm/xlsReader/xls/structure"
)

func (fns *InternalFunctionSet) Xls(fileName string) Value {
	workbook, err := xls.OpenFile(fileName)
	if err != nil {
		runtimeExcption(err)
	}
	obj := &XlsClass{workbook}
	return newClassExecutor("xls", obj, &obj)
}

type XlsClass struct {
	obj xls.Workbook
}

func (clazz *XlsClass) SheetNumber() {
	clazz.obj.GetNumberSheets()
}

func (clazz *XlsClass) Sheet(index int) Value {
	sheet, err := clazz.obj.GetSheet(index)
	if err != nil {
		runtimeExcption(err)
	}
	obj := &XlsSheetClass{sheet}
	return newClassExecutor("XlsSheet", obj, &obj)
}

type XlsSheetClass struct {
	obj *xls.Sheet
}

func (clazz *XlsSheetClass) Name() string {
	return clazz.obj.GetName()
}

func (clazz *XlsSheetClass) Row(index int) Value {
	obj := &XlsRowClass{clazz.obj, index}
	return newClassExecutor("XlsRow", obj, &obj)
}

func (clazz *XlsSheetClass) RowNumber() int {
	return clazz.obj.GetNumberRows()
}

func (clazz *XlsSheetClass) Rows() Value {
	var arr []Value
	size := clazz.obj.GetNumberRows()
	for index := 1; index < size; index++  {
		obj := &XlsRowClass{clazz.obj, index}
		arr = append(arr, newClassExecutor("XlsRow", obj, &obj))
	}
	return toJSONArray(arr)
}

type XlsRowClass struct {
	obj *xls.Sheet
	rowIndex int
}

func (clazz *XlsRowClass) Cols() Value {
	row, err := clazz.obj.GetRow(clazz.rowIndex)
	if err != nil {
		runtimeExcption(err)
	}

	var arr []Value
	for _, cell := range row.GetCols() {
		obj := &XlsCellClass{cell}
		arr = append(arr, newClassExecutor("XlsCell", obj, &obj))
	}
	return toJSONArray(arr)
}

func (clazz *XlsRowClass) Col(index int) Value {
	row, err := clazz.obj.GetRow(clazz.rowIndex)
	if err != nil {
		runtimeExcption(err)
	}
	cell, err := row.GetCol(index)
	if err != nil {
		runtimeExcption(err)
	}
	obj := &XlsCellClass{cell}
	return newClassExecutor("XlsCell", obj, &obj)
}

type XlsCellClass struct {
	obj structure.CellData
}

func (clazz *XlsCellClass) GetInt() int64 {
	return clazz.obj.GetInt64()
}

func (clazz *XlsCellClass) GetFloat() float64 {
	return clazz.obj.GetFloat64()
}

func (clazz *XlsCellClass) Str() string {
	return clazz.obj.GetString()
}

func (clazz *XlsCellClass) GetType() string {
	return clazz.obj.GetType()
}

func (clazz *XlsCellClass) GetXFIndex() int {
	return clazz.obj.GetXFIndex()
}