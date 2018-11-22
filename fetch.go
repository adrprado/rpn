package rapina

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/adrprado/rapina/parsers"
	"github.com/pkg/errors"
)

const dataDir = "data"

//
// FetchYears fetches all statements from a range
// of years
//
func FetchYears(begin, end int) (err error) {
	p, err := parsers.NewParsers()
	for year := begin; year <= end; year++ {
		FetchCVM(&p, "BPA", year)
	}

	return err
}

// FetchCVM will get data from .zip files downloaded
// directly from CVM
func FetchCVM(p *parsers.Parsers, dataType string, year int) (err error) {
	var files []string

	// Check year
	if year < 1900 || year > 2100 {
		return errors.Wrap(err, "ano inválido")
	}

	// Check data type
	switch dataType {
	case "BPA":
		if files, err = fetchFiles(dataType, year); err != nil {
			return err
		}
		if err = p.BPA(dataDir); err != nil {
			return err
		}

	default:
		return errors.Errorf("tipo de informação não existente (%s)", dataType)

	}

	// Clean up
	// for _, f := range files {
	// 	err = os.Remove(f)
	// 	if err != nil {
	// 		fmt.Println("could not delete file", f)
	// 	}
	// }
	_ = fetchFiles

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
// fetchFiles on CVM server
//
func fetchFiles(dataType string, year int) (files []string, err error) {
	dt := strings.ToLower(dataType)
	url := fmt.Sprintf("http://dados.cvm.gov.br/dados/CIA_ABERTA/DOC/DFP/%s/DADOS/%s_cia_aberta_%d.zip", dataType, dt, year)
	outfile := fmt.Sprintf("%s_%d.zip", dt, year)

	fmt.Printf("[x] Baixando %s de %d\n", dataType, year)
	err = downloadFile(outfile, url)
	if err != nil {
		return nil, errors.Wrap(err, "could not download file")
	}

	files, err = Unzip(outfile, dataDir)
	if err != nil {
		return nil, errors.Wrap(err, "could not unzip file")
	}
	files = append(files, outfile)
	files = append(files, dataDir)

	return
}
