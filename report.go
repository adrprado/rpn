package rapina

import (
	"github.com/adrprado/rapina/reports"
	"github.com/pkg/errors"
)

//
// Report a company from DB to Excel
//
func Report(company string) (err error) {
	db, err := openDatabase()
	if err != nil {
		return errors.Wrap(err, "fail to open db")
	}

	return reports.Report(db, company, dataDir+"/"+company+".xlsx")
}
