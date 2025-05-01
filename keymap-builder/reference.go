package main

import (
	"fmt"
	"strings"

	. "keyboard/key"
)

type Ref interface {
	Args() []string
	Name() string
}

func CompileRef(b Ref) string {
	args := b.Args()
	if len(args) > 0 {
		return fmt.Sprintf("&%s %s", b.Name(), strings.Join(args, " "))
	}
	return "&" + b.Name()
}

func ShowReference(b Ref) string {
	args := b.Args()
	if len(args) > 0 {
		return fmt.Sprintf("%s%s", b.Name(), strings.Join(args, ""))
	}
	return b.Name()
}

func EqualRef(a, b Ref) bool {
	return CompileRef(a) == CompileRef(b)
}

var Trans = Ref0("trans")
var None = Ref0("none")
var CapsWord = Ref0("caps_word")

func Lt(layer LayerIndex, tap Key) Ref {
	return Ref2("lt", layer, tap)
}

func To(layer LayerIndex) Ref {
	return Ref1("to", layer)
}

func Mo(layer LayerIndex) Ref {
	return Ref1("mo", layer)
}

func Mt(mod, tap Key) Ref {
	return Ref2("mt", mod, tap)
}

func Kp(k Key) Ref {
	return Ref1("kp", k)
}

func MKp(tap Key) Ref {
	return Ref1("mkp", tap)
}

func Rmt(mod, tap Key) Ref {
	return Ref2("rmt", mod, tap)
}
