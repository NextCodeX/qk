package core

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func (fns *InternalFunctionSet) Xlsx(fileName string) Value {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		runtimeExcption(f)
	}
	obj := &XlsxClass{f}
	return newClass("Xlsx", &obj)
}

func (fns *InternalFunctionSet) NewXlsx() Value {
	f := excelize.NewFile()
	obj := &XlsxClass{f}
	return newClass("Xlsx", &obj)
}

type XlsxClass struct {
	obj *excelize.File
}

// Create a new sheet.
func (clazz *XlsxClass) NewSheet(name string) int {
	return clazz.NewSheet(name)
}

func (clazz *XlsxClass) SetActiveSheet(index int) {
	clazz.obj.SetActiveSheet(index)
}

func (clazz *XlsxClass) SaveAs(path string) {
	err := clazz.obj.SaveAs(path)
	if err != nil {
		runtimeExcption(err)
	}
}

func (clazz *XlsxClass) Save() {
	err := clazz.obj.Save()
	if err != nil {
		runtimeExcption(err)
	}
}

func (clazz *XlsxClass) Sheet(index int) Value {
	sheetName := clazz.obj.GetSheetName(index)
	obj := &XlsxSheetClass{clazz.obj, index, sheetName}
	return newClass("XlsxSheet", &obj)
}


type XlsxSheetClass struct {
	obj *excelize.File
	sheetIndex int
	sheetName string
}

var colIds = []int32{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}


func (clazz *XlsxSheetClass) SetData(headers JSONArray, data JSONArray) {
	for i, item := range headers.values() {
		header := goStr(item)
		axis := string(colIds[i]) + "1"
		clazz.obj.SetCellValue(clazz.sheetName, axis, header)
	}
	for i, item := range data.values() {
		rowIndex := i + 2
		row := goObj(item)
		for index, key := range headers.values() {
			header := goStr(key)
			axis := fmt.Sprintf(`%v%v`, string(colIds[index]), rowIndex)
			clazz.obj.SetCellValue(clazz.sheetName, axis, row.get(header).String())
		}
	}
}

// Set value of a cell.
func (clazz *XlsxSheetClass) SetCellValue(axis string, value string) {
	clazz.obj.SetCellValue(clazz.sheetName, axis, value)
}

func (clazz *XlsxSheetClass) Rows() Value {
	rows, err := clazz.obj.Rows(clazz.sheetName)
	if err != nil {
		runtimeExcption(err)
	}
	obj := &XlsxRowsClass{rows}
	return newClass("XlsxRows", &obj)
}

func (clazz *XlsxSheetClass) CellVals() Value {
	rows, err := clazz.obj.GetRows(clazz.sheetName)
	if err != nil {
		runtimeExcption(err)
	}
	var arr []Value
	for _, row := range rows {
		var subArr []Value
		for _, cell := range row {
			subArr = append(subArr, newQKValue(cell))
		}
		arr = append(arr, array(subArr))
	}
	return array(arr)
}

func (clazz *XlsxSheetClass) CellVal(axis string) string {
	val, err := clazz.obj.GetCellValue(clazz.sheetName, axis)
	if err != nil {
		runtimeExcption(err)
	}
	return val
}

type XlsxRowsClass struct {
	obj *excelize.Rows
}

func (clazz *XlsxRowsClass) Next() bool {
	return clazz.obj.Next()
}

func (clazz *XlsxRowsClass) Columns() Value {
	cols, err := clazz.obj.Columns()
	if err != nil {
		runtimeExcption(err)
	}
	var arr []Value
	for _, val := range cols {
		arr = append(arr, newQKValue(val))
	}
	return array(arr)
}
