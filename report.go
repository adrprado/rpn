package rapina

import (
	"strings"

	"github.com/adrprado/rapina/reports"
	"github.com/pkg/errors"
)

//
// Report a company from DB to Excel
//
func Report(company string, begin, end int) (err error) {
	// Check year
	if begin < 1900 || begin > 2100 || end < 1900 || end > 2100 {
		return errors.Wrap(err, "ano inválido")
	}
	if begin > end {
		aux := end
		end = begin
		begin = aux
	}

	db, err := openDatabase()
	if err != nil {
		return errors.Wrap(err, "fail to open db")
	}

	company = strings.ToUpper(company)

	return reports.Report(db, company, begin, end, dataDir+"/"+company+".xlsx")
}

//
// ListCompanies a company from DB to Excel
//
func ListCompanies() (err error) {
	db, err := openDatabase()
	if err != nil {
		return errors.Wrap(err, "fail to open db")
	}

	return reports.ListCompanies(db)
}
