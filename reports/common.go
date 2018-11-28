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
func Report(db *sql.DB, company, year, filepath string) (err error) {
	e := newExcel()

	s1, _ := e.newSheet(company)
	rowIdx := 1
	items, _ := statementItems(db, company)
	for _, item := range items {
		item = adjustSpace(item)
		fmt.Println(item)
		row := strconv.Itoa(rowIdx)
		s1.printRows("A"+row, &[]string{item})
		rowIdx++
	}
	financialReport(db, e, company, year)

	e.saveAndCloseExcel(filepath)
	return
}

//
// adjustSpace adds spaces in the begining of the statement
// item. E.g.: "1.1 ABC" => "  1.1 ABC"
//
func adjustSpace(str string) string {
	code := strings.SplitN(str, " ", 2)
	c := strings.Count(code[0], ".")
	str = strings.Repeat("  ", c) + str
	return str
}
