package main

import (
	"fmt"
	"os"

	"github.com/adrprado/rapina"
)

var (
	version string
	build   string
)

func main() {

	fmt.Fprint(os.Stderr, "Rapina - Dados Financeiros de Empresas via CVM - ")
	fmt.Fprintf(os.Stderr, "v%s-%s\n", version, build)
	fmt.Fprint(os.Stderr, "(2018) github.com/adrprado/rapina ツ\n\n")

	err := rapina.FetchCVM(2017, 2017)
	if err != nil {
		fmt.Println("[x]", err)
		os.Exit(1)
	}

	rapina.Report("WEG", "2017")
}
