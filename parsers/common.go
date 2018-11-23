package parsers

import (
	"bufio"
	"database/sql"
	"fmt"
	"hash/fnv"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var fnvHash = fnv.New32a()

//
// getHash returns the FNV-1 non-cryptographic hash
//
func getHash(s string) uint32 {
	fnvHash.Write([]byte(s))
	defer fnvHash.Reset()

	return fnvHash.Sum32()
}

//
// populateBPPTable loop thru file and insert its lines into DB
//
func populateTable(db *sql.DB, table, file string) (err error) {
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
			if err = insertLine(tx, table, &header, f, getHash(line)); err != nil {
				fmt.Println("[ ] BPP:", err)
			}
		}

		// fmt.Println("-------------------------------")
		if count++; count%1000 == 0 {
			fmt.Print(".")
		}
	}

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
func insertLine(db *sql.Tx, table string, header *map[string]int, fields []string, hash uint32) error {
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

	insert := fmt.Sprint("INSERT OR IGNORE INTO ", table,
		" (ID, ", strings.Join(names, ","),
		") VALUES (",
		hash, ",", strings.Join(values, ","),
		");")

	// fmt.Println(insert)

	statement, err := db.Prepare(insert)
	if err != nil {
		return errors.Wrapf(err, "erro ao preparar insert em '%s'", table)
	}

	_, err = statement.Exec()
	if err != nil {
		return errors.Wrapf(err, "erro ao inserir dados em '%s'", table)
	}

	return nil
}
