package ref

import (
	"fmt"
	"strings"

	"keyboard/key/keys"
	. "keyboard/util"
)

type T = Ref

type Ref struct {
	Name   string
	Fields []any
}

func (x Ref) String() string {
	if len(x.Fields) == 0 {
		return x.Name
	}
	return x.Name + strings.Join(x.Args(), " ")
}

func (x Ref) Args() []string {
	return Map(x.Fields, ToString)
}

func (r Ref) Show() string {
	args := r.Args()
	if len(args) > 0 {
		return fmt.Sprintf("%s%s", r.Name, strings.Join(args, ""))
	}
	return r.Name
}

func Ref0(name string) Ref {
	return Ref{name, []any{}}
}

func Ref1(name string, a any) Ref {
	return Ref{name, []any{a}}
}

func Ref2(name string, a, b any) Ref {
	return Ref{name, []any{a, b}}
}

func RefN(name string, aa []any) Ref {
	return Ref{name, aa}
}

func Filled(name string, n int, aa ...any) Ref {
	switch n - len(aa) {
	case 2:
		aa = append(aa, keys.ZERO, keys.ZERO)
	case 1:
		aa = append(aa, keys.ZERO)
	case 0:
	default:
		panic("bad custom n")
	}

	if len(aa) > 2 {
		panic("too many custom n")
	}

	return Ref{name, aa}
}

func Compile(b Ref) string {
	args := b.Args()
	if len(args) > 0 {
		return fmt.Sprintf("&%s %s", b.Name, strings.Join(args, " "))
	}
	return "&" + b.Name
}

func Equal(a, b Ref) bool {
	return Compile(a) == Compile(b)
}
