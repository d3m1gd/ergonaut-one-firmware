package behavior

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strings"

	"keyboard/key"
	"keyboard/key/keys"
	"keyboard/ref"
	. "keyboard/util"
)

type T = Behavior

type Type struct {
	Name  string
	Cells int
}

var behaviors []Behavior

var (
	TypeHoldTap     = Type{"behavior-hold-tap", 2}
	TypeStickyKey   = Type{"behavior-sticky-key", 1}
	TypeModMorph    = Type{"behavior-mod-morph", 0}
	TypeToggleLayer = Type{"behavior-toggle-layer", 1}
)

type Props map[string]any

func (p Props) Equal(other Props) bool {
	if len(p) != len(other) {
		return false
	}

	for k, v := range p {
		switch vs := v.(type) {
		case []int:
			switch os := other[k].(type) {
			case []int:
				return slices.Equal(vs, os)
			}
		}

		if v != other[k] {
			return false
		}
	}

	return true
}

type Prop struct {
	Name  string
	Value any
}

func (x Prop) Compile() string {
	switch v := x.Value.(type) {
	case int:
		return fmt.Sprintf("%s = <%d>", x.Name, v)
	case string:
		return fmt.Sprintf(`%s = "%s"`, x.Name, v)
	case bool:
		if v {
			return fmt.Sprintf(`%s`, x.Name)
		} else {
			return fmt.Sprintf(`/delete-property/ %s`, x.Name)
		}
	case []int:
		return fmt.Sprintf(`%s = <%s>`, x.Name, strings.Join(Map(v, ToString), " "))
	case []key.Mod:
		return fmt.Sprintf(`%s = <(%s)>`, x.Name, strings.Join(Map(v, AsString), "|"))
	}
	panic(fmt.Sprintf("unknown device tree property: %T", x.Value))
}

type Behavior struct {
	Name  string
	Label string
	Type  Type
	Refs  []ref.T
	Props Props
}

func (b Behavior) Equal(other Behavior) bool {
	eq := true
	eq = eq && b.Name == other.Name
	eq = eq && b.Label == other.Label
	eq = eq && b.Type == other.Type
	eq = eq && slices.EqualFunc(b.Refs, other.Refs, ref.Equal)
	eq = eq && b.Props.Equal(other.Props)
	return eq
}

func (m Behavior) Bindings() string {
	return strings.Join(Map(m.Refs, ref.Compile), " ")
}

func (m Behavior) Properties() []string {
	keys := slices.Collect(maps.Keys(m.Props))
	slices.Sort(keys)
	return Map(keys, func(k string) string { return Prop{k, m.Props[k]}.Compile() })
}

func Add(b Behavior) {
	i := slices.IndexFunc(behaviors, func(other Behavior) bool {
		return b.Name == other.Name
	})
	if i != -1 {
		Panicif(!b.Equal(behaviors[i]))
		return
	}
	behaviors = append(behaviors, b)
}

func AddX(args []any, b Behavior) ref.T {
	Add(b)
	return ref.Filled(b.Name, b.Type.Cells, args...)
}

func AddY(b Behavior) ref.T {
	cells := b.Type.Cells
	args := []any{}

	for i, r := range b.Refs {
		Panicif(slices.ContainsFunc(r.Fields, func(a any) bool { return a == nil }))

		if r.Filled() > 0 && cells > 0 {
			stripped := r.StripN(cells)
			b.Refs[i] = r
			args = append(args, stripped...)
			cells -= len(stripped)
		}

		b.Name += r.Show()
		b.Label += r.Show()
	}

	Add(b)

	for len(args) < b.Type.Cells {
		args = append(args, keys.ZERO)
	}

	return ref.RefN(b.Name, args)
}

func Render() []Behavior {
	slices.SortFunc(behaviors, func(a, b Behavior) int { return cmp.Compare(a.Name, b.Name) })
	return behaviors
}
