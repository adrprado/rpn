package reports

import (
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/pkg/errors"
)

// Excel instance reachable data
type Excel struct {
	xlsx *excelize.File
}

//
// newExcel creates a new Excel instance
//
func newExcel() (e *Excel) {
	e = &Excel{}
	e.xlsx = excelize.NewFile()
	return
}

//
// saveAndCloseExcel saves to filename (need to set the directory as well)
//
func (e *Excel) saveAndCloseExcel(filename string) (err error) {
	// newFilename = time.Now().Format("02Jan06_150405.000") + ".xlsx" // DDMMMYY
	e.xlsx.DeleteSheet("Sheet1")
	e.xlsx.SetActiveSheet(1)
	err = e.xlsx.SaveAs(filename)
	if err != nil {
		return errors.Wrapf(err, "erro ao salvar planilha %s", filename)
	}
	return
}

// Sheet struct
type Sheet struct {
	e       *Excel
	name    string
	currRow int
}

func (e *Excel) newSheet(name string) (s *Sheet, err error) {
	s = &Sheet{}
	s.name = name
	s.e = e
	s.currRow = 1

	// Create a new sheet.
	// Avoid duplicated sheet
	if index := e.xlsx.GetSheetIndex(name); index > 0 {
		return nil, errors.Wrapf(err, "erro ao criar planilha %s", name)
	}

	e.xlsx.NewSheet(name)

	return
}

//
// printTitle prints the cols titles in Excel
//
func (s *Sheet) printTitle(cell string, title string) (err error) {
	xlsx := s.e.xlsx

	// Print header
	xlsx.SetSheetRow(s.name, cell, &[]string{title})

	// Set styles
	style, err := s.e.xlsx.NewStyle(`{"number_format": 0,"font":{"bold":true},"alignment":{"horizontal":"center"},"border":[{"type":"bottom","color":"333333","style":3}]}`)
	if err != nil {
		return errors.Wrap(err, "style")
	}

	s.e.xlsx.SetCellStyle(s.name, cell, cell, style)

	return
}

//
// printRows prints cols in Excel
//
func (s *Sheet) printRows(startingCel string, slice interface{}) error {

	s.e.xlsx.SetSheetRow(s.name, startingCel, slice)

	return nil
}

//
// printValues prints cols in Excel
//
func (s *Sheet) printValue(cell string, value float32) error {

	s.e.xlsx.SetSheetRow(s.name, cell, &[]float32{value})

	// Set styles
	// style, err := s.e.xlsx.NewStyle(`{"number_format": 38}`)
	style, err := s.e.xlsx.NewStyle(`{"custom_number_format": "_-* #,##0_-;_-* (#,##0);_-* \"-\"_-;_-@_-"}`)
	// _-* #.##0_-;_-* (#.##0);_-* "-"_-;_-@_-
	if err != nil {
		return errors.Wrap(err, "style")
	}

	s.e.xlsx.SetCellStyle(s.name, cell, cell, style)

	return nil
}

//
// autoWidth best effort to automatically adjust the cols width
//
func (s *Sheet) autoWidth() {
	cols := "ABCDEFGHIJKLMONPQRSTUVWXYZ"
	var colMaxWidth [26]int
	for _, row := range s.e.xlsx.GetRows(s.name) {
		for c, colCell := range row {
			if c >= len(colMaxWidth) {
				break
			}
			if len(colCell) > colMaxWidth[c] {
				colMaxWidth[c] = len(colCell)
			}
		}
	}
	for c, width := range colMaxWidth {
		col := string(cols[c])
		if width > 0 {
			w := float64(width)
			if w > 10 {
				w -= 6
			}
			if w > 40 {
				w -= 8
			}
			s.e.xlsx.SetColWidth(s.name, col, col, w+4)
		}
	}
}

//
// axis transforms (2, 3) into "B3"
//
func axis(col, row int) string {
	return excelize.ToAlphaString(col) + strconv.Itoa(row)
}

//
// cell2axis only works from A1 to Z999
//
func cell2axis(cell string) (col, row int) {
	col = int(cell[0] - 'A')
	row, _ = strconv.Atoi(cell[1:])

	return
}
