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

func EqualReference(a, b Reference) bool {
	return CompileReference(a) == CompileReference(b)
}

type Trans struct{}

type None struct{}

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

type Kp struct {
	Tap KeyCode
}

type KpKp struct {
	Hold KeyCode
	Tap  KeyCode
}

type MKp struct {
	Tap KeyCode
}

type Rmt struct {
	Hold KeyCode
	Tap  KeyCode
}
