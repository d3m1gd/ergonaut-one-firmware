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

func Custom0(name string) Custom {
	return Custom{name, []any{}}
}

func Custom1(name string, a any) Custom {
	return Custom{name, []any{a}}
}

func Custom2(name string, a, b any) Custom {
	return Custom{name, []any{a, b}}
}
