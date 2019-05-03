package main

import (
	"excelFilter/util"
	"github.com/lxn/walk"
	"log"
	"strconv"
)

// MyTableView represents
type MyTableView struct {
	*walk.TableView
}

// TableModel represents
type TableModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      [][]string
	// source is the original items (preserved for filter)
	source   []string
}

// RowCount Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *TableModel) RowCount() int {
	return len(m.items)
}

// Value Called by the TableView when it needs the text to display for a given cell.
func (m *TableModel) Value(row, col int) interface{} {
	return m.items[row][col]
}

func (tv *MyTableView) addCol(title string, isAlignLeft bool) {
	col := walk.NewTableViewColumn()
	col.SetTitle(title)
	col.SetWidth(100)

	if isAlignLeft == true{
		col.SetAlignment(walk.AlignNear)
	} else {
		col.SetAlignment(walk.AlignFar)
	}

	tv.TableView.Columns().Add(col)
}

func (tv *MyTableView) handleFile(m *TableModel, file string) {
	csvData, err := util.ReadCsvData(file)
	if err != nil {
		log.Fatalln(err)
	}

	tv.TableView.Columns().Clear()
	for i := 0; i < len(csvData[0]); i++ {
		// long content will be aligned left
		isAlignLeft := len(csvData[0][i]) > 7
		tv.addCol("#" + strconv.Itoa(i + 1), isAlignLeft)
	}

	m.items = csvData
	m.source = util.TwoDimensionTo1DArray(csvData, util.SeperatorStr)
	m.PublishRowsReset()
}
