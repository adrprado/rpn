package main

import (
	"fmt"
	"os"

	"github.com/adrprado/rapina/reports"
)

var (
	version string
	build   string
)

func main() {

	fmt.Fprint(os.Stderr, "Rapina - Dados Financeiros de Empresas via CVM - ")
	fmt.Fprintf(os.Stderr, "v%s-%s\n", version, build)
	fmt.Fprint(os.Stderr, "(2018) github.com/adrprado/rapina\n\n")

	reports.DRE()
	// err := rapina.FetchCVM(2014, 2017)
	// if err != nil {
	// 	fmt.Println("[x]", err)
	// 	os.Exit(1)
	// }

}
