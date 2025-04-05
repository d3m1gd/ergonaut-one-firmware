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
	QUICK  // 3
	REPEAT // 4
	SYS    // 5
	PARENS // 6
	CHAINS // 7
	MAXLAYERINDEX
)

func (l Layer) Render() []string {
	widths := slices.Repeat([]int{0}, 12)
	cells := make([][]string, 4)
	for i := range cells {
		cells[i] = slices.Repeat([]string{""}, 12)
	}
	for n := range 42 {
		rc := RCFrom(n + 1)
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
	Index int
	Name  string
	Rows  []string
}

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
	Level    int
	RC       RC
	Behavior string
}

type Params struct {
	ToBaseAnd []ToMacro
	Layers    []RenderedLayer
	Indices   []RenderedLayer // tmp hack
}

type Kp struct {
	Tap KeyCode
}

func (x Kp) Behavior() string {
	return fmt.Sprintf("&kp %s", x.Tap)
}

type KpKp struct {
	Hold KeyCode
	Tap  KeyCode
}

func (x KpKp) Behavior() string {
	return fmt.Sprintf("&kpkp %s %s", x.Hold, x.Tap)
}

type Rmt struct {
	Hold KeyCode
	Tap  KeyCode
}

func (x Rmt) Behavior() string {
	return fmt.Sprintf("&rmt %s %s", x.Hold, x.Tap)
}

func behaviorArgToString(a any) string {
	switch v := a.(type) {
	case LayerIndex:
		return fmt.Sprintf("%s", v)
	}
	return fmt.Sprintf("%s", a)
}

type Custom2 struct {
	Name string
	A    any
	B    any
}

func (x Custom2) Behavior() string {
	a := behaviorArgToString(x.A)
	b := behaviorArgToString(x.B)
	return fmt.Sprintf("&%s %s %s", x.Name, a, b)
}

type Custom1 struct {
	Name string
	A    any
}

func (x Custom1) Behavior() string {
	a := behaviorArgToString(x.A)
	return fmt.Sprintf("&%s %s", x.Name, a)
}

type Custom0 struct {
	Name string
}

func (x Custom0) Behavior() string {
	return fmt.Sprintf("&%s", x.Name)
}

type Lt struct {
	Layer LayerIndex
	Tap   KeyCode
}

func (x Lt) Behavior() string {
	return fmt.Sprintf("&lt %s %s", x.Layer, x.Tap)
}

type To struct {
	Layer LayerIndex
}

func (x To) Behavior() string {
	return fmt.Sprintf("&to %s", x.Layer)
}

type Mt struct {
	Hold KeyCode
	Tap  KeyCode
}

func (x Mt) Behavior() string {
	return fmt.Sprintf("&mt %s %s", x.Hold, x.Tap)
}

type Trans struct{}

func (_ Trans) Behavior() string {
	return "&trans"
}

type None struct{}

