//go:generate stringer -type=LayerIndex

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

	"github.com/BooleanCat/go-functional/v2/it"
)

type Behavior interface {
	Behavior() string
}

type Layer map[RC]Behavior
type LayerSeq = iter.Seq2[RC, Behavior]
type LayerName string
type LayerIndex int

const (
	BASE LayerIndex = iota
	MAXLAYERINDEX
)

func (l Layer) Render() []string {
	widths := slices.Repeat([]int{0}, 12)
	cells := make([][]string, 4)
	for i := range cells {
		cells[i] = slices.Repeat([]string{""}, 12)
	}
	for n := range 42 {
		n := n + 1
		rc := RCFrom(n)
		b := l[rc]
		rendered := b.Behavior()
		row := rc.Row - 1
		col := rc.Col - 1
		if rc.Side == Right {
			col += 6
		} else {
			if rc.Row == 4 {
				col += 3
			}
		}
		cells[row][col] = rendered
		widths[col] = max(widths[col], len(rendered))
	}

	for _, row := range cells {
		for j := range row {
			row[j] = fmt.Sprintf("%-*s", widths[j], row[j])
		}
	}

	return slices.Collect(it.Map(slices.Values(cells), func(ss []string) string {
		s := strings.Join(ss, "  ")
		s = strings.TrimRight(s, " ")
		return s
	}))
}

func (ln LayerName) Less(other LayerName) int {
	return cmp.Compare(ln, other)
}

var LayerNames = []LayerName{
	"BASE",
}

type RenderedLayer struct {
	Name string
	Rows []string
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

type ToMacro struct {
	Level    int
	RC       RC
	Behavior string
}

type Params struct {
	ToBaseAnd []ToMacro
	Layers    []RenderedLayer
}

type KpKp struct {
	Tap  string
	Hold string
}

func (x KpKp) Behavior() string {
	return fmt.Sprintf("&kpkp %s %s", x.Hold, x.Tap)
}

type Trans struct{}

func (_ Trans) Behavior() string {
	return "&trans"
}

type RC struct {
	Side Side
	Row  int
	Col  int
}

func RCFrom(n int) RC {
	row := (n-1)/12 + 1
	col := (n-1)%12 + 1
	side := Left
	if col > 6 {
		side = Right
		col -= 6
	}
	if row == 4 {
		if col > 3 {
			side = Right
			col -= 3
		}
	}

	return RC{side, row, col}
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

func InitTrans() Layer {
	layer := Layer{}
	for i := range 42 {
		layer[RCFrom(i+1)] = Trans{}
	}

	return layer
}

var layers = make([]Layer, MAXLAYERINDEX)

func init() {
	layers[BASE] = InitTrans()
	// layers["BASEXX"][l(1, 1)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(1, 2)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(1, 3)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(1, 4)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(1, 5)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(1, 6)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(1, 1)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(1, 2)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(1, 3)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(1, 4)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(1, 5)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(1, 6)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(2, 1)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(2, 2)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(2, 3)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(2, 4)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(2, 5)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(2, 6)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(2, 1)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(2, 2)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(2, 3)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(2, 4)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(2, 5)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(2, 6)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(3, 1)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(3, 2)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(3, 3)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(3, 4)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(3, 5)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(3, 6)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(3, 1)] = Key{"a", "b", KpKp}
	layers[BASE][r(3, 2)] = KpKp{"M", "RG(M)"}
	layers[BASE][r(3, 3)] = KpKp{"COMMA", "RG(COMMA)"}
	layers[BASE][r(3, 4)] = KpKp{"DOT", "RG(DOT)"}
	// layers["BASEXX"][r(3, 5)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(3, 6)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(4, 1)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(4, 2)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][l(4, 3)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(4, 1)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(4, 2)] = Key{"a", "b", KpKp}
	// layers["BASEXX"][r(4, 3)] = Key{"a", "b", KpKp}
}

func renderKeymap(path string, params Params) {
	t := must(template.ParseFiles(path + ".tmpl"))
	outFile := must(os.Create(path))
	defer outFile.Close()
	check(t.Execute(outFile, params))
}

func LayerToBaseAndSeq(seq LayerSeq) iter.Seq[ToMacro] {
	return func(yield func(ToMacro) bool) {
		for rc, key := range seq {
			tm := ToMacro{0, rc, key.Behavior()}
			if !yield(tm) {
				return
			}
		}
	}
}

func RenderLayerSeq(seq iter.Seq2[int, Layer]) iter.Seq[RenderedLayer] {
	return func(yield func(RenderedLayer) bool) {
		for n, layer := range seq {
			rl := RenderedLayer{LayerIndex(n).String(), layer.Render()}
			if !yield(rl) {
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
		Layers:    slices.Collect(RenderLayerSeq(slices.All(layers))),
		ToBaseAnd: slices.Collect(LayerToBaseAndSeq(SortedMap(layers[BASE]))),
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
