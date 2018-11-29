package reports

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

//
// Report company from DB to Excel
//
func Report(db *sql.DB, company string, begin, end int, filepath string) (err error) {
	e := newExcel()
	sheet, _ := e.newSheet(company)

	// Print statements code and description
	items, _ := statementItems(db, company)
	row := 2
	for _, it := range items {
		sp := adjustSpace(it.cdConta)
		cell := "A" + strconv.Itoa(row)
		sheet.printRows(cell, &[]string{sp + it.cdConta, sp + it.dsConta})
		row++
	}

	// Print statements values
	cols := "CDEFGHIJKLMONPQRSTUVWXYZ"
	for y := begin; y <= end; y++ {
		col := string(cols[y-begin])
		cell := col + "1"
		sheet.printTitle(cell, "["+strconv.Itoa(y)+"]")

		statements, _ := financialReport(db, company, y)
		row = 2
		for _, it := range items {
			cell := col + strconv.Itoa(row)
			sheet.printValue(cell, statements[it.hash])
			row++
		}
	}

	sheet.autoWidth()
	err = e.saveAndCloseExcel(filepath)

	if err == nil {
		fmt.Printf("[✓] Dados salvos em %s\n", filepath)
	}

	return
}

//
// adjustSpace returns the number of spaces according to the code level, e.g.:
// "1.1 ABC"   => "  " (2 spaces)
// "1.1.1 ABC" => "    " (4 spaces)
//
func adjustSpace(str string) (spaces string) {
	c := strings.Count(str, ".")
	spaces = strings.Repeat("  ", c)

	return
}
