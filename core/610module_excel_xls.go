package core

import (
	"github.com/shakinm/xlsReader/xls"
	"github.com/shakinm/xlsReader/xls/structure"
)

func (this *InternalFunctionSet) Xls(fileName string) Value {
	workbook, err := xls.OpenFile(fileName)
	if err != nil {
		runtimeExcption(err)
	}
	obj := &Xls{workbook}
	return newClass("Xls", &obj)
}

type Xls struct {
	obj xls.Workbook
}

func (clazz *Xls) SheetNumber() int {
	return clazz.obj.GetNumberSheets()
}

func (clazz *Xls) Sheet(index int) Value {
	sheet, err := clazz.obj.GetSheet(index)
	if err != nil {
		runtimeExcption(err)
	}
	obj := &XlsSheet{sheet}
	return newClass("XlsSheet", &obj)
}

type XlsSheet struct {
	obj *xls.Sheet
}

func (clazz *XlsSheet) Name() string {
	return clazz.obj.GetName()
}

func (clazz *XlsSheet) Row(index int) Value {
	obj := &XlsRow{clazz.obj, index}
	return newClass("XlsRow", &obj)
}

func (clazz *XlsSheet) RowNumber() int {
	return clazz.obj.GetNumberRows()
}

func (clazz *XlsSheet) Rows() Value {
	var arr []Value
	size := clazz.obj.GetNumberRows()
	for index := 1; index < size; index++ {
		obj := &XlsRow{clazz.obj, index}
		arr = append(arr, newClass("XlsRow", &obj))
	}
	return array(arr)
}

type XlsRow struct {
	obj      *xls.Sheet
	rowIndex int
}

func (clazz *XlsRow) Cols() Value {
	row, err := clazz.obj.GetRow(clazz.rowIndex)
	if err != nil {
		runtimeExcption(err)
	}

	var arr []Value
	for _, cell := range row.GetCols() {
		obj := &XlsCell{cell}
		arr = append(arr, newClass("XlsCell", &obj))
	}
	return array(arr)
}

func (clazz *XlsRow) Col(index int) Value {
	row, err := clazz.obj.GetRow(clazz.rowIndex)
	if err != nil {
		runtimeExcption(err)
	}
	cell, err := row.GetCol(index)
	if err != nil {
		runtimeExcption(err)
	}
	obj := &XlsCell{cell}
	return newClass("XlsCell", &obj)
}

type XlsCell struct {
	obj structure.CellData
}

func (clazz *XlsCell) GetInt() int64 {
	return clazz.obj.GetInt64()
}

func (clazz *XlsCell) GetFloat() float64 {
	return clazz.obj.GetFloat64()
}

func (clazz *XlsCell) Str() string {
	return clazz.obj.GetString()
}

func (clazz *XlsCell) GetType() string {
	return clazz.obj.GetType()
}

func (clazz *XlsCell) GetXFIndex() int {
	return clazz.obj.GetXFIndex()
}
