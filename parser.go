package rapina

import "fmt"

// ParseCVM will get data from .zip files downloaded
// directly from CVM
func ParseCVM(dataType string, year int) (ok bool) {
	fmt.Printf("[x] Baixando %s de %d\n", dataType, year)
	return true
}
