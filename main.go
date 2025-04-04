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

//go:generate stringer -type=LayerIndex
const (
	BASE LayerIndex = iota
	MOVE
	NUM
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

type Kp struct {
	Tap string
}

func (x Kp) Behavior() string {
	return fmt.Sprintf("&kp %s", x.Tap)
}

type KpKp struct {
	Hold string
	Tap  string
}

func (x KpKp) Behavior() string {
	return fmt.Sprintf("&kpkp %s %s", x.Hold, x.Tap)
}

type Rmt struct {
	Hold string
	Tap  string
}

func (x Rmt) Behavior() string {
	return fmt.Sprintf("&rmt %s %s", x.Hold, x.Tap)
}

type Custom struct {
	Name string
	Hold string
	Tap  string
}

func (x Custom) Behavior() string {
	return fmt.Sprintf("&%s %s %s", x.Name, x.Hold, x.Tap)
}

type Custom1 struct {
	Name string
	Hold string
}

func (x Custom1) Behavior() string {
	return fmt.Sprintf("&%s %s", x.Name, x.Hold)
}

type Custom0 struct {
	Name string
}

func (x Custom0) Behavior() string {
	return fmt.Sprintf("&%s", x.Name)
}

type Lt struct {
	Layer LayerIndex
	Tap   string
}

func (x Lt) Behavior() string {
	return fmt.Sprintf("&lt %d %s", x.Layer, x.Tap)
}

type Mt struct {
	Hold string
	Tap  string
}

func (x Mt) Behavior() string {
	return fmt.Sprintf("&mt %s %s", x.Hold, x.Tap)
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

func InitToLevelTrans(level int) Layer {
	layer := Layer{}
	for i := range 42 {
		rc := RCFrom(i + 1)
		layer[rc] = Custom0{fmt.Sprintf("to%d%s", level, rc)}
	}

	return layer
}

var layers = make([]Layer, MAXLAYERINDEX)

func init() {
	layers[BASE] = InitTrans()
	layers[BASE][l(1, 1)] = Kp{"TAB"}
	layers[BASE][l(1, 2)] = Kp{"Q"}
	layers[BASE][l(1, 3)] = Kp{"W"}
	layers[BASE][l(1, 4)] = Kp{"E"}
	layers[BASE][l(1, 5)] = Kp{"R"}
	layers[BASE][l(1, 6)] = KpKp{"RG(T)", "T"}
	layers[BASE][r(1, 1)] = Kp{"Y"}
	layers[BASE][r(1, 2)] = Kp{"U"}
	layers[BASE][r(1, 3)] = Kp{"I"}
	layers[BASE][r(1, 4)] = Kp{"O"}
	layers[BASE][r(1, 5)] = Kp{"P"}
	layers[BASE][r(1, 6)] = Kp{"LBKT"}

	layers[BASE][l(2, 1)] = Mt{"LSHIFT", "BACKSPACE"}
	layers[BASE][l(2, 2)] = Kp{"A"}
	layers[BASE][l(2, 3)] = Mt{"LSHIFT", "S"}
	layers[BASE][l(2, 4)] = Mt{"LGUI", "D"}
	layers[BASE][l(2, 5)] = Mt{"LALT", "F"}
	layers[BASE][l(2, 6)] = Kp{"G"}
	layers[BASE][r(2, 1)] = Kp{"H"}
	layers[BASE][r(2, 2)] = Rmt{"LALT", "J"}
	layers[BASE][r(2, 3)] = Rmt{"LEFT_WIN", "K"} // todo
	layers[BASE][r(2, 4)] = Rmt{"LSHIFT", "L"}
	layers[BASE][r(2, 5)] = KpKp{"RG(SEMI)", "SEMI"}
	layers[BASE][r(2, 6)] = KpKp{"RG(SINGLE_QUOTE)", "SINGLE_QUOTE"}

	layers[BASE][l(3, 1)] = Mt{"LCTRL", "MINUS"}
	layers[BASE][l(3, 2)] = Kp{"Z"}
	layers[BASE][l(3, 3)] = Kp{"X"}
	layers[BASE][l(3, 4)] = Custom{"kpConfig", "0", "C"}
	layers[BASE][l(3, 5)] = Kp{"V"}
	layers[BASE][l(3, 6)] = Kp{"B"}
	layers[BASE][r(3, 1)] = Kp{"N"}
	layers[BASE][r(3, 2)] = KpKp{"RG(M)", "M"}
	layers[BASE][r(3, 3)] = KpKp{"RG(COMMA)", "COMMA"}
	layers[BASE][r(3, 4)] = KpKp{"RG(DOT)", "DOT"}
	layers[BASE][r(3, 5)] = Kp{"SLASH"}
	layers[BASE][r(3, 6)] = Kp{"BACKSLASH"}

	layers[BASE][l(4, 1)] = Custom{"lslxl", "3", "7"}           // todo layer
	layers[BASE][l(4, 2)] = Custom{"lmmNumMoveUnder", "2", "0"} // todo layer
	layers[BASE][l(4, 3)] = Mt{"LCTRL", "ESCAPE"}
	layers[BASE][r(4, 1)] = Mt{"LCTRL", "RETURN"}
	layers[BASE][r(4, 2)] = Lt{NUM, "SPACE"}     // todo layer
	layers[BASE][r(4, 3)] = Custom1{"slxl", "7"} // todo layer

	layers[MOVE] = InitToLevelTrans(0)

	layers[MOVE][r(2, 1)] = Kp{"LEFT"}
	layers[MOVE][r(2, 2)] = Rmt{"LALT", "DOWN"}
	layers[MOVE][r(2, 3)] = Rmt{"LGUI", "UP"}
	layers[MOVE][r(2, 4)] = Rmt{"LSHIFT", "RIGHT"}
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
		Layers:    slices.Collect(RenderLayerSeq(slices.All([]Layer{layers[BASE], layers[MOVE]}))),
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
