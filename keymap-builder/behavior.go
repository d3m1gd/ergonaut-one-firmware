package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	. "keyboard/key"
	. "keyboard/util"
)

type Number int

func (n Number) String() string {
	return strconv.Itoa(int(n))
}

var ZERO = Number(0)

type BehaviorType struct {
	Name  string
	Cells int
}

var BehaviorTypeHoldTap = BehaviorType{"behavior-hold-tap", 2}
var BehaviorTypeStickyKey = BehaviorType{"behavior-sticky-key", 1}
var BehaviorTypeModMorph = BehaviorType{"behavior-mod-morph", 0}

type DeviceTreeProperty struct {
	Name  string
	Value any
}

func (x DeviceTreeProperty) CompileProperty() string {
	switch v := x.Value.(type) {
	case int:
		return fmt.Sprintf("%s = <%d>", x.Name, v)
	case string:
		return fmt.Sprintf(`%s = "%s"`, x.Name, v)
	case []Mod:
		// mods = <(MOD_RSFT|MOD_LSFT)>;
		return fmt.Sprintf(`%s = <(%s)>`, x.Name, strings.Join(Map(v, AsString), "|"))
	}
	panic(fmt.Sprintf("unknown device tree property: %T", x.Value))
}

type Behavior struct {
	Name  string
	Label string
	Cells int
	Type  string
	Refs  []Ref
	Props []DeviceTreeProperty
}

func (b Behavior) Equal(other Behavior) bool {
	eq := true
	eq = eq && b.Name == other.Name
	eq = eq && b.Label == other.Label
	eq = eq && b.Type == other.Type
	eq = eq && slices.EqualFunc(b.Refs, other.Refs, EqualRef)
	eq = eq && slices.EqualFunc(b.Props, other.Props, func(a, b DeviceTreeProperty) bool {
		return a.CompileProperty() == b.CompileProperty()
	})
	return eq
}

func (m Behavior) Bindings() string {
	return strings.Join(Map(m.Refs, CompileRef), " ")
}

func (m Behavior) Properties() []string {
	return Map(m.Props, CompileProperty)
}

func CompileProperty(p DeviceTreeProperty) string {
	return p.CompileProperty()
}

func AddBehavior(b Behavior) {
	i := slices.IndexFunc(behaviors, func(other Behavior) bool {
		return b.Name == other.Name
	})
	if i != -1 {
		Panicif(!b.Equal(behaviors[i]))
		return
	}
	behaviors = append(behaviors, b)
}

func KpKp(a, b Key) Ref {
	name := "kpkp"
	_ = TapNoRepeat(a) // instantiate macro
	AddBehavior(Behavior{
		Name:  name,
		Label: "kpkp",
		Cells: BehaviorTypeHoldTap.Cells,
		Type:  BehaviorTypeHoldTap.Name,
		Refs:  []Ref{Ref0("TapNoRepeat"), Ref0("kp")},
		Props: []DeviceTreeProperty{
			{"flavor", "tap-preferred"},
			{"tapping-term-ms", 200},
			{"quick-tap-ms", 200},
		},
	})

	return RefN(name, BehaviorTypeHoldTap.Cells, a, b)
}

func MoTo(mo, to LayerIndex) Ref {
	refs := []Ref{Ref0("mo"), Ref0("to")}
	name := "moto"
	AddBehavior(Behavior{
		Name:  name,
		Label: "Momentary/To",
		Type:  BehaviorTypeHoldTap.Name,
		Cells: BehaviorTypeHoldTap.Cells,
		Refs:  refs,
		Props: []DeviceTreeProperty{
			{"flavor", "balanced"},
			{"tapping-term-ms", 300},
		},
	})

	return RefN(name, BehaviorTypeHoldTap.Cells, mo, to)
}

func MoX(mo LayerIndex, x Ref) Ref {
	refs := []Ref{Ref0("mo"), x}
	name := "mo" + ShowReference(x)
	AddBehavior(Behavior{
		Name:  name,
		Label: "Momentary " + name,
		Type:  BehaviorTypeHoldTap.Name,
		Cells: BehaviorTypeHoldTap.Cells,
		Refs:  refs,
		Props: []DeviceTreeProperty{
			{"flavor", "balanced"},
			{"tapping-term-ms", 300},
		},
	})

	return RefN(name, BehaviorTypeHoldTap.Cells, mo)
}

func ModX(mod Key, x Ref) Ref {
	name := "m" + ShowReference(x)
	AddBehavior(Behavior{
		Name:  name,
		Label: "Mod " + x.Name(),
		Type:  BehaviorTypeHoldTap.Name,
		Cells: BehaviorTypeHoldTap.Cells,
		Refs:  []Ref{Ref0("kp"), x},
		Props: []DeviceTreeProperty{
			{"flavor", "tap-preferred"},
			{"tapping-term-ms", 200},
			{"quick-tap-ms", 200},
		},
	})

	return RefN(name, BehaviorTypeHoldTap.Cells, mod)
}

func ModMorph(a, b Ref, mods []Mod, keep []Mod) Ref {
	refs := []Ref{a, b}
	name := "mm" + ShowReference(a) + ShowReference(b)
	props := []DeviceTreeProperty{
		{"mods", mods},
	}
	if len(keep) > 0 {
		props = append(props, DeviceTreeProperty{
			"keep-mods", keep,
		})
	}

	AddBehavior(Behavior{
		Name:  name,
		Label: "ModMorph " + a.Name() + " " + b.Name(),
		Cells: BehaviorTypeModMorph.Cells,
		Type:  BehaviorTypeModMorph.Name,
		Refs:  refs,
		Props: props,
	})

	return RefN(name, BehaviorTypeModMorph.Cells)
}
