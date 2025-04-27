package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
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
	case []KeyMod:
		// mods = <(MOD_RSFT|MOD_LSFT)>;
		return fmt.Sprintf(`%s = <(%s)>`, x.Name, strings.Join(Map(v, asString), "|"))
	}
	panic(fmt.Sprintf("unknown device tree property: %T", x.Value))
}

type Behavior struct {
	Name  string
	Label string
	Cells int
	Type  string
	Refs  []Reference
	Props []DeviceTreeProperty
}

func (b Behavior) Equal(other Behavior) bool {
	eq := true
	eq = eq && b.Name == other.Name
	eq = eq && b.Label == other.Label
	eq = eq && b.Type == other.Type
	eq = eq && slices.EqualFunc(b.Refs, other.Refs, EqualReference)
	eq = eq && slices.EqualFunc(b.Props, other.Props, func(a, b DeviceTreeProperty) bool {
		return a.CompileProperty() == b.CompileProperty()
	})
	return eq
}

func (m Behavior) Bindings() string {
	return strings.Join(Map(m.Refs, func(r Reference) string { return "&" + r.Name() }), " ")
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
		panicif(!b.Equal(behaviors[i]))
		return
	}
	behaviors = append(behaviors, b)
}

func ModRef(key KeyCode, ref Reference) Reference {
	name := fmt.Sprintf("m%s", ref.Name())
	refs := []Reference{Kp{}, ref}
	AddBehavior(Behavior{
		Name:  name,
		Label: fmt.Sprintf("Mod %s", ref.Name()),
		Cells: BehaviorTypeHoldTap.Cells,
		Type:  BehaviorTypeHoldTap.Name,
		Refs:  refs,
		Props: []DeviceTreeProperty{
			{"flavor", "tap-preferred"},
			{"tapping-term-ms", 200},
		},
	})

	return CustomN(name, BehaviorTypeHoldTap.Cells-CountSlots(refs), key)
}

func MoTo(mo, to LayerIndex) Reference {
	name := "moto"
	refs := []Reference{Mo{}, To{}}
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

	return CustomN(name, BehaviorTypeHoldTap.Cells-CountSlots(refs), mo, to)
}

func MoX(name string, mo LayerIndex, x Reference) Reference {
	refs := []Reference{Mo{}, x}
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

	return CustomN(name, BehaviorTypeHoldTap.Cells-CountSlots(refs), mo)
}

func CountSlots(rr []Reference) int {
	return Reduce(rr, 0, func(acc int, r Reference) int {
		return acc + r.Slots()
	})
}

func ModMorph(name string, a, b Reference, mods []KeyMod) Reference {
	refs := []Reference{a, b}
	AddBehavior(Behavior{
		Name:  name,
		Label: "ModMorph " + name,
		Cells: BehaviorTypeModMorph.Cells,
		Type:  BehaviorTypeModMorph.Name,
		Refs:  refs,
		Props: []DeviceTreeProperty{
			{"mods", mods},
		},
	})

	return CustomN(name, BehaviorTypeModMorph.Cells-CountSlots(refs))
}

// mmMoveUnder: mmMoveUnder {
//     compatible = "zmk,behavior-mod-morph";
//     label = "mm Move Under";
//     bindings = <&to MOVER>, <&kp UNDERSCORE>;
//     // TODO
//
//     #binding-cells = <0>;
//     mods = <(MOD_RSFT|MOD_LSFT)>;
// };
//
