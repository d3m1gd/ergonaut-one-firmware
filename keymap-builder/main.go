package main

import (
	"cmp"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

//go:generate stringer -type=LayerIndex
const (
	BASE LayerIndex = iota
	MOVER
	NUMER
	QUICK  // 3
	REPEAT // 4
	SYS    // 5
	PARENS // 6
	CHAINS // 7
	MAXLAYERINDEX
)

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

type Params struct {
	Macros    []Macro
	Combos    []Combo
	Behaviors []Behavior
	Layers    []RenderedLayer
}

var layers = make([]Layer, MAXLAYERINDEX)
var macros = make([]Macro, 0, 64)
var behaviors = make([]Behavior, 0, 64)
var combos []Combo

func init() {
	layers[BASE] = InitWith(Trans)
	layers[BASE][l(1, 1)] = Kp(TAB) // row 1
	layers[BASE][l(1, 2)] = Kp(Q)
	layers[BASE][l(1, 3)] = Kp(W)
	layers[BASE][l(1, 4)] = Kp(E)
	layers[BASE][l(1, 5)] = Kp(R)
	layers[BASE][l(1, 6)] = KpKp(RG(T), T)
	layers[BASE][l(2, 1)] = Mt(LSHIFT, BACKSPACE) // row 2
	layers[BASE][l(2, 2)] = Kp(A)
	layers[BASE][l(2, 3)] = Mt(LSHIFT, S)
	layers[BASE][l(2, 4)] = Mt(LGUI, D)
	layers[BASE][l(2, 5)] = Mt(LALT, F)
	layers[BASE][l(2, 6)] = Kp(G)
	layers[BASE][l(3, 1)] = ModX(LCTRL, ModMorph(Kp(MINUS), Kp(PLUS), []KeyMod{MOD_LSFT, MOD_RSFT})) // row 3
	layers[BASE][l(3, 2)] = Kp(Z)
	layers[BASE][l(3, 3)] = Kp(X)
	layers[BASE][l(3, 4)] = Ref2("kpConfig", ZERO, C)
	layers[BASE][l(3, 5)] = Kp(V)
	layers[BASE][l(3, 6)] = Kp(B)
	layers[BASE][l(4, 1)] = MoTo(QUICK, CHAINS) // row 4
	layers[BASE][l(4, 2)] = MoX(NUMER, ModMorph(To(MOVER), Kp(UNDER), []KeyMod{MOD_LSFT, MOD_RSFT}))
	layers[BASE][l(4, 3)] = Mt(LCTRL, ESCAPE)

	layers[BASE][r(1, 1)] = Kp(Y) // row 1
	layers[BASE][r(1, 2)] = Kp(U)
	layers[BASE][r(1, 3)] = Kp(I)
	layers[BASE][r(1, 4)] = Kp(O)
	layers[BASE][r(1, 5)] = Kp(P)
	layers[BASE][r(1, 6)] = Kp(LBKT)
	layers[BASE][r(2, 1)] = Kp(H) // row 2
	layers[BASE][r(2, 2)] = Rmt(LALT, J)
	layers[BASE][r(2, 3)] = Rmt(LGUI, K)
	layers[BASE][r(2, 4)] = Rmt(LSHIFT, L)
	layers[BASE][r(2, 5)] = KpKp(RG(SEMI), SEMI)
	layers[BASE][r(2, 6)] = KpKp(RG(SQT), SQT)
	layers[BASE][r(3, 1)] = Kp(N) // row 3
	layers[BASE][r(3, 2)] = KpKp(RG(M), M)
	layers[BASE][r(3, 3)] = KpKp(RG(COMMA), COMMA)
	layers[BASE][r(3, 4)] = KpKp(RG(DOT), DOT)
	layers[BASE][r(3, 5)] = Kp(SLASH)
	layers[BASE][r(3, 6)] = Kp(BACKSLASH)
	layers[BASE][r(4, 1)] = Mt(LCTRL, RETURN) // row 4
	layers[BASE][r(4, 2)] = Lt(NUMER, SPACE)
	layers[BASE][r(4, 3)] = Wrap(MoTo(QUICK, CHAINS))

	layers[MOVER] = InitToLevelTrans(BASE)
	layers[MOVER][l(4, 3)] = To(BASE) // row 4
	layers[MOVER][r(2, 1)] = Kp(LEFT) // row 2
	layers[MOVER][r(2, 2)] = Rmt(LALT, DOWN)
	layers[MOVER][r(2, 3)] = Rmt(LGUI, UP)
	layers[MOVER][r(2, 4)] = Rmt(LSHIFT, RIGHT)

	layers[NUMER] = InitWith(Trans)
	layers[NUMER][l(1, 1)] = Kp(LS(TAB))
	layers[NUMER][l(1, 6)] = Kp(TILDE)
	layers[NUMER][l(2, 1)] = Kp(DELETE) // row 2
	// layers[NUMER][l(2, 3)] = ModX(LSHIFT, Brackets())
	// layers[NUMER][l(2, 4)] = ModX(LGUI, Parens())
	// layers[NUMER][l(2, 5)] = ModX(LALT, Curlies())
	layers[NUMER][l(3, 1)] = Mt(LCTRL, PLUS) // row 3
	layers[NUMER][l(3, 5)] = Kp(LS(INSERT))
	layers[NUMER][l(4, 2)] = Kp(UNDERSCORE)

	layers[NUMER][r(1, 1)] = Kp(N0) // row 1
	layers[NUMER][r(1, 2)] = Kp(N1)
	layers[NUMER][r(1, 3)] = Kp(N2)
	layers[NUMER][r(1, 4)] = Kp(N3)
	layers[NUMER][r(1, 6)] = Kp(RBKT)
	layers[NUMER][r(2, 1)] = Ref0("mmEquals") // row 2
	layers[NUMER][r(2, 2)] = Mt(LALT, N4)
	layers[NUMER][r(2, 3)] = Mt(LGUI, N5)
	layers[NUMER][r(2, 4)] = Mt(LSHIFT, N6)
	layers[NUMER][r(2, 5)] = Kp(COLON)
	layers[NUMER][r(2, 6)] = Ref0("mmQuoteGrave")
	layers[NUMER][r(3, 1)] = Kp(PLUS) // row 3
	layers[NUMER][r(3, 2)] = Kp(N7)
	layers[NUMER][r(3, 3)] = KpKp(RG(COMMA), N8)
	layers[NUMER][r(3, 4)] = KpKp(RG(DOT), N9)
	layers[NUMER][r(3, 5)] = Kp(LS(SLASH))
	layers[NUMER][r(3, 6)] = Kp(PIPE)

	layers[QUICK] = InitWith(Trans)
	layers[QUICK][l(1, 5)] = Kp(LG(C_VOL_UP)) // row 1
	layers[QUICK][l(1, 6)] = Kp(C_VOL_UP)
	layers[QUICK][l(2, 5)] = Kp(LG(C_VOL_DN)) // row 2
	layers[QUICK][l(2, 6)] = Kp(C_VOL_DN)
	layers[QUICK][r(1, 5)] = Kp(PSCRN) // row 1
	layers[QUICK][r(1, 6)] = Kp(LC(RBKT))
	layers[QUICK][r(2, 1)] = Kp(HOME) // row 2
	layers[QUICK][r(2, 2)] = Rmt(LALT, PG_DN)
	layers[QUICK][r(2, 3)] = Rmt(LGUI, PG_UP)
	layers[QUICK][r(2, 4)] = Rmt(LSHIFT, END)
	layers[QUICK][r(4, 1)] = Rmt(LCTRL, F10) // row 4
	layers[QUICK][r(4, 2)] = Kp(F11)
	layers[QUICK][r(4, 3)] = Kp(F12)

	layers[REPEAT] = InitWith(Trans)

	layers[SYS] = InitWith(None)
	layers[SYS][l(1, 1)] = Ref0("bootloader")     // tab
	layers[SYS][l(1, 5)] = Ref0("sys_reset")      // r
	layers[SYS][r(1, 2)] = Ref1("out", "OUT_USB") // u
	layers[SYS][l(3, 5)] = Ref1("out", "OUT_USB") // v - left only backup
	layers[SYS][l(3, 6)] = Ref1("out", "OUT_BLE") // b
	layers[SYS][l(2, 5)] = Ref2("bt", "BT_SEL", "0")
	layers[SYS][l(2, 4)] = Ref2("bt", "BT_SEL", "1")
	layers[SYS][l(2, 3)] = Ref2("bt", "BT_SEL", "2")
	layers[SYS][l(2, 2)] = Ref2("bt", "BT_SEL", "3")
	layers[SYS][l(2, 1)] = Ref2("bt", "BT_SEL", "4")
	layers[SYS][l(3, 4)] = Ref1("bt", "BT_CLR")     // c
	layers[SYS][l(3, 3)] = Ref1("bt", "BT_CLR_ALL") // x
	layers[SYS][r(3, 1)] = Ref1("bt", "BT_CLR_ALL") // n - nuke

	layers[PARENS] = InitToLevelTrans(BASE)
	layers[PARENS][l(4, 3)] = To(BASE)
	layers[PARENS][l(2, 1)] = BackspaceDelete()
	layers[PARENS][r(2, 4)] = XThenLayer(Kp(RIGHT), BASE)

	layers[CHAINS] = InitWith(To(BASE))
	// layers[CHAINS][l(1, 4)] = Custom1("slxl", 18) // todo
	// layers[CHAINS][l(1, 5)] = Custom1("slxl", 10) // todo
	// layers[CHAINS][l(2, 3)] = To(8)               // todo
	// layers[CHAINS][l(2, 3)] = To(8)               // todo
	// layers[CHAINS][l(2, 3)] = To(8)               // todo
	// layers[CHAINS][l(2, 3)] = To(8)               // todo
	// layers[CHAINS][l(2, 3)] = To(8)               // todo
	// layers[CHAINS][l(2, 3)] = To(8)               // todo
	// layers[CHAINS][l(2, 3)] = To(8)               // todo
	// layers[CHAINS][l(2, 3)] = To(8)               // todo
	// &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &slxl 18      &slxl 10      &kp K_CANCEL
	// &kp K_CANCEL  &kp K_CANCEL  &slxl 8       &slxl 11      &kp K_CANCEL  &slxl 15
	// &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &slxl 9
	// //                                           &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL

	// &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL
	// &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL
	// &kp K_CANCEL  &slxl 17      &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL
	// &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL

	combos = []Combo{
		{
			Name:   "System",
			Refs:   []Ref{Ref1("sll", SYS)}, // todo
			Keys:   []RC{l(1, 5), l(1, 6)},
			IdleMs: 500,
		},
		{
			Name: "LeftEnter",
			Refs: []Ref{Kp(RETURN)},
			Keys: []RC{l(4, 2), l(4, 3)},
		},
		{
			Name: "LeftSpace",
			Refs: []Ref{Kp(SPACE)},
			Keys: []RC{l(4, 1), l(4, 2)},
		},
		{
			Name: "RightCaps",
			Refs: []Ref{CapsWord},
			Keys: []RC{r(2, 3), r(2, 4)},
		},
		{
			Name:     "MiddleMouse",
			Refs:     []Ref{MKp(MCLK)},
			Keys:     []RC{l(3, 4), l(3, 5)},
			IdleMs:   200,
			TimoutMs: 100,
		},
		{
			Name:     "Curlies",
			Refs:     []Ref{Curlies()},
			Keys:     []RC{l(2, 5), r(2, 2)},
			IdleMs:   100,
			TimoutMs: 80,
		},
		{
			Name:     "Parens",
			Refs:     []Ref{Parens()},
			Keys:     []RC{l(2, 4), r(2, 3)},
			IdleMs:   110,
			TimoutMs: 90,
		},
		{
			Name:     "Brackets",
			Refs:     []Ref{Brackets()},
			Keys:     []RC{l(2, 3), r(2, 4)},
			IdleMs:   120,
			TimoutMs: 100,
		},
	}
}

func renderKeymap(path string, params Params) {
	tmplPath := path + ".tmpl"
	var funcs = template.FuncMap{"join": strings.Join}
	t := must(template.New(filepath.Base(path + ".tmpl")).Funcs(funcs).ParseFiles(tmplPath))
	outFile := must(os.Create(path))
	defer outFile.Close()
	check(t.Execute(outFile, params))
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

func main() {
	slices.SortFunc(behaviors, func(a, b Behavior) int { return cmp.Compare(a.Name, b.Name) })
	slices.SortFunc(macros, func(a, b Macro) int { return cmp.Compare(a.Name, b.Name) })
	renderKeymap("../config/ergonaut_one.keymap", Params{
		Behaviors: behaviors,
		Macros:    macros,
		Combos:    combos,
		Layers:    slices.Collect(RenderLayerSeq(slices.All(layers[:MAXLAYERINDEX]))),
	})

	fmt.Println("good")
}
