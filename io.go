package rapina

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/adrprado/rapina/parsers"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const dataDir = ".data"

//
// FetchYears fetches all statements from a range
// of years
//
func FetchYears(begin, end int) (err error) {
	db, err := openDatabase()
	if err != nil {
		return err
	}

	for year := begin; year <= end; year++ {
		if err = FetchCVM(db, "BPA", year); err != nil {
			fmt.Printf("[ ] Erro ao processar BPA de %d: %v\n", year, err)
		}
	}

	return err
}

// FetchCVM will get data from .zip files downloaded
// directly from CVM
func FetchCVM(db *sql.DB, dataType string, year int) (err error) {
	var file string

	// Check year
	if year < 1900 || year > 2100 {
		return errors.Wrap(err, "ano inválido")
	}

	// Check data type
	switch dataType {
	case "BPA":
		if file, err = fetchFile(dataType, year); err != nil {
			return err
		}
		if err = parsers.BPA(db, file); err != nil {
			return err
		}

	default:
		return errors.Errorf("tipo de informação não existente (%s)", dataType)

	}

	return nil
}

//
// downloadFile source: https://stackoverflow.com/a/33853856/276311
//
func downloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

//
// fetchFile on CVM server
//
func fetchFile(dataType string, year int) (reqFile string, err error) {
	dt := strings.ToLower(dataType)
	url := fmt.Sprintf("http://dados.cvm.gov.br/dados/CIA_ABERTA/DOC/DFP/%s/DADOS/%s_cia_aberta_%d.zip", dataType, dt, year)
	outfile := fmt.Sprintf("%s_%d.zip", dt, year)

	fmt.Printf("[x] Baixando %s de %d\n", dataType, year)
	err = downloadFile(outfile, url)
	if err != nil {
		return "", errors.Wrap(err, "could not download file")
	}

	var files []string
	files, err = Unzip(outfile, dataDir)
	if err != nil {
		return "", errors.Wrap(err, "could not unzip file")
	}
	files = append(files, outfile)

	// File pattern:
	// xxx_cia_aberta_con_yyy.csv
	reqFile = fmt.Sprintf("%s/%s_cia_aberta_con_%d.csv", dataDir, dt, year)
	idx := find(files, reqFile)
	if idx == -1 {
		filesCleanup(files)
		return "", errors.Errorf("file %s not found", reqFile)
	}

	files[idx] = files[len(files)-1] // Replace it with the last one.
	files = files[:len(files)-1]     // Chop off the last one.
	filesCleanup(files)

	return
}

//
// filesCleanup
//
func filesCleanup(files []string) {
	// Clean up
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			fmt.Println("could not delete file", f)
		}
	}
}

//
// openDatabase to be used by parsers and reporting
//
func openDatabase() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", "./cvm.db")
	if err != nil {
		return db, errors.Wrap(err, "database open failed")
	}

	return
}

//
// find returns the smallest index i at which x == a[i],
// or -1 if there is no such index.
//
func find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}
