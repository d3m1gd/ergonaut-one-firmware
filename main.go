package main

import (
	"cmp"
	"fmt"
	"iter"
	"maps"
	"os"
	"slices"
	"strings"
	"text/template"
)

type LevelName string

var LevelNames = []LevelName{
	"BASE",
}

const Miss = "MISSING"
const Left Side = "left"
const Right Side = "right"
const L = Left
const R = Right

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

type Layer struct {
	Name string
	Keys []Key
}

type ToMacro struct {
	Level int
	RC    RC
	Key   Key
}

type Params struct {
	ToBaseAnd []ToMacro
}

type Key struct {
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

func (rc RC) Offset() int {
	offset := 0
	if rc.Side == Right {
		offset = 6
		if rc.Row == 4 {
			offset = 3
		}
	}

	return offset
}

func (rc RC) Serial() int {
	return rc.Offset() + 12*(rc.Row-1) + rc.Col
}

func (rc RC) Less(other RC) int {
	return cmp.Compare(rc.Serial(), other.Serial())
}

func KpKp(a, b string) string {
	return fmt.Sprintf("&kpkp %s %s", b, a)
}

func Trans(_, _ string) string {
	return "&trans"
}

var all = []Key{}

func InitAllTrans() []Key {
	keys := make([]Key, 42)
	for i := range keys {
		keys[i] = Key{"", "", Trans}
	}

	return keys
}

var layers = map[string]map[RC]Key{}

func init() {
	layers["BASE"] = map[RC]Key{}
	layers["BASE"][l(1, 1)] = Key{"a", "b", KpKp}
	layers["BASE"][l(1, 2)] = Key{"a", "b", KpKp}
	layers["BASE"][l(1, 3)] = Key{"a", "b", KpKp}
	layers["BASE"][l(1, 4)] = Key{"a", "b", KpKp}
	layers["BASE"][l(1, 5)] = Key{"a", "b", KpKp}
	layers["BASE"][l(1, 6)] = Key{"a", "b", KpKp}
	layers["BASE"][r(1, 1)] = Key{"a", "b", KpKp}
	layers["BASE"][r(1, 2)] = Key{"a", "b", KpKp}
	layers["BASE"][r(1, 3)] = Key{"a", "b", KpKp}
	layers["BASE"][r(1, 4)] = Key{"a", "b", KpKp}
	layers["BASE"][r(1, 5)] = Key{"a", "b", KpKp}
	layers["BASE"][r(1, 6)] = Key{"a", "b", KpKp}
	layers["BASE"][l(2, 1)] = Key{"a", "b", KpKp}
	layers["BASE"][l(2, 2)] = Key{"a", "b", KpKp}
	layers["BASE"][l(2, 3)] = Key{"a", "b", KpKp}
	layers["BASE"][l(2, 4)] = Key{"a", "b", KpKp}
	layers["BASE"][l(2, 5)] = Key{"a", "b", KpKp}
	layers["BASE"][l(2, 6)] = Key{"a", "b", KpKp}
	layers["BASE"][r(2, 1)] = Key{"a", "b", KpKp}
	layers["BASE"][r(2, 2)] = Key{"a", "b", KpKp}
	layers["BASE"][r(2, 3)] = Key{"a", "b", KpKp}
	layers["BASE"][r(2, 4)] = Key{"a", "b", KpKp}
	layers["BASE"][r(2, 5)] = Key{"a", "b", KpKp}
	layers["BASE"][r(2, 6)] = Key{"a", "b", KpKp}
	layers["BASE"][l(3, 1)] = Key{"a", "b", KpKp}
	layers["BASE"][l(3, 2)] = Key{"a", "b", KpKp}
	layers["BASE"][l(3, 3)] = Key{"a", "b", KpKp}
	layers["BASE"][l(3, 4)] = Key{"a", "b", KpKp}
	layers["BASE"][l(3, 5)] = Key{"a", "b", KpKp}
	layers["BASE"][l(3, 6)] = Key{"a", "b", KpKp}
	layers["BASE"][r(3, 1)] = Key{"a", "b", KpKp}
	layers["BASE"][r(3, 2)] = Key{"M", "RG(M)", KpKp}
	layers["BASE"][r(3, 3)] = Key{"COMMA", "RG(COMMA)", KpKp}
	layers["BASE"][r(3, 4)] = Key{"DOT", "RG(DOT)", KpKp}
	layers["BASE"][r(3, 5)] = Key{"a", "b", KpKp}
	layers["BASE"][r(3, 6)] = Key{"a", "b", KpKp}
	layers["BASE"][l(4, 1)] = Key{"a", "b", KpKp}
	layers["BASE"][l(4, 2)] = Key{"a", "b", KpKp}
	layers["BASE"][l(4, 3)] = Key{"a", "b", KpKp}
	layers["BASE"][r(4, 1)] = Key{"a", "b", KpKp}
	layers["BASE"][r(4, 2)] = Key{"a", "b", KpKp}
	layers["BASE"][r(4, 3)] = Key{"a", "b", KpKp}

	// BASE.Keys[r(3, 2).Serial()]={r(3, 2), "M", "RG(M)", KpKp}

	// all = append(all, r32, r33, r34)
}

func renderKeymap(path string, params Params) {
	t := must(template.ParseFiles(path + ".tmpl"))
	outFile := must(os.Create(path))
	defer outFile.Close()
	check(t.Execute(outFile, params))
}

func MapToBaseAnd(seq iter.Seq2[RC, Key]) iter.Seq[ToMacro] {
	return func(yield func(ToMacro) bool) {
		for rc, key := range seq {
			tm := ToMacro{0, rc, key}
			if !yield(tm) {
				return
			}
		}
	}
}

type Lesser[K any] interface {
	comparable
	Less(K) int
}

func SortedMap[K Lesser[K], V any](m map[K]V) iter.Seq2[K, V] {
	keys := slices.Collect(maps.Keys(m))
	slices.SortFunc(keys, func(a, b K) int {
		return a.Less(b)
	})
	return func(yield func(K, V) bool) {
		for _, key := range keys {
			if !yield(key, m[key]) {
				return
			}
		}
	}
}

func main() {
	params := Params{
		ToBaseAnd: slices.Collect(MapToBaseAnd(SortedMap(layers["BASE"]))),
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
