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
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var fnvHash = fnv.New32a()

//
// GetHash returns the FNV-1 non-cryptographic hash
//
func GetHash(s string) uint32 {
	fnvHash.Write([]byte(s))
	defer fnvHash.Reset()

	return fnvHash.Sum32()
}

//
// whatTable for the data type
//
func whatTable(dataType string) (table string, err error) {
	switch dataType {
	case "BPA", "BPP", "DRE", "DFC_MD", "DFC_MI", "DVA":
		table = "dfp"
	default:
		return "", errors.Wrapf(err, "tipo de informação inexistente: %s", dataType)
	}

	return
}

//
// createTable creates the table if not created yet
//
func createTable(db *sql.DB, dataType string) (err error) {

	table, err := whatTable(dataType)
	if err != nil {
		return err
	}

	statement, err := db.Prepare(createTableMap[table])
	if err != nil {
		return errors.Wrap(err, "erro ao preparar tabela")
	}

	_, err = statement.Exec()
	if err != nil {
		return errors.Wrap(err, "erro ao criar tabela")
	}

	return nil

}

//
// populateBPPTable loop thru file and insert its lines into DB
//
func populateTable(db *sql.DB, dataType, file string) (err error) {
	progress := []string{"/", "-", "\\", "|", "-", "\\"}
	p := 0

	fh, err := os.Open(file)
	if err != nil {
		return errors.Wrapf(err, "erro ao abrir arquivo %s", file)
	}
	defer fh.Close()

	dec := transform.NewReader(fh, charmap.ISO8859_1.NewDecoder())

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
	scanner := bufio.NewScanner(dec)
	count := 0

	// Loop thru file, line by line
	for scanner.Scan() {
		// line := toUtf8([]byte(scanner.Text()))
		line := scanner.Text()
		f := strings.FieldsFunc(line, sep)

		if len(header) == 0 {
			// Get header positioning
			for i, h := range f {
				header[h] = i
			}
		} else {
			if err = insertLine(tx, dataType, &header, f, GetHash(line)); err != nil {
				fmt.Printf("[x] %s: %v\n", dataType, err)
			}
		}

		// fmt.Println("-------------------------------")
		if count++; count%1000 == 0 {
			fmt.Printf("\r[%s", progress[p%6])
			p++
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
func insertLine(db *sql.Tx, dataType string, header *map[string]int, fields []string, hash uint32) (err error) {
	var names, values []string
	for h, i := range *header {
		names = append(names, "`"+h+"`")

		// Change date from 'YYYY-MM-DD' to Unix epoch
		// To convert back from sqlite: strftime('%Y-%m-%d', DT_REFER, 'unixepoch')
		f := ""
		switch h {
		case "DT_REFER", "DT_INI_EXERC", "DT_FIM_EXERC":
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

	table, err := whatTable(dataType)
	if err != nil {
		return err
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

//
// Exec start the data import process, including the database creation
// if necessary
//
func Exec(db *sql.DB, dataType string, file string) (err error) {
	err = createTable(db, dataType)
	if err != nil {
		return err
	}

	fmt.Print("[ ] Processando arquivo ", dataType)
	err = populateTable(db, dataType, file)
	if err == nil {
		fmt.Print("\r[✓")
	}
	fmt.Println()

	return err
}

//
// toUtf8 convert iso8859-1 to utf8
// https://stackoverflow.com/a/13511463/276311
//
func toUtf8(iso8859_1_buf []byte) string {
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}
