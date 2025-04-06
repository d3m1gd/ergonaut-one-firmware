package main

import (
	"fmt"
	"slices"
	"strings"
)

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
	return strings.Join(Map(m.Refs, CompileBehavior), " ")
}

func (m Macro) Equal(other Macro) bool {
	eq := true
	eq = eq && m.Name == other.Name
	eq = eq && m.Label == other.Label
	eq = eq && m.Cells == other.Cells
	eq = eq && slices.EqualFunc(m.Refs, other.Refs, func(a, b Reference) bool {
		return a.Reference() == b.Reference()
	})
	return eq
}

func BackspaceDeleteMacro() Reference {
	name := "bspcdel"
	AddMacro(Macro{
		Name:  name,
		Label: "Backspace Delete",
		Cells: 0,
		Refs:  []Reference{Kp{BSPC}, Kp{DEL}},
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

func XThenTransMacro(beh Reference, index LayerIndex, rc RC) Reference {
	layer := layers[index]
	name := fmt.Sprintf("xThenTrans%d", index)
	args := beh.Args()
	anyArgs := MapToAnyStatic(args, "MACRO_PLACEHOLDER")
	inner := Custom{beh.Name(), anyArgs}
	AddMacro(Macro{
		Name:  name,
		Label: fmt.Sprintf("X Then Trans %s", index),
		Cells: 1,
		Refs:  []Reference{MacroParam{1, 1}, inner, To{index}, layer[rc]},
	})

	return Custom{name, MapToAny(args)}
}

func XThenLayerMacro(r Reference, index LayerIndex) Reference {
	name := fmt.Sprintf("xThenLayer%d", index)
	args := r.Args()
	anyArgs := MapToAnyStatic(args, "MACRO_PLACEHOLDER")
	inner := Custom{r.Name(), anyArgs}
	AddMacro(Macro{
		Name:  name,
		Label: fmt.Sprintf("X Then Layer %s", index),
		Cells: 1,
		Refs:  []Reference{MacroParam{1, 1}, inner, To{index}},
	})

	return Custom{name, MapToAny(args)}
}
