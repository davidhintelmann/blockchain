package ttmpl

import (
	"os"
	"text/template"

	"github.com/davidhintelmann/blockchain/bparser"
)

// Output a single blocks details to the terminal.
// Used in ParseBlocks function.
func PrintBlock(block bparser.BlockData) {
	tmpl, err := template.ParseFiles("block.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, block)
	if err != nil {
		panic(err)
	}
}
