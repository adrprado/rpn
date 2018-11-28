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
	// s1, _ := e.newSheet(company)
	// rowIdx := 1
	// items, _ := statementItems(db, company)
	// for _, item := range items {
	// 	item = adjustSpace(item)
	// 	fmt.Println(item)
	// 	row := strconv.Itoa(rowIdx)
	// 	s1.printRows("A"+row, &[]string{item})
	// 	rowIdx++
	// }

	e := newExcel()
	s1, _ := e.newSheet(company)

	// Print statements code and description
	items, _ := statementItems(db, company)
	row := 2
	for _, it := range items {
		sp := adjustSpace(it.cdConta)
		cell := "A" + strconv.Itoa(row)
		s1.printRows(cell, &[]string{sp + it.cdConta, sp + it.dsConta})
		row++
	}

	// Print statements values
	cols := "CDEFGHIJKLMONPQRSTUVWXYZ"
	for y := begin; y <= end; y++ {
		col := string(cols[y-begin])
		cell := col + "1"
		s1.printRows(cell, &[]int{y})

		statements, _ := financialReport(db, company, y)
		row = 2
		for _, it := range items {
			fmt.Println(y, it.cdConta, it.dsConta, statements[it.hash])
			cell := col + strconv.Itoa(row)
			s1.printRows(cell, &[]float32{statements[it.hash]})
			row++
		}
	}

	e.saveAndCloseExcel(filepath)

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
