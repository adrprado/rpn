package reports

import (
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
// printHeader prints the cols titles in Excel
//
func (s *Sheet) printHeader(titles []string) (err error) {
	xlsx := s.e.xlsx

	// Set styles
	styleHeader, err := xlsx.NewStyle(`{"font":{"bold":true,"color":"#333333","size":8},"border":[{"type":"bottom","color":"333333","style":3}]}`)
	if err != nil {
		return errors.Wrap(err, "styleHeader")
	}
	xlsx.SetCellStyle(s.name, "A1", "Z1", styleHeader)

	// Print header
	xlsx.SetSheetRow(s.name, "A1", &titles)

	return
}

//
// printCols prints cols in Excel
//
func printCols(cols []string) error {
	return nil
}

//
// autoWidth best effort to automatically adjust the cols width
//
func (s *Sheet) autoWidth(sheetName string) {
	cols := "ABCDEFGHIJKLMONPQRSTUVWXYZ"
	var colMaxWidth [26]int
	for _, row := range s.e.xlsx.GetRows(sheetName) {
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
			s.e.xlsx.SetColWidth(sheetName, col, col, float64(width+4))
		}
	}
}
