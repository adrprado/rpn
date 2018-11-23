package parsers

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const createTableBpa = `CREATE TABLE IF NOT EXISTS bpa
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
// BPA = Balanço Patrimonial Ativo
//
func BPA(db *sql.DB, file string) (err error) {
	fmt.Println("[ ] Criando/conferindo banco de dados...")
	err = createBPATable(db)
	if err != nil {
		return err
	}
	fmt.Println("[x] ok")

	fmt.Println("[ ] Processando arquivo da CVM...")
	err = populateBPATable(db, file)

	return err
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

	// rows, _ := db.Query("SELECT CNPJ_CIA, DT_REFER, DENOM_CIA FROM bpa")
	// var cnpj string
	// var dtRef int
	// var cia string
	// for rows.Next() {
	// 	rows.Scan(&cnpj, &dtRef, &cia)
	// 	fmt.Println(strconv.Itoa(dtRef) + ": " + cnpj + " " + cia)
	// }

	return nil

}

func populateBPATable(db *sql.DB, file string) (err error) {
	fh, err := os.Open(file)
	if err != nil {
		return errors.Wrapf(err, "erro ao abrir arquivo %s", file)
	}
	defer fh.Close()

	// dec := transform.NewReader(file, charmap.Windows1252.NewDecoder())

	// BEGIN TRANSACTION
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "Failed to begin transaction")
	}

	// Data used inside loop
	sep := func(r rune) bool {
		return r == ';'
	}
	header := make(map[string]int) // stores the header item position (e.g., DT_FIM_EXERC:9)
	scanner := bufio.NewScanner(fh)
	count := 0

	// Loop thru file, line by line
	for scanner.Scan() {
		line := scanner.Text()
		f := strings.FieldsFunc(line, sep)

		if len(header) == 0 {
			// Get header positioning
			for i, h := range f {
				header[h] = i
			}
		} else {
			if err = insertLine(tx, &header, f, getHash(line)); err != nil {
				fmt.Println("[ ] BPA:", err)
			}
		}

		// fmt.Println("-------------------------------")
		if count++; count%1000 == 0 {
			fmt.Print(".")
		}
	}
	fmt.Print("\n")

	// END TRANSACTION
	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "Failed to commit transaction")
	}

	if err := scanner.Err(); err != nil {
		return errors.Wrapf(err, "erro ao ler arquivo %s", file)
	}

	return err
}

//
// insertLine into DB
//
func insertLine(db *sql.Tx, header *map[string]int, fields []string, hash uint32) error {
	var names, values []string
	for h, i := range *header {
		names = append(names, "`"+h+"`")

		// Change date from 'YYYY-MM-DD' to Unix epoch
		// To convert back from sqlite: strftime('%Y-%m-%d', DT_REFER, 'unixepoch')
		f := ""
		switch h {
		case "DT_REFER", "DT_FIM_EXERC":
			layout := "2006-01-02"
			t, err := time.Parse(layout, fields[i])
			if err != nil {
				return errors.Wrap(err, "data invalida "+fields[i])
			}
			f = fmt.Sprintf("%v", t.Unix())

		default:
			f = "'" + fields[i] + "'"
		}

		values = append(values, f)
	}

	insert := fmt.Sprint("INSERT OR IGNORE INTO bpa (ID, ",
		strings.Join(names, ", "),
		") VALUES (", hash, ",",
		strings.Join(values, ", "),
		");")

	// fmt.Println(insert)

	statement, err := db.Prepare(insert)
	if err != nil {
		return errors.Wrap(err, "erro ao preparar insert em `bpa`")
	}

	_, err = statement.Exec()
	if err != nil {
		return errors.Wrap(err, "erro ao inserir dados em `bpa`")
	}

	return nil
}
