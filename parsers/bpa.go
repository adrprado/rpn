package parsers

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

const createTableBpa = `CREATE TABLE IF NOT EXISTS bpa
	(
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

const insertIntoBpa = `INSERT INTO bpa
    ("CNPJ_CIA", "DT_REFER", "VERSAO", "DENOM_CIA", "CD_CVM", "GRUPO_DFP", "MOEDA", "ESCALA_MOEDA", "ORDEM_EXERC", "DT_FIM_EXERC", "CD_CONTA", "DS_CONTA", "VL_CONTA")
VALUES
    ('84.429.695/0001-11', '31/12/2017', 1, 'WEG SA', 5410, 'DF Consolidado - Balanço Patrimonial Ativo', 'REAL', 'MILHAR', 'ÚLTIMO', '31/12/2017', 1, 'Ativo Total', 13985987.00)
;`

//
// BPA = Balanço Patrimonial Ativo
//
func BPA(db *sql.DB, dir string) (err error) {
	fmt.Println("[ ] Criando/conferindo banco de dados...")
	err = createBPATable(db)
	if err != nil {
		return err
	}
	fmt.Println("[x] ok")

	fmt.Println("[ ] Lendo arquivo obtido do servidor da CVM...")
	err = populateBPATable(db, dir)

	return nil
}

//
// createBPATable creates the table if not created yet
//
func createBPATable(db *sql.DB) (err error) {
	statement, err := db.Prepare(createTableBpa)
	if err != nil {
		return errors.Wrap(err, "erro ao preparar tabela")
	}

	_, err = statement.Exec()
	if err != nil {
		return errors.Wrap(err, "erro ao criar tabela")
	}

	statement, err = db.Prepare(insertIntoBpa)
	statement.Exec()
	if err != nil {
		return errors.Wrap(err, "erro ao inserir dados em `bpa`")
	}

	rows, _ := db.Query("SELECT CNPJ_CIA, DT_REFER, DENOM_CIA FROM bpa")
	var cnpj string
	var dtRef int
	var cia string
	for rows.Next() {
		rows.Scan(&cnpj, &dtRef, &cia)
		fmt.Println(strconv.Itoa(dtRef) + ": " + cnpj + " " + cia)
	}

	return nil

}

func populateBPATable(db *sql.DB, dir string) (err error) {

	return err
}
