package ref

import (
	"slices"
	"strings"

	. "keyboard/util"
)

type T = Ref

type Ref struct {
	Name   string
	Fields []any
}

// todo delme
func (x Ref) Strip() Ref {
	Panicif(len(x.Fields) > 2)
	return Ref{
		Name:   x.Name,
		Fields: Map(x.Fields, func(a any) any { return nil }),
	}
}

// todo rename
func (x Ref) Strip2() Ref {
	Panicif(len(x.Fields) != 2)
	return Ref{
		Name:   x.Name,
		Fields: Map(x.Fields, func(a any) any { return nil }),
	}
}

func (x Ref) Filled() int {
	return Reduce(x.Fields, 0, func(acc int, a any) int {
		if a == nil {
			return acc
		} else {
			return acc + 1
		}
	})
}

func (x *Ref) StripN(n int) []any {
	Panicif(len(x.Fields) > 2)
	Panicif(n > 2)
	Panicif(slices.ContainsFunc(x.Fields, func(f any) bool { return f == nil }))

	switch n {
	case 2:
		fields := x.Fields
		x.Fields = Map(fields, func(f any) any { return nil })
		return fields
	case 1:
		if len(x.Fields) > 0 {
			f := x.Fields[len(x.Fields)-1]
			x.Fields[len(x.Fields)-1] = nil
			return []any{f}
		}
	}

	return []any{}
}

func (x Ref) String() string {
	b := strings.Builder{}
	b.WriteString(x.Name)
	for _, field := range x.Fields {
		if field != nil {
			b.WriteString(" ")
			b.WriteString(ToString(field))
		}
	}
	return b.String()
}

func (x Ref) Compile() string {
	return "&" + x.String()
}

func (x Ref) Args() []string {
	return Map(x.Fields, ToString)
}

func (r Ref) Show() string {
	s := r.String()
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}

func Ref0(name string) Ref {
	return RefN(name, []any{})
}

func Ref1(name string, a any) Ref {
	return RefN(name, []any{a})
}

func Ref2(name string, a, b any) Ref {
	return RefN(name, []any{a, b})
}

func RefN(name string, aa []any) Ref {
	return Ref{name, aa}
}

func Filled(name string, n int, aa ...any) Ref {
	Panicif(n != len(aa))
	return RefN(name, aa)
}

func Compile(r Ref) string {
	return r.Compile()
}

func Equal(a, b Ref) bool {
	return Compile(a) == Compile(b)
}
