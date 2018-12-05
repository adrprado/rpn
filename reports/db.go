package reports

import (
	"database/sql"
	"fmt"

	"github.com/adrprado/rapina/parsers"
	"github.com/pkg/errors"
)

type accItems struct {
	hash    uint32
	cdConta string
	dsConta string
}

//
// accountsItems returns all accounts codes and descriptions, e.g.:
// [1 Ativo Total, 1.01 Ativo Circulante, ...]
//
func accountsItems(db *sql.DB, company string) (items []accItems, err error) {
	selectItems := fmt.Sprintf(`
	SELECT DISTINCT
		CD_CONTA, DS_CONTA
	FROM
		dfp
	WHERE
		DENOM_CIA LIKE "%s%%"
		AND ORDEM_EXERC LIKE "_LTIMO"

	ORDER BY
		CD_CONTA, DS_CONTA
	;`, company)

	rows, err := db.Query(selectItems)
	if err != nil {
		panic(err)
	}

	var item accItems
	for rows.Next() {
		rows.Scan(&item.cdConta, &item.dsConta)
		item.hash = parsers.GetHash(item.cdConta + item.dsConta)
		items = append(items, item)
	}

	// genericPrint(rows)

	return
}

type account struct {
	date     string
	denomCia string
	escala   string
	cdConta  string
	dsConta  string
	vlConta  float32
}

//
// accountsValues stores the values for each account into a map using a hash
// of the account code and description as its key
//
func accountsValues(db *sql.DB, company string, year int) (values map[uint32]float32, err error) {

	selectReport := fmt.Sprintf(`
	SELECT
		strftime('%%Y-%%m-%%d', DT_REFER, 'unixepoch') AS DT,
		DENOM_CIA,
		CD_CONTA,
		DS_CONTA,
		VL_CONTA
	FROM
		dfp
	WHERE
		DENOM_CIA LIKE "%s%%"
		AND ORDEM_EXERC LIKE "_LTIMO"
		AND DT = "%d-12-31"

	ORDER BY
		DT, CD_CONTA
	;`, company, year)

	values = make(map[uint32]float32)
	st := account{}

	rows, err := db.Query(selectReport)
	for rows.Next() {
		rows.Scan(
			&st.date,
			&st.denomCia,
			&st.cdConta,
			&st.dsConta,
			&st.vlConta,
		)

		values[parsers.GetHash(st.cdConta+st.dsConta)] = st.vlConta
	}

	return
}

//
// genericPrint prints the entire row
//
func genericPrint(rows *sql.Rows) (err error) {
	limit := 0
	cols, _ := rows.Columns()
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		// m := make(map[string]interface{})
		for i := range cols {
			val := columnPointers[i].(*interface{})
			// m[colName] = *val
			// fmt.Println(colName, *val)

			switch (*val).(type) {
			default:
				fmt.Print(*val, ";")
			case []uint8:
				y := *val
				var x = y.([]uint8)
				fmt.Print(string(x[:]), ";")
			}
		}
		fmt.Println()

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		// fmt.Println(m)
		limit++
		if limit >= 4000 {
			break
		}
	}

	return
}

//
// companies returns available companies in the DB
//
func companies(db *sql.DB) (list []string, err error) {

	selectCompanies := `
		SELECT DISTINCT
			DENOM_CIA
		FROM
			dfp
		ORDER BY
			DENOM_CIA
		;`

	rows, err := db.Query(selectCompanies)
	if err != nil {
		return nil, errors.Wrap(err, "falha ao ler banco de dados")
	}

	list = make([]string, 0, 10)
	var companyName string
	for rows.Next() {
		rows.Scan(&companyName)
		list = append(list, companyName)
	}

	return
}
