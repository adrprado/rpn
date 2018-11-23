package parsers

import (
	"database/sql"

	"github.com/pkg/errors"
)

const createTableBpp = `CREATE TABLE IF NOT EXISTS bpp
	(
		"ID" PRIMARY KEY,
		"CNPJ_CIA" varchar(20),
		"DT_REFER" integer,
		"VERSAO" integer,
		"DENOM_CIA" varchar(100),
		"CD_CVM" integer,
		"GRUPO_DFP" varchar(206),
		"MOEDA" varchar(4),
		"ESCALA_MOEDA" varchar(7),
		"ORDEM_EXERC" varchar(9),
		"DT_FIM_EXERC" integer,
		"CD_CONTA" varchar(18),
		"DS_CONTA" varchar(100),
		"VL_CONTA" real
	)
;`

//
// createBPPTable creates the table if not created yet
//
func createBPPTable(db *sql.DB) (err error) {
	statement, err := db.Prepare(createTableBpp)
	if err != nil {
		return errors.Wrap(err, "erro ao preparar tabela")
	}

	_, err = statement.Exec()
	if err != nil {
		return errors.Wrap(err, "erro ao criar tabela")
	}

	// rows, _ := db.Query("SELECT CNPJ_CIA, DT_REFER, DENOM_CIA FROM bpp")
	// var cnpj string
	// var dtRef int
	// var cia string
	// for rows.Next() {
	// 	rows.Scan(&cnpj, &dtRef, &cia)
	// 	fmt.Println(strconv.Itoa(dtRef) + ": " + cnpj + " " + cia)
	// }

	return nil

}
