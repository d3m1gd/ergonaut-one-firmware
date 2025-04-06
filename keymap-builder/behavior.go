package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type BehaviorType struct {
	Name  string
	Cells int
}

var BehaviorTypeHoldTap = BehaviorType{"behavior-hold-tap", 2}
var BehaviorTypeStickyKey = BehaviorType{"behavior-sticky-key", 1}

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
	Type  BehaviorType
	Refs  []Reference
	Props []DeviceTreeProperty
}

func (b Behavior) Equal(other Behavior) bool {
	eq := true
	eq = eq && b.Name == other.Name
	eq = eq && b.Label == other.Label
	eq = eq && b.Type == other.Type
	eq = eq && slices.EqualFunc(b.Refs, other.Refs, func(a, b Reference) bool {
		return a.Reference() == b.Reference()
	})
	eq = eq && slices.EqualFunc(b.Props, other.Props, func(a, b DeviceTreeProperty) bool {
		return a.CompileProperty() == b.CompileProperty()
	})
	return eq
}

func (m Behavior) Cells() string {
	return strconv.Itoa(m.Type.Cells)
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
		Type:  BehaviorTypeHoldTap,
		Refs:  []Reference{Kp{}, ref},
		Props: []DeviceTreeProperty{
			{"flavor", "tap-preferred"},
			{"tapping-term-ms", 200},
		},
	})

	return Custom1(name, key)
}

func MoTo(mo, to LayerIndex) Reference {
	name := "moto"
	AddBehavior(Behavior{
		Name:  name,
		Label: "Momentary/To",
		Type:  BehaviorTypeHoldTap,
		Refs:  []Reference{Mo{}, To{}},
		Props: []DeviceTreeProperty{
			{"flavor", "balanced"},
			{"tapping-term-ms", 300},
		},
	})

	return Custom2(name, mo, to)
}
