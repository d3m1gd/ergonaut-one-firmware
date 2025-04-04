package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

const Miss = "MISSING"
const Left Side = "left"
const Right Side = "right"

type Side string

func (s Side) Short() string {
	switch s {
	case Left:
		return "l"
	case Right:
		return "r"
	}
	panic("unhandled side: " + s)
}

type ToMacro struct {
	Level int
	Key   Key
}

type Params struct {
	ToMacros []ToMacro
}

type Key struct {
	Number      int
	RowCol      RC
	Action      string
	Other       string
	Constructor func(string, string) string
}

func (b Key) Behavior() string {
	return b.Constructor(b.Action, b.Other)
}

type RC struct {
	Side Side
	Row  int
	Col  int
}

func r(a, b int) RC {
	return RC{Right, a, b}
}

func l(a, b int) RC {
	return RC{Left, a, b}
}

func (rc RC) String() string {
	return fmt.Sprintf("%s%d%d", rc.Side.Short(), rc.Row, rc.Col)
}

func (rc RC) Pretty() string {
	return fmt.Sprintf("%s%d%d", strings.ToUpper(rc.Side.Short()), rc.Row, rc.Col)
}

func KpKp(a, b string) string {
	return fmt.Sprintf("&kpkp %s %s", b, a)
}

var l11 = Key{1, l(1, 1), "a", "b", KpKp}
var l12 = Key{2, l(1, 2), "a", "b", KpKp}
var l13 = Key{3, l(1, 3), "a", "b", KpKp}
var l14 = Key{4, l(1, 4), "a", "b", KpKp}
var l15 = Key{5, l(1, 5), "a", "b", KpKp}
var l16 = Key{6, l(1, 6), "a", "b", KpKp}
var r11 = Key{7, r(1, 1), "a", "b", KpKp}
var r12 = Key{8, r(1, 2), "a", "b", KpKp}
var r13 = Key{9, r(1, 3), "a", "b", KpKp}
var r14 = Key{10, r(1, 4), "a", "b", KpKp}
var r15 = Key{11, r(1, 5), "a", "b", KpKp}
var r16 = Key{12, r(1, 6), "a", "b", KpKp}
var l21 = Key{13, l(2, 1), "a", "b", KpKp}
var l22 = Key{14, l(2, 2), "a", "b", KpKp}
var l23 = Key{15, l(2, 3), "a", "b", KpKp}
var l24 = Key{16, l(2, 4), "a", "b", KpKp}
var l25 = Key{17, l(2, 5), "a", "b", KpKp}
var l26 = Key{18, l(2, 6), "a", "b", KpKp}
var r21 = Key{19, r(2, 1), "a", "b", KpKp}
var r22 = Key{20, r(2, 2), "a", "b", KpKp}
var r23 = Key{21, r(2, 3), "a", "b", KpKp}
var r24 = Key{22, r(2, 4), "a", "b", KpKp}
var r25 = Key{23, r(2, 5), "a", "b", KpKp}
var r26 = Key{24, r(2, 6), "a", "b", KpKp}
var l31 = Key{25, l(3, 1), "a", "b", KpKp}
var l32 = Key{26, l(3, 2), "a", "b", KpKp}
var l33 = Key{27, l(3, 3), "a", "b", KpKp}
var l34 = Key{28, l(3, 4), "a", "b", KpKp}
var l35 = Key{29, l(3, 5), "a", "b", KpKp}
var l36 = Key{30, l(3, 6), "a", "b", KpKp}
var r31 = Key{31, r(3, 1), "a", "b", KpKp}
var r32 = Key{32, r(3, 2), "M", "RG(M)", KpKp}
var r33 = Key{33, r(3, 3), "COMMA", "RG(COMMA)", KpKp}
var r34 = Key{34, r(3, 4), "DOT", "RG(DOT)", KpKp}
var r35 = Key{35, r(3, 5), "a", "b", KpKp}
var r36 = Key{36, r(3, 6), "a", "b", KpKp}
var l41 = Key{37, l(4, 1), "a", "b", KpKp}
var l42 = Key{38, l(4, 2), "a", "b", KpKp}
var l43 = Key{39, l(4, 3), "a", "b", KpKp}
var r41 = Key{40, r(4, 1), "a", "b", KpKp}
var r42 = Key{41, r(4, 2), "a", "b", KpKp}
var r43 = Key{42, r(4, 3), "a", "b", KpKp}

var all = []Key{}

func init() {
	all = append(all, r32, r33, r34)
}

func renderKeymap(path string, params Params) {
	t := must(template.ParseFiles(path + ".tmpl"))
	outFile := must(os.Create(path))
	defer outFile.Close()
	check(t.Execute(outFile, params))
}

func main() {
	params := Params{
		ToMacros: Map(all, func(k Key) ToMacro {
			return ToMacro{0, k}
		}),
	}

	renderKeymap("config/ergonaut_one.keymap", params)
	fmt.Println("good")
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Map[T, U any](s []T, f func(T) U) []U {
	us := make([]U, len(s))
	for i := range s {
		us[i] = f(s[i])
	}

	return us
}
