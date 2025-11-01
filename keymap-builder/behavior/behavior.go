package behavior

import (
	"cmp"
	"crypto/sha256"
	"fmt"
	"maps"
	"slices"
	"strings"

	"keyboard/key"
	"keyboard/ref"
	"keyboard/util"
	"keyboard/util/indenter"
)

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
			return x.Name
		} else {
			return fmt.Sprintf("/delete-property/ %s", x.Name)
		}
	case []int:
		return fmt.Sprintf("%s = <%s>", x.Name, strings.Join(util.Map(v, util.ToString), " "))
	case []key.Mod:
		return fmt.Sprintf("%s = <(%s)>", x.Name, strings.Join(util.Map(v, util.AsString), "|"))
	}
	panic(fmt.Sprintf("unknown device tree property: %T", x.Value))
}

type Behavior struct {
	Name  string
	Type  Type
	Refs  []ref.T
	Props Props
}

func (b Behavior) Equal(other Behavior) bool {
	eq := true
	eq = eq && b.Name == other.Name
	eq = eq && b.Type == other.Type
	eq = eq && slices.EqualFunc(b.Refs, other.Refs, ref.Equal)
	eq = eq && b.Props.Equal(other.Props)
	return eq
}

func (b Behavior) Bindings() string {
	return strings.Join(util.Map(b.Refs, ref.Compile), " ")
}

func (b Behavior) Properties() []string {
	keys := slices.Collect(maps.Keys(b.Props))
	slices.Sort(keys)
	return util.Map(keys, func(k string) string { return Prop{k, b.Props[k]}.Compile() })
}

func Add(b Behavior) ref.T {
	for _, r := range b.Refs {
		util.Panicif(slices.ContainsFunc(r.Fields, func(a any) bool { return a == nil }))
	}

	cells := b.Type.Cells
	args := []any{}

	switch cells {
	case 2:
		if len(b.Refs) == 2 {
			stripped := b.Refs[0].StripN(1) // todo 1
			if len(stripped) > 0 {
				args = append(args, stripped...)
			} else {
				args = append(args, key.ZERO)
			}
			stripped = b.Refs[1].StripN(1) // todo 1
			if len(stripped) > 0 {
				args = append(args, stripped...)
			} else {
				args = append(args, key.ZERO)
			}
		}
	case 1:
		if len(b.Refs) >= 1 {
			stripped := b.Refs[0].StripN(1) // todo 1
			if len(stripped) > 0 {
				args = append(args, stripped...)
			} else {
				args = append(args, key.ZERO)
			}
		}
	case 0:
	}

	for _, r := range b.Refs {
		b.Name += r.Show()
	}

	if len(b.Name) > 20 {
		sum := sha256.Sum256([]byte(b.Name))
		b.Name = string(fmt.Sprintf("%x", sum[:10]))
	}

	i := slices.IndexFunc(behaviors, func(other Behavior) bool {
		return b.Name == other.Name
	})
	if i != -1 {
		util.Panicif(!b.Equal(behaviors[i]))
		return ref.RefN(b.Name, args)
	}
	behaviors = append(behaviors, b)

	return ref.RefN(b.Name, args)
}

func Render() []Behavior {
	slices.SortFunc(behaviors, func(a, b Behavior) int { return cmp.Compare(a.Name, b.Name) })
	return behaviors
}

func (b Behavior) Compile(indent, level int) string {
	ir := indenter.New(indent)

	ir.Sprintf(0, "\n")
	ir.Sprintf(level, "%s: %s {\n", b.Name, b.Name)
	ir.Sprintf(level+1, "compatible = \"zmk,%s\";\n", b.Type.Name)
	ir.Sprintf(level+1, "#binding-cells = <%d>;\n", b.Type.Cells)
	if len(b.Bindings()) > 0 {
		ir.Sprintf(level+1, "bindings = <%s>;\n", b.Bindings())
	}
	for _, p := range b.Properties() {
		ir.Sprintf(level+1, "%s;\n", p)
	}
	ir.Sprintf(level, "};\n")
	return ir.String()
}
