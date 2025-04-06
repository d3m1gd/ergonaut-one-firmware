package main

import (
	"cmp"
	"errors"
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
	Args() []string
	Name() string
	Behavior() string
}

func CompileBehavior(b Behavior) string {
	return b.Behavior()
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
	Macros    []Macro
	ToBaseAnd []ToMacro
	Layers    []RenderedLayer
	Indices   []RenderedLayer // tmp hack
}

type Custom struct {
	Label  string
	Fields []any
}

func (x Custom) Behavior() string {
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

func InitBy(f func(RC) Behavior) Layer {
	layer := Layer{}
	for rc := range RCs() {
		layer[rc] = f(rc)
	}

	return layer
}

func InitToLevelAndTrans(index LayerIndex) Layer {
	base := layers[BASE]
	return InitBy(func(rc RC) Behavior {
		name := fmt.Sprintf("to%d%s", index, rc)
		macro := Macro{
			Name:      name,
			Label:     fmt.Sprintf("To %d, %s", index, rc.Pretty()),
			Cells:     0,
			Behaviors: []Behavior{To{index}, base[rc]},
		}
		AddMacro(macro)
		return Custom0(name)
	})
}

func AddMacro(macro Macro) {
	i := slices.IndexFunc(macros, func(other Macro) bool {
		return macro.Name == other.Name
	})
	if i != -1 {
		panicif(!macro.Equal(macros[i]))
		return
	}
	macros = append(macros, macro)
}

type MacroParam struct {
	a, b int
}

func (mp MacroParam) Behavior() string {
	return fmt.Sprintf("&%s", mp.Name())
}

func (mp MacroParam) Name() string {
	return fmt.Sprintf("macro_param_%dto%d", mp.a, mp.b)
}

func (mp MacroParam) Args() []string {
	return []string{}
}

func XThenTransMacro(beh Behavior, index LayerIndex, rc RC) Behavior {
	layer := layers[index]
	name := fmt.Sprintf("xThenTrans%d", index)
	inner := Custom1(beh.Name(), "MACRO_PLACEHOLDER")
	AddMacro(Macro{
		Name:      name,
		Label:     fmt.Sprintf("X Then Trans %s", index),
		Cells:     1,
		Behaviors: []Behavior{MacroParam{1, 1}, inner, To{index}, layer[rc]},
	})

	anyArgs := Map(beh.Args(), func(a string) any { return a })

	return Custom{name, anyArgs}
}

func XThenLayerMacro(beh Behavior, index LayerIndex) Behavior {
	name := fmt.Sprintf("xThenLayer%d", index)
	args := beh.Args()
	inner := Custom1(beh.Name(), "MACRO_PLACEHOLDER")
	AddMacro(Macro{
		Name:      name,
		Label:     fmt.Sprintf("X Then Layer %s", index),
		Cells:     1,
		Behaviors: []Behavior{MacroParam{1, 1}, inner, To{index}},
	})

	anyArgs := Map(args, func(a string) any { return a })

	return Custom{name, anyArgs}
}

func BackspaceDeleteMacro() Behavior {
	name := "bspcdel"
	AddMacro(Macro{
		Name:      name,
		Label:     "Backspace Delete",
		Cells:     0,
		Behaviors: []Behavior{Kp{BSPC}, Kp{DEL}},
	})

	return Custom0(name)
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

// todo     bindings
//
//	= <&macro_param_1to1>
//	, <&macro_press &mo MACRO_PLACEHOLDER>
//	, <&macro_param_2to1>
//	, <&macro_press &kp MACRO_PLACEHOLDER>
//	, <&macro_pause_for_release>
//	, <&macro_param_2to1>
//	, <&macro_release &kp MACRO_PLACEHOLDER>
//	, <&macro_param_1to1>
//	, <&macro_release &mo MACRO_PLACEHOLDER>
//	;

type Macro struct {
	Name      string
	Label     string
	Cells     int
	Behaviors []Behavior
}

func (m Macro) Compatible() string {
	return map[int]string{
		0: "behavior-macro",
		1: "behavior-macro-one-param",
		2: "behavior-macro-two-param",
	}[m.Cells]
}

func (m Macro) Bindings() string {
	return strings.Join(Map(m.Behaviors, CompileBehavior), " ")
}

func (m Macro) Equal(other Macro) bool {
	eq := true
	eq = eq && m.Name == other.Name
	eq = eq && m.Label == other.Label
	eq = eq && m.Cells == other.Cells
	eq = eq && slices.EqualFunc(m.Behaviors, other.Behaviors, func(a, b Behavior) bool {
		return a.Behavior() == b.Behavior()
	})
	return eq
}

var layers = make([]Layer, MAXLAYERINDEX)
var macros = make([]Macro, 0, 64)

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
	layers[BASE][l(3, 4)] = Custom2("kpConfig", "0", "C")
	layers[BASE][l(3, 5)] = Kp{V}
	layers[BASE][l(3, 6)] = Kp{B}
	layers[BASE][l(4, 1)] = Custom2("lslxl", QUICK, CHAINS) // row 4
	layers[BASE][l(4, 2)] = Custom2("lmmNumMoveUnder", NUM, "0")
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
	layers[BASE][r(4, 3)] = Custom1("slxl", CHAINS)

	layers[MOVE] = InitToLevelAndTrans(BASE)
	layers[MOVE][l(4, 3)] = To{BASE} // row 4
	layers[MOVE][r(2, 1)] = Kp{LEFT} // row 2
	layers[MOVE][r(2, 2)] = Rmt{LALT, DOWN}
	layers[MOVE][r(2, 3)] = Rmt{LGUI, UP}
	layers[MOVE][r(2, 4)] = Rmt{LSHIFT, RIGHT}

	layers[NUM] = InitWith(Trans{})
	layers[NUM][l(1, 1)] = Kp{LS(TAB)}
	layers[NUM][l(1, 6)] = Kp{TILDE}
	layers[NUM][l(2, 1)] = Kp{DELETE} // row 2
	layers[NUM][l(2, 3)] = Custom2("mtBracket", "LSHIFT", "0")
	layers[NUM][l(2, 4)] = Custom2("mtParen", "LGUI", "0")
	layers[NUM][l(2, 5)] = Custom2("mtCurly", "LALT", "0")
	layers[NUM][l(3, 5)] = Kp{LS(INSERT)}
	layers[NUM][l(4, 2)] = Kp{UNDERSCORE}

	layers[NUM][r(1, 1)] = Kp{N0} // row 1
	layers[NUM][r(1, 2)] = Kp{N1}
	layers[NUM][r(1, 3)] = Kp{N2}
	layers[NUM][r(1, 4)] = Kp{N3}
	layers[NUM][r(1, 6)] = Kp{RBKT}
	layers[NUM][r(2, 1)] = Custom0("mmEquals") // row 2
	layers[NUM][r(2, 2)] = Mt{LALT, N4}
	layers[NUM][r(2, 3)] = Mt{LGUI, N5}
	layers[NUM][r(2, 4)] = Mt{LSHIFT, N6}
	layers[NUM][r(2, 5)] = Kp{COLON}
	layers[NUM][r(2, 6)] = Custom0("mmQuoteGrave")
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
	layers[SYS][l(1, 1)] = Custom0("bootloader")     // tab
	layers[SYS][l(1, 5)] = Custom0("sys_reset")      // r
	layers[SYS][r(1, 2)] = Custom1("out", "OUT_USB") // u
	layers[SYS][l(3, 5)] = Custom1("out", "OUT_USB") // v - single half backup
	layers[SYS][l(3, 6)] = Custom1("out", "OUT_BLE") // b
	layers[SYS][l(2, 5)] = Custom2("bt", "BT_SEL", "0")
	layers[SYS][l(2, 4)] = Custom2("bt", "BT_SEL", "1")
	layers[SYS][l(2, 3)] = Custom2("bt", "BT_SEL", "2")
	layers[SYS][l(2, 2)] = Custom2("bt", "BT_SEL", "3")
	layers[SYS][l(2, 1)] = Custom2("bt", "BT_SEL", "4")
	layers[SYS][l(3, 4)] = Custom1("bt", "BT_CLR")     // c
	layers[SYS][l(3, 3)] = Custom1("bt", "BT_CLR_ALL") // x
	layers[SYS][r(3, 1)] = Custom1("bt", "BT_CLR_ALL") // n - nuke

	layers[PARENS] = InitToLevelAndTrans(BASE)
	layers[PARENS][l(4, 3)] = To{BASE}
	layers[PARENS][l(2, 1)] = BackspaceDeleteMacro()
	layers[PARENS][r(2, 4)] = XThenLayerMacro(Kp{RIGHT}, 0)
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
	renderKeymap("../config/ergonaut_one.keymap", Params{
		Macros:    macros,
		Layers:    slices.Collect(RenderLayerSeq(slices.All(layers[:PARENS+1]))),
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

func errif(cond bool) error {
	if cond {
		return errors.New("condition failed")
	}

	return nil
}

func panicif(cond bool) {
	if cond {
		panic("condition failed")
	}
}

func Map[T, U any](s []T, f func(T) U) []U {
	us := make([]U, len(s))
	for i := range s {
		us[i] = f(s[i])
	}

	return us
}

func toString(x any) string {
	return fmt.Sprintf("%s", x)
}
