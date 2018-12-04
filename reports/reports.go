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
		cell = fmt.Sprintf("%s%d", col, row)
		sheet.printTitle(cell, "["+strconv.Itoa(y)+"]") // Print year as title
		row++
		seq := []string{
			"Patrim. Líq.",
			"",
			"Receita Líq.",
			"EBITDA",
			"D&A",
			"EBIT",
			"Lucro Líq.",
			"",
			"Marg. EBITDA",
			"Marg. EBIT",
			"Marg. Líq.",
			"ROE",
			"",
			"Caixa",
			"Dívida Bruta",
			"Dívida Líq.",
			"Dív. Bru./PL",
			"Dív.Líq./EBITDA",
			"",
			"FCO",
			"FCI",
			"FCF",
			"Fluxo de Caixa Total",
			"",
			"Proventos",
			"Payout",
		}
		metrics, _ := financialMetric(accounts, values)
		for _, m := range seq {
			if m != "" {
				if col == "C" {
					sheet.printRows("B"+strconv.Itoa(row), &[]string{m}, false)
				}
				cell := col + strconv.Itoa(row)
				sheet.printValue(cell, metrics[m], false)
			}
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
		"Patrim. Líq.":                  "(?i)Patrim.nio L.quido.*",
		"D&A":                           "(?i)Deprecia..o.*",
		"Receita Líq.":                  "(?i)Receita de Venda de Bens e/ou Serviços.*",
		"EBIT":                          "(?i)Resultado Antes do Resultado Financeiro e dos Tributos.*",
		"Lucro Líq.":                    "(?i)Lucro/Preju.zo Consolidado do Per.odo.*",
		"Juros sobre o Capital Próprio": "(?i)Juros sobre o Capital Pr.prio.*",
		"Dividendos":                    "(?)Dividendos",
		"Caixa e Equiv.":                "(?)Caixa e Equivalentes de Caixa.*",
		"Aplic. Financeiras":            "(?)Aplic.*Financeiras$",
		"FCO":                           "(?)Caixa L.quido Atividades Operacionais",
		"FCI":                           "(?)Caixa L.quido Atividades de Investimento",
		"FCF":                           "(?)Caixa L.quido Atividades de Financiamento",
	}

	val = make(map[string]float32, len(list)+5)
	for _, acct := range accounts {
		for key := range list {
			match, _ := regexp.MatchString(list[key], acct.dsConta)
			if match {
				val[key] = values[acct.hash]
				continue
			} else
			// There are 4 entries for "Empréstimos e Financiamentos"
			// Need this workaround to get the right ones
			if acct.cdConta+acct.dsConta == "2.01.04Empréstimos e Financiamentos" {
				val["Dívida Circulante"] = values[acct.hash]
			} else if acct.cdConta+acct.dsConta == "2.02.01Empréstimos e Financiamentos" {
				val["Dívida Não Circulante"] = values[acct.hash]
			}
		}
	}
	val["EBITDA"] = val["EBIT"] - val["D&A"]
	val["Marg. EBITDA"] = val["EBITDA"] / val["Receita Líq."]
	val["Marg. EBIT"] = val["EBIT"] / val["Receita Líq."]
	val["ROE"] = val["Lucro Líq."] / val["Patrim. Líq."]
	val["Marg. Líq."] = val["Lucro Líq."] / val["Receita Líq."]
	val["Proventos"] = val["Dividendos"] + val["Juros sobre o Capital Próprio"]
	val["Payout"] = val["Proventos"] / val["Lucro Líq."]
	val["Caixa"] = val["Caixa e Equiv."] + val["Aplic. Financeiras"]
	val["Dívida Bruta"] = val["Dívida Circulante"] + val["Dívida Não Circulante"]
	val["Dívida Líq."] = val["Dívida Bruta"] - val["Caixa"]
	val["Dív. Bru./PL"] = 0.0
	if val["Dívida Bruta"] > 0 {
		val["Dív. Bru./PL"] = val["Dívida Bruta"] / val["Patrim. Líq."]
	}
	val["Dív.Líq./EBITDA"] = 0.0
	if val["Dívida Líq."] > 0 {
		val["Dív.Líq./EBITDA"] = val["Dívida Líq."] / val["EBITDA"]
	}
	val["Fluxo de Caixa Total"] = val["FCO"] + val["FCI"] + val["FCF"]

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
