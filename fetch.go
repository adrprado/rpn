package rapina

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/adrprado/rapina/parsers"
	"github.com/pkg/errors"
)

const dataDir = "data"

// FetchCVM will get data from .zip files downloaded
// directly from CVM
func FetchCVM(dataType, year string) (err error) {
	var files []string
	p, _ := parsers.NewParsers()

	// Check year
	y, err := strconv.Atoi(year)
	if err != nil || y < 1900 || y > 2100 {
		return errors.Wrap(err, "ano inválido")
	}

	// Check data type
	switch dataType {
	case "BPA":
		if files, err = fetchFiles(dataType, year); err != nil {
			return err
		}
		if err = p.BPA(); err != nil {
			return err
		}

	default:
		return errors.Errorf("tipo de informação não existente (%s)", dataType)

	}

	// Clean up
	for _, f := range files {
		err = os.Remove(f)
		if err != nil {
			fmt.Println("could not delete file", f)
		}
	}

	return nil
}

// downloadFile source: https://stackoverflow.com/a/33853856/276311
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

// fetchFiles on CVM server
func fetchFiles(dataType, year string) (files []string, err error) {
	dt := strings.ToLower(dataType)
	url := fmt.Sprintf("http://dados.cvm.gov.br/dados/CIA_ABERTA/DOC/DFP/%s/DADOS/%s_cia_aberta_%s.zip", dataType, dt, year)
	outfile := fmt.Sprintf("%s_%s.zip", dt, year)

	fmt.Printf("[x] Baixando %s de %s\n", dataType, year)
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