func (_ None) Behavior() string {
	return "&none"
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

func InitWith(b Behavior) Layer {
	layer := Layer{}
	for rc := range RCs() {
		layer[rc] = b
	}

	return layer
}

func InitToLevelAndTrans(level int) Layer {
	layer := Layer{}
	for rc := range RCs() {
		// replace with proper trans when zmk ready
		layer[rc] = Custom0{fmt.Sprintf("to%d%s", level, rc)}
	}

	return layer
}

func RCs() iter.Seq[RC] {
	return func(yield func(RC) bool) {
		for i := range 42 {
			if !yield(RCFrom(i + 1)) {
				return
			}
		}
	}
}

var layers = make([]Layer, MAXLAYERINDEX)

func init() {
	layers[BASE] = InitWith(Trans{})
	layers[BASE][l(1, 1)] = Kp{TAB} // row 1
	layers[BASE][l(1, 2)] = Kp{Q}
	layers[BASE][l(1, 3)] = Kp{W}
	layers[BASE][l(1, 4)] = Kp{E}
	layers[BASE][l(1, 5)] = Kp{R}
	layers[BASE][l(1, 6)] = KpKp{RG(T), T}
	layers[BASE][l(2, 1)] = Mt{LSHIFT, BACKSPACE} // row 2
	layers[BASE][l(2, 2)] = Kp{A}
	layers[BASE][l(2, 3)] = Mt{LSHIFT, S}
	layers[BASE][l(2, 4)] = Mt{LGUI, D}
	layers[BASE][l(2, 5)] = Mt{LALT, F}
	layers[BASE][l(2, 6)] = Kp{G}
	layers[BASE][l(3, 1)] = Mt{LCTRL, MINUS} // row 3
	layers[BASE][l(3, 2)] = Kp{Z}
	layers[BASE][l(3, 3)] = Kp{X}
	layers[BASE][l(3, 4)] = Custom2{"kpConfig", "0", "C"}
	layers[BASE][l(3, 5)] = Kp{V}
	layers[BASE][l(3, 6)] = Kp{B}
	layers[BASE][l(4, 1)] = Custom2{"lslxl", QUICK, CHAINS} // row 4
	layers[BASE][l(4, 2)] = Custom2{"lmmNumMoveUnder", NUM, "0"}
	layers[BASE][l(4, 3)] = Mt{LCTRL, ESCAPE}

	layers[BASE][r(1, 1)] = Kp{Y} // row 1
	layers[BASE][r(1, 2)] = Kp{U}
	layers[BASE][r(1, 3)] = Kp{I}
	layers[BASE][r(1, 4)] = Kp{O}
	layers[BASE][r(1, 5)] = Kp{P}
	layers[BASE][r(1, 6)] = Kp{LBKT}
	layers[BASE][r(2, 1)] = Kp{H} // row 2
	layers[BASE][r(2, 2)] = Rmt{LALT, J}
	layers[BASE][r(2, 3)] = Rmt{LGUI, K}
	layers[BASE][r(2, 4)] = Rmt{LSHIFT, L}
	layers[BASE][r(2, 5)] = KpKp{RG(SEMI), SEMI}
	layers[BASE][r(2, 6)] = KpKp{RG(SQT), SQT}
	layers[BASE][r(3, 1)] = Kp{N} // row 3
	layers[BASE][r(3, 2)] = KpKp{RG(M), M}
	layers[BASE][r(3, 3)] = KpKp{RG(COMMA), COMMA}
	layers[BASE][r(3, 4)] = KpKp{RG(DOT), DOT}
	layers[BASE][r(3, 5)] = Kp{SLASH}
	layers[BASE][r(3, 6)] = Kp{BACKSLASH}
	layers[BASE][r(4, 1)] = Mt{LCTRL, RETURN} // row 4
	layers[BASE][r(4, 2)] = Lt{NUM, SPACE}
	layers[BASE][r(4, 3)] = Custom1{"slxl", CHAINS}

	layers[MOVE] = InitToLevelAndTrans(0)
	layers[MOVE][l(4, 3)] = To{BASE} // row 4
	layers[MOVE][r(2, 1)] = Kp{LEFT} // row 2
	layers[MOVE][r(2, 2)] = Rmt{LALT, DOWN}
	layers[MOVE][r(2, 3)] = Rmt{LGUI, UP}
	layers[MOVE][r(2, 4)] = Rmt{LSHIFT, RIGHT}

	layers[NUM] = InitWith(Trans{})
	layers[NUM][l(1, 1)] = Kp{LS(TAB)}
	layers[NUM][l(1, 6)] = Kp{TILDE}
	layers[NUM][l(2, 1)] = Kp{DELETE} // row 2
	layers[NUM][l(2, 3)] = Custom2{"mtBracket", "LSHIFT", "0"}
	layers[NUM][l(2, 4)] = Custom2{"mtParen", "LGUI", "0"}
	layers[NUM][l(2, 5)] = Custom2{"mtCurly", "LALT", "0"}
	layers[NUM][l(3, 5)] = Kp{LS(INSERT)}
	layers[NUM][l(4, 2)] = Kp{UNDERSCORE}

	layers[NUM][r(1, 1)] = Kp{N0} // row 1
	layers[NUM][r(1, 2)] = Kp{N1}
	layers[NUM][r(1, 3)] = Kp{N2}
	layers[NUM][r(1, 4)] = Kp{N3}
	layers[NUM][r(1, 6)] = Kp{RBKT}
	layers[NUM][r(2, 1)] = Custom0{"mmEquals"} // row 2
	layers[NUM][r(2, 2)] = Mt{LALT, N4}
	layers[NUM][r(2, 3)] = Mt{LGUI, N5}
	layers[NUM][r(2, 4)] = Mt{LSHIFT, N6}
	layers[NUM][r(2, 5)] = Kp{COLON}
	layers[NUM][r(2, 6)] = Custom0{"mmQuoteGrave"}
	layers[NUM][r(3, 1)] = Kp{PLUS} // row 3
	layers[NUM][r(3, 2)] = Kp{N7}
	layers[NUM][r(3, 3)] = KpKp{RG(COMMA), N8}
	layers[NUM][r(3, 4)] = KpKp{RG(DOT), N9}
	layers[NUM][r(3, 5)] = Kp{LS(SLASH)}
	layers[NUM][r(3, 6)] = Kp{PIPE}

	layers[QUICK] = InitWith(Trans{})
	layers[QUICK][l(1, 5)] = Kp{LG(C_VOL_UP)} // row 1
	layers[QUICK][l(1, 6)] = Kp{C_VOL_UP}
	layers[QUICK][l(2, 5)] = Kp{LG(C_VOL_DN)} // row 2
	layers[QUICK][l(2, 6)] = Kp{C_VOL_DN}
	layers[QUICK][r(1, 5)] = Kp{PSCRN} // row 1
	layers[QUICK][r(1, 6)] = Kp{LC(RBKT)}
	layers[QUICK][r(2, 1)] = Kp{HOME} // row 2
	layers[QUICK][r(2, 2)] = Rmt{LALT, PG_DN}
	layers[QUICK][r(2, 3)] = Rmt{LGUI, PG_UP}
	layers[QUICK][r(2, 4)] = Rmt{LSHIFT, END}
	layers[QUICK][r(4, 1)] = Rmt{LCTRL, F10} // row 4
	layers[QUICK][r(4, 2)] = Kp{F11}
	layers[QUICK][r(4, 3)] = Kp{F12}

	layers[REPEAT] = InitWith(Trans{})
	layers[SYS] = InitWith(None{})
	layers[SYS][l(1, 1)] = Custom0{"bootloader"}     // tab
	layers[SYS][l(1, 5)] = Custom0{"sys_reset"}      // r
	layers[SYS][r(1, 2)] = Custom1{"out", "OUT_USB"} // u
	layers[SYS][l(3, 5)] = Custom1{"out", "OUT_USB"} // v - single half backup
	layers[SYS][l(3, 6)] = Custom1{"out", "OUT_BLE"} // b
	layers[SYS][l(2, 5)] = Custom2{"bt", "BT_SEL", "0"}
	layers[SYS][l(2, 4)] = Custom2{"bt", "BT_SEL", "1"}
	layers[SYS][l(2, 3)] = Custom2{"bt", "BT_SEL", "2"}
	layers[SYS][l(2, 2)] = Custom2{"bt", "BT_SEL", "3"}
	layers[SYS][l(2, 1)] = Custom2{"bt", "BT_SEL", "4"}
	layers[SYS][l(3, 4)] = Custom1{"bt", "BT_CLR"}     // c
	layers[SYS][l(3, 3)] = Custom1{"bt", "BT_CLR_ALL"} // x
	layers[SYS][r(3, 1)] = Custom1{"bt", "BT_CLR_ALL"} // n - nuke
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
			rl := RenderedLayer{n, LayerIndex(n).String(), layer.Render()}
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
	renderKeymap("config/ergonaut_one.keymap", Params{
		Layers:    slices.Collect(RenderLayerSeq(slices.All(layers[:SYS+1]))),
		ToBaseAnd: slices.Collect(LayerToBaseAndSeq(SortedMap(layers[BASE]))),
		Indices: func() []RenderedLayer {
			a := []RenderedLayer{}
			for i := range MAXLAYERINDEX {
				a = append(a, RenderedLayer{
					Index: int(i),
					Name:  i.String(),
				})
			}
			return a
		}(),
	})

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
