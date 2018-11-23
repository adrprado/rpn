package reports

//
// DRE report - Demonstração de Resultado
//
func DRE() {
	e := newExcel()
	s1, _ := e.newSheet("Um")
	s1.printHeader([]string{"a", "b", "c"})
	s2, _ := e.newSheet("Dois")
	s2.printHeader([]string{"a2", "b2", "c2"})
	e.saveAndCloseExcel("demo.xlsx")
}
