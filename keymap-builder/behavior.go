package main

import (
	"fmt"
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

// lslxl: lslxl {
//     compatible = "zmk,behavior-hold-tap";
//     label = "LSLXL";
//     #binding-cells = <2>;
//     tapping-term-ms = <300>;
//     flavor = "balanced";
//     bindings = <&mo>, <&slxl>;
// };
