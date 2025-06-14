package macro

import (
	"cmp"
	"slices"
	"strings"

	"keyboard/ref"
	. "keyboard/util"
)

type T = Macro

var macros []Macro

type Macro struct {
	Name  string
	Label string
	Cells int
	Refs  []ref.T
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

func Add(macro Macro) {
	i := slices.IndexFunc(macros, func(other Macro) bool {
		return macro.Name == other.Name
	})
	if i != -1 {
		Panicif(!macro.Equal(macros[i]))
		return
	}
	macros = append(macros, macro)
}

func Placeholder(r ref.T) ref.T {
	return ref.RefN(r.Name, Map(r.Args(), func(string) any { return "MACRO_PLACEHOLDER" }))
}

func Render() []Macro {
	slices.SortFunc(macros, func(a, b Macro) int { return cmp.Compare(a.Name, b.Name) })
	return macros
}
