package main

import (
	"fmt"
	"slices"
	"strings"
)

var MacroPress = Custom0("macro_press")
var MacroRelease = Custom0("macro_release")
var MacroWait = Custom0("macro_pause_for_release")

type Macro struct {
	Name  string
	Label string
	Cells int
	Refs  []Reference
}

func (m Macro) Type() string {
	return map[int]string{
		0: "behavior-macro",
		1: "behavior-macro-one-param",
		2: "behavior-macro-two-param",
	}[m.Cells]
}

func (m Macro) Bindings() string {
	return strings.Join(Map(m.Refs, CompileReference), " ")
}

func (m Macro) Equal(other Macro) bool {
	eq := true
	eq = eq && m.Name == other.Name
	eq = eq && m.Label == other.Label
	eq = eq && m.Cells == other.Cells
	eq = eq && slices.EqualFunc(m.Refs, other.Refs, EqualReference)
	return eq
}

func BackspaceDelete() Reference {
	name := "bspcdel"
	AddMacro(Macro{
		Name:  name,
		Label: "Backspace Delete",
		Cells: 0,
		Refs:  []Reference{Kp(BSPC), Kp(DEL)},
	})

	return Custom0(name)
}

func Parens() Reference {
	return OpenCloseMacro("parens", LPAR, RPAR)
}

func Brackets() Reference {
	return OpenCloseMacro("brackets", LBKT, RBKT)
}

func Curlies() Reference {
	return OpenCloseMacro("curlies", LBRC, RBRC)
}

func OpenCloseMacro(name string, left, right KeyCode) Reference {
	AddMacro(Macro{
		Name:  name,
		Label: fmt.Sprintf("OpenClose %s", name),
		Cells: 0,
		Refs:  []Reference{Kp(left), Kp(right), Kp(LEFT), To{PARENS}},
	})

	return Custom0(name)
}

func AddMacro(macro Macro) {
	i := slices.IndexFunc(macros, func(other Macro) bool {
		return macro.Name == other.Name
	})
	if i != -1 {
		panicif(!macro.Equal(macros[i]))
		return
	}
	macros = append(macros, macro)
}

type MacroParam struct {
	a, b int
}

func (mp MacroParam) Reference() string {
	return fmt.Sprintf("&%s", mp.Name())
}

func (mp MacroParam) Name() string {
	return fmt.Sprintf("macro_param_%dto%d", mp.a, mp.b)
}

func (mp MacroParam) Args() []string {
	return []string{}
}

func (mp MacroParam) Slots() int {
	return 0
}

func Curry(r Reference) Custom {
	return Custom{r.Name(), MapToMacroPlaceholder(r.Args())}
}

func MapToMacroPlaceholder(args []string) []any {
	return MapToAnyStatic(args, "MACRO_PLACEHOLDER")
}

func XThenTrans(r Reference, index LayerIndex, rc RC) Reference {
	layer := layers[index]
	name := fmt.Sprintf("xThenTrans%d", index)
	AddMacro(Macro{
		Name:  name,
		Label: fmt.Sprintf("X Then Trans %s", index),
		Cells: 1,
		Refs:  []Reference{MacroParam{1, 1}, Curry(r), To{index}, layer[rc]},
	})

	return Custom{name, MapToAny(r.Args())}
}

func XThenLayer(r Reference, index LayerIndex) Reference {
	name := fmt.Sprintf("xThenLayer%d", index)
	AddMacro(Macro{
		Name:  name,
		Label: fmt.Sprintf("X Then Layer %s", index),
		Cells: 1,
		Refs:  []Reference{MacroParam{1, 1}, Curry(r), To{index}},
	})

	return Custom{name, MapToAny(r.Args())}
}

func MapParams(n int) []Reference {
	switch n {
	case 0:
		return []Reference{}
	case 1:
		return []Reference{MacroParam{1, 1}}
	case 2:
		return []Reference{MacroParam{1, 1}, MacroParam{2, 2}}
	}
	panic(fmt.Sprintf("bad n: %d", n))
}

func Wrap(r Reference) Reference {
	name := fmt.Sprintf("W%s", r.Name())
	params := MapParams(len(r.Args()))
	refs := []Reference{}
	refs = append(refs, MacroPress)
	refs = append(refs, params...)
	refs = append(refs, Curry(r))
	refs = append(refs, MacroWait)
	refs = append(refs, MacroRelease)
	refs = append(refs, params...)
	refs = append(refs, Curry(r))
	AddMacro(Macro{
		Name:  name,
		Label: fmt.Sprintf("Wrap %s", r.Name()),
		Cells: len(r.Args()),
		Refs:  refs,
	})

	return Custom{name, MapToAny(r.Args())}
}
