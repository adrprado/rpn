package reports

import (
	"database/sql"
	"fmt"
)

//
// Exec process all reports from DB to Excel
//
func Exec(db *sql.DB) (err error) {
	e := newExcel()

	dre(db, e)
	bpa(db, e)

	e.saveAndCloseExcel("demo.xlsx")
	return
}

/*
(
		"ID" PRIMARY KEY,
		"CNPJ_CIA" varchar(20),
		"DT_REFER" integer,
		"VERSAO" integer,
		"DENOM_CIA" varchar(100),
		"CD_CVM" integer,
		"GRUPO_DFP" varchar(206),
		"ESCALA_DRE" varchar(7),
		"ORDEM_EXERC" varchar(9),
		"DT_INI_EXERC" integer,
		"DT_FIM_EXERC" integer,
		"CD_CONTA" varchar(18),
		"DS_CONTA" varchar(100),
		"VL_CONTA" real
		)
*/

//
// DRE report - Demonstração de Resultado
//
func dre(db *sql.DB, e *Excel) (err error) {
	s1, _ := e.newSheet("Um")
	s1.printHeader([]string{"a", "b", "c"})
	s2, _ := e.newSheet("Dois")
	s2.printHeader([]string{"a2", "b2", "c2"})

	limit := 0

	rows, _ := db.Query("SELECT * FROM dre") // Note: Ignoring errors for brevity
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
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			// m[colName] = *val
			// fmt.Println(colName, *val)

			switch v := (*val).(type) {
			default:
				fmt.Printf("%s [%T] %v\n", colName, v, *val)
			case []uint8:
				y := *val
				var x = y.([]uint8)
				fmt.Println(colName, "[]uint8", string(x[:]))
			}
		}
		fmt.Println()

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		// fmt.Println(m)
		limit++
		if limit >= 4 {
			break
		}
	}

	return
}

//
// BPA report - Balanço Patrimonial Ativo
//
func bpa(db *sql.DB, e *Excel) {
	s1, _ := e.newSheet("BPA")
	s1.printHeader([]string{"CNPJ", "VAL", "X"})
}
