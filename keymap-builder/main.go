package main

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"keyboard/behavior"
	"keyboard/combo"
	"keyboard/layer"
	"keyboard/macro"
	. "keyboard/util"
)

func renderKeymap(path string, params Params) {
	tmplPath := path + ".tmpl"
	var funcs = template.FuncMap{"join": strings.Join}
	t := Must(template.New(filepath.Base(path + ".tmpl")).Funcs(funcs).ParseFiles(tmplPath))
	outFile := Must(os.Create(path))
	defer outFile.Close()
	Check(t.Execute(outFile, params))
}

type Params struct {
	Layers    []layer.T
	Macros    []macro.T
	Combos    []combo.T
	Behaviors []behavior.T
}

func main() {
	renderKeymap("../config/ergonaut_one.keymap", Params{
		Behaviors: behavior.Render(),
		Macros:    macro.Render(),
		Combos:    combo.Render(),
		Layers:    layer.All(),
	})
}
