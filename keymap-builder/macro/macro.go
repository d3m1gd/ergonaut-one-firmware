package macro

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"keyboard/ref"
	. "keyboard/util"
)

type T = Macro

var macros []Macro

var Press = ref.Ref0("macro_press")
var Release = ref.Ref0("macro_release")
var Pause = ref.Ref0("macro_pause_for_release")
var Wait = ref.Ref0("macro_pause_for_release")

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
	return strings.Join(Map(m.Refs, ref.CompileRef), " ")
}

func (m Macro) Equal(other Macro) bool {
	eq := true
	eq = eq && m.Name == other.Name
	eq = eq && m.Label == other.Label
	eq = eq && m.Cells == other.Cells
	eq = eq && slices.EqualFunc(m.Refs, other.Refs, ref.EqualRef)
	return eq
}

func AddMacro(macro Macro) {
	i := slices.IndexFunc(macros, func(other Macro) bool {
		return macro.Name == other.Name
	})
	if i != -1 {
		Panicif(!macro.Equal(macros[i]))
		return
	}
	macros = append(macros, macro)
}

var Param11 = macroParamBuilder(1, 1)
var Param12 = macroParamBuilder(1, 2)
var Param21 = macroParamBuilder(2, 1)
var Param22 = macroParamBuilder(2, 2)

func macroParamBuilder(a, b int) ref.T {
	return ref.Ref0(fmt.Sprintf("macro_param_%dto%d", a, b))
}

func Placeholder(r ref.T) ref.T {
	return ref.RefN(r.Name, MapToMacroPlaceholder(r.Args()))
}

func MapToMacroPlaceholder(args []string) []any {
	return Map(args, func(_ string) any { return "MACRO_PLACEHOLDER" })
}

func MapParams(n int) []ref.T {
	switch n {
	case 0:
		return []ref.T{}
	case 1:
		return []ref.T{Param11}
	case 2:
		return []ref.T{Param11, Param22}
	}
	panic(fmt.Sprintf("bad n: %d", n))
}

func Render() []Macro {
	slices.SortFunc(macros, func(a, b Macro) int { return cmp.Compare(a.Name, b.Name) })
	return macros
}
