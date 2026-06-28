// Command yze-go-anonstruct runs the anonstruct analyzer as a standalone
// go/analysis checker (text, -json, and -fix output, and as a `go vet -vettool`).
package main

import (
	anonstruct "github.com/gomatic/yze-go-anonstruct"
	"golang.org/x/tools/go/analysis/singlechecker"
)

// run is the analysis entry point, indirected so the binary's wiring is testable
// without invoking the real driver (which loads packages and exits the process).
var run = singlechecker.Main

func main() { run(anonstruct.Analyzer) }
