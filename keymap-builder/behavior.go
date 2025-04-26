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

type BehaviorType string

var BehaviorTypeHoldTap BehaviorType = "behavior-hold-tap"
var BehaviorTypeStickyKey BehaviorType = "behavior-sticky-key"
var BehaviorTypeModMorph BehaviorType = "behavior-mod-morph"

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
	}
	panic(fmt.Sprintf("unknown type: %T", x.Value))
}

type Behavior struct {
	Name  string
	Label string
	Cells int
	Type  BehaviorType
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
	AddBehavior(Behavior{
		Name:  name,
		Label: fmt.Sprintf("Mod %s", ref.Name()),
		Cells: 2,
		Type:  BehaviorTypeHoldTap,
		Refs:  []Reference{Kp{}, ref},
		Props: []DeviceTreeProperty{
			{"flavor", "tap-preferred"},
			{"tapping-term-ms", 200},
		},
	})

	return Custom2(name, key, ZERO)
}

func MoTo(mo, to LayerIndex) Reference {
	name := "moto"
	AddBehavior(Behavior{
		Name:  name,
		Label: "Momentary/To",
		Type:  BehaviorTypeHoldTap,
		Cells: 2,
		Refs:  []Reference{Mo{}, To{}},
		Props: []DeviceTreeProperty{
			{"flavor", "balanced"},
			{"tapping-term-ms", 300},
		},
	})

	return Custom2(name, mo, to)
}

func MoX(name string, mo LayerIndex, x Reference) Reference {
	AddBehavior(Behavior{
		Name:  name,
		Label: "Momentary " + name,
		Type:  BehaviorTypeHoldTap,
		Cells: 1,
		Refs:  []Reference{Mo{}, x},
		Props: []DeviceTreeProperty{
			{"flavor", "balanced"},
			{"tapping-term-ms", 300},
		},
	})

	return Custom2(name, mo, to)
}

func ModMorph(name string, l LayerIndex, k Reference, mods []KeyMod) Reference {
	AddBehavior(Behavior{
		Name:  name,
		Label: "ModMorph " + name,
		Type:  BehaviorTypeModMorph,
		Cells: 1,
		Refs:  []Reference{To{}, k},
		Props: []DeviceTreeProperty{
			{"mods", mods},
			{"tapping-term-ms", 300},
		},
	})

	return Custom1(name, l)
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
