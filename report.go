package rapina

import (
	"github.com/adrprado/rapina/reports"
	"github.com/pkg/errors"
)

//
// Exec process all reports from DB to Excel
//
func Exec() (err error) {
	db, err := openDatabase()
	if err != nil {
		return errors.Wrap(err, "fail to open db")
	}

	return reports.Exec(db)
}
