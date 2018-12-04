package reports

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//
// Report company from DB to Excel
//
func Report(db *sql.DB, company string, begin, end int, filepath string) (err error) {
	e := newExcel()
	sheet, _ := e.newSheet(company)

	// Print accounts codes and descriptions in columns A and B
	// starting on row 2. Adjust space related to the group, e.g.:
	// 3.02 ABC <== print in bold if base item and stores the row position in baseItems[]
	//   3.02.01 ABC
	accounts, _ := accountsItems(db, company)
	row := 2
	baseItems := make([]bool, len(accounts)+row)
	for _, it := range accounts {
		var sp string
		sp, baseItems[row] = adjustSpace(it.cdConta)
		cell := "A" + strconv.Itoa(row)
		sheet.printRows(cell, &[]string{sp + it.cdConta, sp + it.dsConta}, baseItems[row])
		row++
	}

	// Print accounts values one year per columns, starting from C, row 2
	var values map[uint32]float32
	cols := "CDEFGHIJKLMONPQRSTUVWXYZ"
	for y := begin; y <= end; y++ {
		if y-begin >= len(cols) {
			break
		}
		col := string(cols[y-begin])
		cell := col + "1"
		sheet.printTitle(cell, "["+strconv.Itoa(y)+"]") // Print year as title in row 1

		values, _ = accountsValues(db, company, y)
		row = 2
		for _, acct := range accounts {
			cell := col + strconv.Itoa(row)
			sheet.printValue(cell, values[acct.hash], baseItems[row])
			row++
		}

		// Print financial metrics
		row++
		if col == "C" {
			sheet.printRows("B"+strconv.Itoa(row), &[]string{"ROE"}, true)
		}
		metrics, _ := financialMetric(accounts, values)
		fmt.Println("[>] Metrics: ", metrics)
		for _, m := range metrics {
			cell := col + strconv.Itoa(row)
			sheet.printValue(cell, m, true)
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

func financialMetric(accounts []accItems, values map[uint32]float32) (val map[string]float32, err error) {
	list := map[string]string{
		"EQUITY":                 "(?i)Patrim.nio L.quido.*",
		"DDA":                    "(?i)Deprecia..o.*",
		"REVENUE":                "(?i)Receita de Venda de Bens e/ou Serviços.*",
		"EBIT":                   "(?i)Resultado Antes do Resultado Financeiro e dos Tributos.*",
		"NET_INCOME":             "(?i)Lucro/Preju.zo Consolidado do Per.odo.*",
		"INTEREST_ON_NET_EQUITY": "(?i)Juros sobre o Capital Pr.prio.*",
		"DIVIDENDS":              "(?)Dividendos",
	}
	val = make(map[string]float32, len(list)+5)

	for _, acct := range accounts {
		for key := range list {
			match, _ := regexp.MatchString(list[key], acct.dsConta)
			if match {
				fmt.Println("[>]", acct.cdConta, acct.dsConta)
				val[key] = values[acct.hash]
				continue
			}
		}
	}
	val["EBITDA"] = val["EBIT"] - val["DDA"]
	val["MARGIN_EBITDA"] = val["EBITDA"] / val["REVENUE"]
	val["MARGIN_EBIT"] = val["EBIT"] / val["REVENUE"]
	val["ROE"] = val["NET_INCOME"] / val["EQUITY"]

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
