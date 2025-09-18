package main

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	. "keyboard/behavior"
	"keyboard/chain"
	"keyboard/combo"
	. "keyboard/key"
	"keyboard/layer"
	. "keyboard/layout"
	"keyboard/macro"
	"keyboard/ref"
	"keyboard/rowcol"
	"keyboard/util"
)

const (
	comboBothSidesTimeout = 90
	comboBothSidesIdle    = 300
	shortSticky           = 250
	longSticky            = 500
	symbolSticky          = 750
)

const (
	a = LALT
	w = LGUI
	s = LSHIFT
	c = LCTRL
)

var (
	BASE   = layer.New("BASE", layer.InitWith(Trans))
	MOVER  = layer.New("MOVER", layer.InitWith(Trans))
	NUMER  = layer.New("NUMER", layer.InitWith(Trans))
	QUICK  = layer.New("QUICK", layer.InitWith(Trans))
	SYS    = layer.New("SYS", layer.InitWith(To(BASE)))
	SYMBOL = layer.New("SYMBOL", layer.InitWith(Trans))
	CHAINS = layer.New(chain.Name(""), layer.InitWith(To(BASE)))
)

func init() {
	defer util.Reporter(nil)

	BASE.Extend(layer.Cells{
		L11: Kp(TAB), // row 1
		L12: Kp(Q),
		L13: Kp(W),
		L14: Kp(E),
		L15: Kp(R),
		L16: KpKp(RG(T), T),
		L21: Mt(LSHIFT, BACKSPACE), // row 2
		L22: Kp(A),
		L23: Mt(LSHIFT, S),
		L24: Mt(LGUI, D),
		L25: Mt(LALT, F),
		L26: Kp(G),
		L31: ModX(LCTRL, ModMorph(Kp(MINUS), Kp(PLUS), ShiftsAndCtrls, Ctrls)), // row 3
		L32: Kp(Z),
		L33: Kp(X),
		L34: XKp(Text("XdgConfig", `"$HOME/.config"/`), C),
		L35: Kp(V),
		L36: Kp(B),
		L41: MoTo(MOVER, CHAINS), // row 4
		L42: MoX(NUMER, ModMorph(Sl(SYMBOL, longSticky), Kp(UNDER), Shifts, nil)),
		L43: Mt(LCTRL, ESCAPE),
		// ------------------------
		R11: Kp(Y), // row 1
		R12: Kp(U),
		R13: Kp(I),
		R14: Kp(O),
		R15: Kp(P),
		R16: Kp(LBKT),
		R21: Kp(H), // row 2
		R22: Mt(LALT, J),
		R23: Mt(LGUI, K),
		R24: Rmt(LSHIFT, L),
		R25: KpKp(RG(SEMI), SEMI),
		R26: KpKp(RG(SQT), SQT),
		R31: Kp(N), // row 3
		R32: KpKp(RG(M), M),
		R33: KpKp(RG(COMMA), COMMA),
		R34: KpKp(RG(DOT), DOT),
		R35: Kp(SLASH),
		R36: Kp(BACKSLASH),
		R41: Mt(LCTRL, RETURN), // row 4
		R42: Lt(NUMER, SPACE),
		R43: MoTo(MOVER, CHAINS),
	})

	MOVER.Extend(layer.Cells{
		// L42: Sl(MOVER, symbolSticky),
		// L43: Off(MOVER),
		// ------------------------
		R21: Kp(LEFT),
		R22: Kp(DOWN),
		R23: Kp(UP),
		R24: Kp(RIGHT),
		R11: Kp(HOME),
		R12: Kp(PG_DN),
		R13: Kp(PG_UP),
		R14: Kp(END),
		R15: Kp(PSCRN),
		// R21: KpSl(LEFT, MOVER, moverDuration),
		// R22: HoldTap(Kp(LALT), KpSl(DOWN, MOVER, moverDuration)),
		// R23: HoldTap(Kp(LGUI), KpSl(UP, MOVER, moverDuration)),
		// R24: HoldTap(Kp(LSHIFT), KpSl(RIGHT, MOVER, moverDuration)),
		// R23: Mt(LGUI, UP),
		// R24: Mt(LSHIFT, RIGHT),
		// R11: KpSl(RPAR, MOVER, symbolSticky),
		// R12: KpSl(EXCL, MOVER, symbolSticky),
		// R13: KpSl(AT, MOVER, shortSticky),
		// R14: KpSl(HASH, MOVER, shortSticky),
		// R21: KpSl(LEFT, MOVER, symbolSticky),
		// R22: KpSl(DOWN, MOVER, symbolSticky),
		// R23: KpSl(UP, MOVER, symbolSticky),
		// R24: KpSl(RIGHT, MOVER, symbolSticky),
		// R26: KpSl(GRAVE, MOVER, shortSticky),
		// R32: KpSl(AMPS, MOVER, shortSticky),
		// R33: KpSl(STAR, MOVER, shortSticky),
		// R34: KpSl(LPAR, MOVER, shortSticky),
		// R36: KpSl(PIPE, MOVER, shortSticky),
		// R23: Mt(LGUI, UP),
	})

	NUMER.Extend(layer.Cells{
		L11: Kp(LS(TAB)),
		L21: Kp(DELETE),      // row 2
		L31: Mt(LCTRL, PLUS), // row 3
		L35: Kp(LS(INSERT)),
		L42: Kp(UNDER),
		// ------------------------
		R11: Kp(N0), // row 1
		R12: Kp(N1),
		R13: Kp(N2),
		R14: Kp(N3),
		R16: Kp(RBKT),
		R21: ModMorph(Kp(EQUAL), Kp(EQUAL), Shifts, nil), // row 2
		R22: Mt(LALT, N4),
		R23: Mt(LGUI, N5),
		R24: Mt(LSHIFT, N6),
		R25: Kp(COLON),
		R26: ModMorph(Kp(DQT), Kp(GRAVE), Shifts, nil),
		R31: Kp(PLUS), // row 3
		R32: Kp(N7),
		R33: KpKp(RG(COMMA), N8),
		R34: KpKp(RG(DOT), N9),
		R35: Kp(LS(SLASH)),
		R36: Kp(PIPE),
	})

	QUICK.Extend(layer.Cells{
		// L15: Kp(LG(C_VOL_UP)), // row 1
		// L16: Kp(C_VOL_UP),
		// L25: Kp(LG(C_VOL_DN)), // row 2
		// L26: Kp(C_VOL_DN),
		// ------------------------
		R15: Kp(PSCRN), // row 1
		R16: Kp(LC(RBKT)),
		R21: Kp(HOME), // row 2
		// R22: Rmt(LALT, PG_DN),
		// R23: Rmt(LGUI, PG_UP),
		// R24: Rmt(LSHIFT, END),
		// R41: Rmt(LCTRL, F10), // row 4
		R22: Mt(LALT, PG_DN),
		R23: Mt(LGUI, PG_UP),
		R24: Mt(LSHIFT, END),
		R41: Mt(LCTRL, F10), // row 4
		R42: Kp(F11),
		R43: Kp(F12),
	})

	SYS.Extend(layer.Cells{
		L11: ref.Ref0("bootloader"),     // tab
		L15: ref.Ref0("sys_reset"),      // r
		L35: ref.Ref1("out", "OUT_USB"), // v - left only backup
		L36: ref.Ref1("out", "OUT_BLE"), // b
		L25: ref.Ref2("bt", "BT_SEL", "0"),
		L24: ref.Ref2("bt", "BT_SEL", "1"),
		L23: ref.Ref2("bt", "BT_SEL", "2"),
		L22: ref.Ref2("bt", "BT_SEL", "3"),
		L21: ref.Ref2("bt", "BT_SEL", "4"),
		L34: ref.Ref1("bt", "BT_CLR"),     // c
		L33: ref.Ref1("bt", "BT_CLR_ALL"), // x
		// ------------------------
		R12: ref.Ref1("out", "OUT_USB"),   // u
		R31: ref.Ref1("bt", "BT_CLR_ALL"), // n - nuke
	})

	SYMBOL.Extend(layer.Cells{
		// L23: Skl(LSHIFT),
		// L24: Skl(LGUI),
		// L25: Skl(LALT),
		// L43: Skl(LCTRL),
		L11: Kp(TILDE),
		// ------------------------
		R11: KpSl(RPAR, SYMBOL, shortSticky),
		R12: KpSl(EXCL, SYMBOL, symbolSticky),
		R13: Kp(AT),
		R14: KpSl(HASH, SYMBOL, shortSticky),
		R16: Kp(RBKT),
		R21: KpSl(EQUAL, SYMBOL, symbolSticky),
		R22: Kp(DLLR),
		R23: Kp(PRCNT),
		R24: Kp(CARET),
		R25: Kp(COLON),
		R26: KpSl(GRAVE, SYMBOL, shortSticky),
		R31: KpSl(PLUS, SYMBOL, shortSticky),
		R32: KpSl(AMPS, SYMBOL, shortSticky),
		R33: KpSl(STAR, SYMBOL, shortSticky),
		R34: Kp(LPAR),
		R36: KpSl(PIPE, SYMBOL, shortSticky),
	})

	chain.Add(CHAINS, layer.InitWith(To(BASE)), chain.KeyRefs{
		"sdf": Kp(X),
		"gie": Text("GoIfError", CursorAt("go if%fine", "%")),
		"fu":  Kp(F1),
		"fi":  Kp(F2), // todo i -> 2 alias
		"fo":  Kp(F3),
		"fj":  Kp(F4),
		"fk":  Kp(F5),
		"fl":  Kp(F6),
		"fm":  Kp(F7),
		"f,":  Kp(F8),
		"f.":  Kp(F9),
		"fR":  Kp(F10),
		"fS":  Kp(F11),
		"f>":  Kp(F12),

		"c1": Kp(RG(LS(F1))),
		"c2": Kp(RG(LS(F2))),
		"c3": Kp(RG(LS(F3))),
		"c4": Kp(RA(RG(LS(F4)))),
		"c5": Kp(RA(RG(LS(F5)))),
		"c6": Kp(RA(RG(LS(F6)))),
	})

	combo.Add(combo.T{
		Name:   "System",
		Ref:    To(SYS),
		RCs:    []rowcol.T{L15, L16},
		Idle:   500,
		Timout: 100,
	})

	combo.Add(combo.T{
		Name: "LeftEnter",
		Ref:  Kp(RETURN),
		RCs:  []rowcol.T{L42, L43},
	})

	combo.Add(combo.T{
		Name: "LeftSpace",
		Ref:  Kp(SPACE),
		RCs:  []rowcol.T{L41, L42},
	})

	combo.Add(combo.T{
		Name: "BottomLeftCtrlShift",
		Ref:  Kp(LC(LSHIFT)),
		RCs:  []rowcol.T{L21, L31},
	})

	LeftStickyCombo(a, w, s, c)

	LeftStickyCombo(a, w, s)
	LeftStickyCombo(a, w)
	LeftStickyCombo(a, s)
	LeftStickyCombo(w, s)

	LeftStickyCombo(a, w, c)
	LeftStickyCombo(a, c)
	LeftStickyCombo(w, c)

	LeftStickyCombo(a, s, c)
	LeftStickyCombo(s, c)

	LeftStickyCombo(w, s, c)

	combo.Add(combo.T{
		Name: "RightAltWinShift",
		Ref:  Kp(LA(LW(LSHIFT))),
		RCs:  []rowcol.T{R22, R23, R24},
	})

	combo.Add(combo.T{
		Name: "RightWinShift_Caps",
		Ref:  HoldTap(Kp(LG(LSHIFT)), CapsWord),
		RCs:  []rowcol.T{R23, R24},
	})

	combo.Add(combo.T{
		Name:   "MiddleMouse",
		Ref:    MKp(MCLK),
		RCs:    []rowcol.T{L34, L35},
		Idle:   200,
		Timout: 100,
	})

	combo.Add(combo.T{
		Name:   "Parens",
		Ref:    Parens(BASE),
		RCs:    []rowcol.T{L25, R22},
		Timout: comboBothSidesTimeout + 30,
		Idle:   comboBothSidesIdle,
	})

	combo.Add(combo.T{
		Name:   "Curlies",
		Ref:    Curlies(BASE),
		RCs:    []rowcol.T{L24, R23},
		Timout: comboBothSidesTimeout + 40,
		Idle:   comboBothSidesIdle,
	})

	combo.Add(combo.T{
		Name:   "Brackets",
		Ref:    Brackets(BASE),
		RCs:    []rowcol.T{L23, R24},
		Timout: comboBothSidesTimeout - 30,
		Idle:   comboBothSidesIdle,
	})

	combo.Add(combo.T{
		Name:   "DoubleQuotes",
		Ref:    DoubleQuotes(BASE),
		RCs:    []rowcol.T{L15, R12},
		Timout: comboBothSidesTimeout + 30,
		Idle:   comboBothSidesIdle,
	})

	combo.Add(combo.T{
		Name:   "SingleQuotes",
		Ref:    SingleQuotes(BASE),
		RCs:    []rowcol.T{L14, R13},
		Timout: comboBothSidesTimeout,
		Idle:   comboBothSidesIdle,
	})

	combo.Add(combo.T{
		Name:   "BackQuotes",
		Ref:    HoldTap(Text("CodeBlock", CursorAt("```%```", "%")), BackQuotes(BASE)),
		RCs:    []rowcol.T{L13, R14},
		Timout: comboBothSidesTimeout,
		Idle:   comboBothSidesIdle,
	})

	// combo.Add(combo.T{
	// 	Name:     "CodeQuotes",
	// 	Refs:     []ref.T{},
	// 	Keys:     []rowcol.T{L16, R11},
	// 	TimoutMs: comboBothSidesTimeout,
	// 	IdleMs:   comboBothSidesIdle,
	// })
}

