package main

import (
	"fmt"
	"iter"
	"os"
	"slices"
	"text/template"
)

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
	Behaviors []Behavior
	Layers    []RenderedLayer
	Indices   []RenderedLayer // tmp hack
}

var layers = make([]Layer, MAXLAYERINDEX)
var macros = make([]Macro, 0, 64)
var behaviors = make([]Behavior, 0, 64)

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
	layers[BASE][l(4, 1)] = MoTo(QUICK, CHAINS) // row 4
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
	layers[BASE][r(4, 3)] = MoTo(QUICK, CHAINS)

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
	layers[NUM][l(2, 3)] = ModRef("LSHIFT", BracketsMacro())
	layers[NUM][l(2, 4)] = ModRef("LGUI", ParensMacro())
	layers[NUM][l(2, 5)] = ModRef("LALT", CurliesMacro())
	layers[NUM][l(3, 5)] = Kp{LS(INSERT)}
	layers[NUM][l(4, 2)] = Kp{UNDERSCORE}

	layers[NUM][r(1, 1)] = Kp{N0} // row 1
	layers[NUM][r(1, 2)] = Kp{N1}
	layers[NUM][r(1, 3)] = Kp{N2}
	layers[NUM][r(1, 4)] = Kp{N3}
	layers[NUM][r(1, 6)] = Kp{RBKT}
	layers[NUM][r(2, 1)] = Custom0("mmEquals") // row 2
	// layers[NUM][r(2, 1)] = Kp{EQUAL}
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

	layers[CHAINS] = InitWith(To{BASE})
	// layers[CHAINS][l(1, 4)] = Custom1("slxl", 18) // todo
	// layers[CHAINS][l(1, 5)] = Custom1("slxl", 10) // todo
	// layers[CHAINS][l(2, 3)] = To{8}               // todo
	// layers[CHAINS][l(2, 3)] = To{8}               // todo
	// layers[CHAINS][l(2, 3)] = To{8}               // todo
	// layers[CHAINS][l(2, 3)] = To{8}               // todo
	// layers[CHAINS][l(2, 3)] = To{8}               // todo
	// layers[CHAINS][l(2, 3)] = To{8}               // todo
	// layers[CHAINS][l(2, 3)] = To{8}               // todo
	// layers[CHAINS][l(2, 3)] = To{8}               // todo
	// &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &slxl 18      &slxl 10      &kp K_CANCEL
	// &kp K_CANCEL  &kp K_CANCEL  &slxl 8       &slxl 11      &kp K_CANCEL  &slxl 15
	// &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &slxl 9
	// //                                           &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL

	// &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL
	// &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL
	// &kp K_CANCEL  &slxl 17      &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL
	// &kp K_CANCEL  &kp K_CANCEL  &kp K_CANCEL

}

func renderKeymap(path string, params Params) {
	t := must(template.ParseFiles(path + ".tmpl"))
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
	renderKeymap("../config/ergonaut_one.keymap", Params{
		Behaviors: behaviors,
		Macros:    macros,
		Layers:    slices.Collect(RenderLayerSeq(slices.All(layers[:MAXLAYERINDEX]))),
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
