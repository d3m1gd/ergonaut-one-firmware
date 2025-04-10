package main

import (
	"bytes"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	output = flag.String("output", "reference_gen.go", "Output file name")
)

type StructInfo struct {
	Name   string
	Fields []string
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("referencegen: ")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("no input files")
	}

	structs := []StructInfo{}

	for _, filename := range flag.Args() {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
		if err != nil {
			log.Fatalf("parsing %s: %v", filename, err)
		}

		for _, decl := range node.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				fields := []string{}
				for _, f := range st.Fields.List {
					for _, n := range f.Names {
						fields = append(fields, n.Name)
					}
				}
				structs = append(structs, StructInfo{Name: ts.Name.Name, Fields: fields})
			}
		}
	}

	if len(structs) == 0 {
		log.Fatal("no structs with fields found")
	}

	var buf bytes.Buffer
	tmpl := template.Must(template.New("tmpl").Funcs(template.FuncMap{
		"lower": strings.ToLower,
	}).Parse(strings.Trim(`
		// Code generated by referencegen; DO NOT EDIT.

		package main

		import "fmt"
		{{range .}}
		func (x {{.Name}}) Reference() string {
			return fmt.Sprintf("&%s{{range .Fields}} %s{{end}}", x.Name(){{range $i, $f := .Fields}}, x.{{$f}}{{end}})
		}

		func (x {{.Name}}) Name() string {
			return "{{lower .Name}}"
		}

		func (x {{.Name}}) Args() []string {
			return []string{ {{range .Fields}}fmt.Sprintf("%s", x.{{.}}), {{end}}}
		}
		{{end}}`, "\n")))

	err := tmpl.Execute(&buf, structs)
	if err != nil {
		log.Fatalf("template execute: %v", err)
	}

	out := *output
	if !filepath.IsAbs(out) {
		out = filepath.Join(".", out)
	}
	if err := os.WriteFile(out, []byte(dedent(buf.String())), 0644); err != nil {
		log.Fatalf("writing output: %v", err)
	}
}

func dedent(s string) string {
	lines := strings.Split(s, "\n")
	// Skip empty lines
	var prefix string
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if trimmed == "" {
			continue
		}
		leading := line[:len(line)-len(trimmed)]
		if prefix == "" || len(leading) < len(prefix) {
			prefix = leading
		}
	}
	if prefix == "" {
		return s
	}
	for i, line := range lines {
		lines[i] = strings.TrimPrefix(line, prefix)
	}
	return strings.Join(lines, "\n")
}
