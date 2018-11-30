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
	baseItems := make([]bool, len(items)+row)
	for _, it := range items {
		var sp string
		sp, baseItems[row] = adjustSpace(it.cdConta)
		cell := "A" + strconv.Itoa(row)
		sheet.printRows(cell, &[]string{sp + it.cdConta, sp + it.dsConta}, baseItems[row])
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
			sheet.printValue(cell, statements[it.hash], baseItems[row])
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
// For items equal or above 3, only returns spaces after 2nd level:
// "3.01 ABC"    => ""
// "3.01.01 ABC" => "  "
//
func adjustSpace(str string) (spaces string, baseItem bool) {
	num := strings.SplitN(str, ".", 2)[0]
	c := strings.Count(str, ".")
	if num != "1" && num != "2" && c > 0 {
		c--
	}
	if c > 0 {
		spaces = strings.Repeat("  ", c)
	}

	if num == "1" || num == "2" {
		baseItem = c <= 1
	} else {
		baseItem = c == 0
	}

	return
}
