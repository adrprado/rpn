package parsers

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

// Parsers struct
type Parsers struct {
	db *sql.DB
}

// NewParsers creates a new Parsers
func NewParsers() (Parsers, error) {
	p := Parsers{}
	var err error

	p.db, err = sql.Open("sqlite3", "./cvm.db")
	if err != nil {
		return p, errors.Wrap(err, "database open failed")
	}

	return p, nil
}

// func (p *Parsers) createDB() {
// 	statement, _ := p.db.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
// 	statement.Exec()
// 	statement, _ = p.db.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
// 	statement.Exec("Nic", "Noc")
// 	rows, _ := p.db.Query("SELECT id, firstname, lastname FROM people")
// 	var id int
// 	var firstname string
// 	var lastname string
// 	for rows.Next() {
// 		rows.Scan(&id, &firstname, &lastname)
// 		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
// 	}
// }