func renderKeymap(path string, params Params) {
	tmplPath := path + ".tmpl"
	var funcs = template.FuncMap{"join": strings.Join}
	t := util.Must(template.New(filepath.Base(path + ".tmpl")).Funcs(funcs).ParseFiles(tmplPath))
	outFile := util.Must(os.Create(path))
	defer outFile.Close()
	util.Check(t.Execute(outFile, params))
}

type Params struct {
	Layers    []layer.T
	Macros    []macro.T
	Combos    []combo.T
	Behaviors []Behavior
}

func main() {
	renderKeymap("../config/ergonaut_one.keymap", Params{
		Behaviors: Render(), // todo
		Macros:    macro.Render(),
		Combos:    combo.Render(),
		Layers:    layer.All(),
	})
}

func LeftStickyCombo(keys ...Key) {
	rowcols := util.Map(keys, func(k Key) rowcol.T {
		rc, ok := LeftModPositions[k]
		util.Panicif(!ok)
		return rc
	})

	names := util.Map(keys, func(k Key) string {
		return string(k)
	})

	combo.Add(combo.T{
		Name: "LeftStickyCombo_" + strings.Join(names, "_"),
		Ref:  HoldTapStickyLong(ModMap(keys)),
		RCs:  rowcols,
	})
}

var LeftModPositions = map[Key]rowcol.T{
	a: L25,
	w: L24,
	s: L23,
	c: L43,
}

var RightModPositions = map[Key]rowcol.T{
	a: R21,
	w: R22,
	s: R23,
	c: R41,
}
