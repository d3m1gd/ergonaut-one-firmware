package main

import (
	"fmt"
	"strings"
)

type Reference interface {
	Args() []string
	Name() string
}

func CompileReference(b Reference) string {
	args := b.Args()
	if len(args) > 0 {
		return fmt.Sprintf("&%s %s", b.Name(), strings.Join(args, " "))
	}
	return "&" + b.Name()
}

func ShowReference(b Reference) string {
	args := b.Args()
	if len(args) > 0 {
		return fmt.Sprintf("%s%s", b.Name(), strings.Join(args, ""))
	}
	return b.Name()
}

func EqualReference(a, b Reference) bool {
	return CompileReference(a) == CompileReference(b)
}

var Trans = Custom0("trans")
var None = Custom0("none")

type Lt struct {
	Layer LayerIndex
	Tap   KeyCode
}

type To struct {
	Layer LayerIndex
}

type Mo struct {
	Layer LayerIndex
}

type Mt struct {
	Hold KeyCode
	Tap  KeyCode
}

func Kp(k KeyCode) Reference {
	return Custom1("kp", k)
}

type MKp struct {
	Tap KeyCode
}

type Rmt struct {
	Hold KeyCode
	Tap  KeyCode
}
