package main

import (
	"fmt"
	"slices"
	"strconv"
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
	return strings.Join(Map(m.Refs, CompileReference), " ")
}

func (m Macro) Equal(other Macro) bool {
	eq := true
	eq = eq && m.Name == other.Name
	eq = eq && m.Label == other.Label
	eq = eq && m.Cells == other.Cells
	eq = eq && slices.EqualFunc(m.Refs, other.Refs, func(a, b Reference) bool {
		return CompileReference(a) == CompileReference(b)
	})
	return eq
}

func BackspaceDelete() Reference {
	name := "bspcdel"
	AddMacro(Macro{
		Name:  name,
		Label: "Backspace Delete",
		Cells: 0,
		Refs:  []Reference{Kp{BSPC}, Kp{DEL}},
	})

	return Custom0(name)
}

func Parens() Reference {
	return OpenCloseMacro("parens", Kp{LPAR}, Kp{RPAR})
}

func Brackets() Reference {
	return OpenCloseMacro("brackets", Kp{LBKT}, Kp{RBKT})
}

func Curlies() Reference {
	return OpenCloseMacro("curlies", Kp{LBRC}, Kp{RBRC})
}

func OpenCloseMacro(name string, left, right Kp) Reference {
	AddMacro(Macro{
		Name:  name,
		Label: fmt.Sprintf("OpenClose %s", name),
		Cells: 0,
		Refs:  []Reference{left, right, Kp{LEFT}, To{PARENS}},
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

type MacroStatePressRelease int

const (
	MacroStatePress MacroStatePressRelease = iota
	MacroStateRelease
	MacroStateWait
)

var MacroPress = MacroStateBase{MacroStatePress}
var MacroRelease = MacroStateBase{MacroStateRelease}
var MacroWait = MacroStateBase{MacroStateWait}

type MacroStateBase struct {
	press MacroStatePressRelease
}

func (mp MacroStateBase) Reference() string {
	return "&" + mp.Name()
}

func (mp MacroStateBase) Name() string {
	switch mp.press {
	case MacroStatePress:
		return "macro_press"
	case MacroStateRelease:
		return "macro_release"
	case MacroStateWait:
		return "macro_pause_for_release"
	}
	panic("unhandled macro press state: " + strconv.Itoa(int(mp.press)))
}

func (mp MacroStateBase) Args() []string {
	return []string{}
}

func XThenTrans(r Reference, index LayerIndex, rc RC) Reference {
	layer := layers[index]
	name := fmt.Sprintf("xThenTrans%d", index)
	args := r.Args()
	inner := Custom{r.Name(), MapToAnyStatic(args, "MACRO_PLACEHOLDER")}
	AddMacro(Macro{
		Name:  name,
		Label: fmt.Sprintf("X Then Trans %s", index),
		Cells: 1,
		Refs:  []Reference{MacroParam{1, 1}, inner, To{index}, layer[rc]},
	})

	return Custom{name, MapToAny(args)}
}

func XThenLayer(r Reference, index LayerIndex) Reference {
	name := fmt.Sprintf("xThenLayer%d", index)
	args := r.Args()
	inner := Custom{r.Name(), MapToAnyStatic(args, "MACRO_PLACEHOLDER")}
	AddMacro(Macro{
		Name:  name,
		Label: fmt.Sprintf("X Then Layer %s", index),
		Cells: 1,
		Refs:  []Reference{MacroParam{1, 1}, inner, To{index}},
	})

	return Custom{name, MapToAny(args)}
}
