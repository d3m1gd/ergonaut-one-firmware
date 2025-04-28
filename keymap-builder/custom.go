package main

import (
	"fmt"
	"strings"
)

type Custom struct {
	Label  string
	Fields []any
}

func (x Custom) Reference() string {
	if len(x.Fields) == 0 {
		return fmt.Sprintf("&%s", x.Name())
	}
	return fmt.Sprintf("&%s %s", x.Name(), strings.Join(x.Args(), " "))
}

func (x Custom) Name() string {
	return x.Label
}

func (x Custom) Args() []string {
	return Map(x.Fields, toString)
}

func (x Custom) Slots() int {
	return 0
}

func Custom0(name string) Custom {
	return Custom{name, []any{}}
}

func Custom1(name string, a any) Custom {
	return Custom{name, []any{a}}
}

func Custom2(name string, a, b any) Custom {
	return Custom{name, []any{a, b}}
}

func CustomN(name string, n int, aa ...any) Custom {
	switch n - len(aa) {
	case 2:
		aa = append(aa, ZERO, ZERO)
	case 1:
		aa = append(aa, ZERO)
	case 0:
	default:
		panic("bad custom n")
	}

	if len(aa) > 2 {
		panic("too many custom n")
	}

	return Custom{name, aa}
}
