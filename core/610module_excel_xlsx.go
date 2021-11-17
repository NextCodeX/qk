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
	obj := &Xlsx{f}
	return newClass("Xlsx", &obj)
}

func (fns *InternalFunctionSet) NewXlsx() Value {
	f := excelize.NewFile()
	obj := &Xlsx{f}
	return newClass("Xlsx", &obj)
}

type Xlsx struct {
	obj *excelize.File
}

// Create a new sheet.
func (clazz *Xlsx) NewSheet(name string) int {
	return clazz.NewSheet(name)
}

func (clazz *Xlsx) SetActiveSheet(index int) {
	clazz.obj.SetActiveSheet(index)
}

func (clazz *Xlsx) SaveAs(path string) {
	err := clazz.obj.SaveAs(path)
	if err != nil {
		runtimeExcption(err)
	}
}

func (clazz *Xlsx) Save() {
	err := clazz.obj.Save()
	if err != nil {
		runtimeExcption(err)
	}
}

func (clazz *Xlsx) Sheet(index int) Value {
	sheetName := clazz.obj.GetSheetName(index)
	obj := &XlsxSheet{clazz.obj, index, sheetName}
	return newClass("XlsxSheet", &obj)
}

type XlsxSheet struct {
	obj        *excelize.File
	sheetIndex int
	sheetName  string
}

var colIds = []int32{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func (clazz *XlsxSheet) SetData(headers JSONArray, data JSONArray) {
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
func (clazz *XlsxSheet) SetCellValue(axis string, value string) {
	clazz.obj.SetCellValue(clazz.sheetName, axis, value)
}

func (clazz *XlsxSheet) SetValue(axisX, axisY int, value string) {
	x := colIds[axisX]
	y := axisY + 1
	axis := fmt.Sprintf("%v%v", x, y)
	clazz.obj.SetCellValue(clazz.sheetName, axis, value)
}

func (clazz *XlsxSheet) Rows() Value {
	rows, err := clazz.obj.Rows(clazz.sheetName)
	if err != nil {
		runtimeExcption(err)
	}
	obj := &XlsxRows{rows}
	return newClass("XlsxRows", &obj)
}

func (clazz *XlsxSheet) CellVals() Value {
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

func (clazz *XlsxSheet) CellVal(axis string) string {
	val, err := clazz.obj.GetCellValue(clazz.sheetName, axis)
	if err != nil {
		runtimeExcption(err)
	}
	return val
}

func (clazz *XlsxSheet) Val(axisX, axisY int) string {
	x := colIds[axisX]
	y := axisY + 1
	axis := fmt.Sprintf("%v%v", x, y)
	val, err := clazz.obj.GetCellValue(clazz.sheetName, axis)
	if err != nil {
		runtimeExcption(err)
	}
	return val
}

type XlsxRows struct {
	obj *excelize.Rows
}

func (clazz *XlsxRows) Next() bool {
	return clazz.obj.Next()
}

func (clazz *XlsxRows) Cols() Value {
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
