package macro

import (
	"cmp"
	"slices"
	"strings"

	"keyboard/ref"
	. "keyboard/util"
	"keyboard/util/indenter"
)

type T = Macro

var macros []Macro

type Macro struct {
	Name   string
	Cells  int
	WaitMs int
	Refs   []ref.T
}

func (m Macro) Type() string {
	return map[int]string{
		0: "behavior-macro",
		1: "behavior-macro-one-param",
		2: "behavior-macro-two-param",
	}[m.Cells]
}

func (m Macro) Bindings() string {
	return strings.Join(Map(m.Refs, ref.Compile), " ")
}

func (m Macro) Equal(other Macro) bool {
	eq := true
	eq = eq && m.Name == other.Name
	eq = eq && m.Cells == other.Cells
	eq = eq && slices.EqualFunc(m.Refs, other.Refs, ref.Equal)
	return eq
}

func Add(macro Macro) Macro {
	i := slices.IndexFunc(macros, func(other Macro) bool {
		return macro.Name == other.Name
	})
	if i != -1 {
		Panicif(!macro.Equal(macros[i]))
		return macro
	}
	macros = append(macros, macro)
	return macro
}

func (m Macro) Invoke(args ...any) ref.Ref {
	Panicif(len(args) != m.Cells)
	return ref.RefN(m.Name, args)
}

func Render() []Macro {
	slices.SortFunc(macros, func(a, b Macro) int { return cmp.Compare(a.Name, b.Name) })
	return macros
}

func (m Macro) Compile(indent, level int) string {
	ir := indenter.New(indent)

	ir.Sprintf(0, "\n")
	ir.Sprintf(level, "%s: %s {\n", m.Name, m.Name)
	ir.Sprintf(level+1, "compatible = \"zmk,%s\";\n", m.Type())
	ir.Sprintf(level+1, "#binding-cells = <%d>;\n", m.Cells)
	if m.WaitMs != 0 {
		ir.Sprintf(level+1, "wait-ms = <%d>;\n", m.WaitMs)
	}
	ir.Sprintf(level+1, "bindings = <%s>;\n", m.Bindings())
	ir.Sprintf(level, "};\n")
	return ir.String()
}
