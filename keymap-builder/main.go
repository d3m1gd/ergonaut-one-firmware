package main

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"keyboard/behavior"
	"keyboard/chain"
	"keyboard/combo"
	. "keyboard/instance"
	"keyboard/key"
	. "keyboard/key/keys"
	"keyboard/layer"
	"keyboard/macro"
	. "keyboard/util"
)

var BASE = layer.New("BASE", InitWith(Trans))
var MOVER = layer.New("MOVER", InitWith(Trans))
var NUMER = layer.New("NUMER", InitWith(Trans))
var QUICK = layer.New("QUICK", InitWith(Trans))
var REPEAT = layer.New("REPEAT", InitWith(Trans))
var SYS = layer.New("SYS", InitWith(None))
var PARENS = layer.New("PARENS", InitWith(Trans))
var CHAINS = layer.New("CHAINS", InitWith(To(BASE)))

type Params struct {
	Layers    []layer.R
	Macros    []macro.T
	Combos    []combo.T
	Behaviors []behavior.T
}

func init() {
	BASE.Extend(layer.T{
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
		L31: ModX(LCTRL, ModMorph(Kp(MINUS), Kp(PLUS), key.ShiftsCtrls, key.Ctrls)), // row 3
		L32: Kp(Z),
		L33: Kp(X),
		L34: ref2("kpConfig", ZERO, C),
		L35: Kp(V),
		L36: Kp(B),
		L41: MoTo(QUICK, CHAINS), // row 4
		L42: MoX(NUMER, ModMorph(To(MOVER), Kp(UNDER), key.Shifts, nil)),
		L43: Mt(LCTRL, ESCAPE),
	})

	BASE[R11] = Kp(Y) // row 1
	BASE[R12] = Kp(U)
	BASE[R13] = Kp(I)
	BASE[R14] = Kp(O)
	BASE[R15] = Kp(P)
	BASE[R16] = Kp(LBKT)
	BASE[R21] = Kp(H) // row 2
	BASE[R22] = Rmt(LALT, J)
	BASE[R23] = Rmt(LGUI, K)
	BASE[R24] = Rmt(LSHIFT, L)
	BASE[R25] = KpKp(RG(SEMI), SEMI)
	BASE[R26] = KpKp(RG(SQT), SQT)
	BASE[R31] = Kp(N) // row 3
	BASE[R32] = KpKp(RG(M), M)
	BASE[R33] = KpKp(RG(COMMA), COMMA)
	BASE[R34] = KpKp(RG(DOT), DOT)
	BASE[R35] = Kp(SLASH)
	BASE[R36] = Kp(BACKSLASH)
	BASE[R41] = Mt(LCTRL, RETURN) // row 4
	BASE[R42] = Lt(NUMER, SPACE)
	BASE[R43] = Wrap(MoTo(QUICK, CHAINS))

	MOVER.Fill(InitToLevelTrans(BASE))
	MOVER[L43] = To(BASE) // row 4
	MOVER[R21] = Kp(LEFT) // row 2
	MOVER[R22] = Rmt(LALT, DOWN)
	MOVER[R23] = Rmt(LGUI, UP)
	MOVER[R24] = Rmt(LSHIFT, RIGHT)

	NUMER[L11] = Kp(LS(TAB))
	NUMER[L16] = Kp(TILDE)
	NUMER[L21] = Kp(DELETE)      // row 2
	NUMER[L31] = Mt(LCTRL, PLUS) // row 3
	NUMER[L35] = Kp(LS(INSERT))
	NUMER[L42] = Kp(UNDERSCORE)

	NUMER[R11] = Kp(N0) // row 1
	NUMER[R12] = Kp(N1)
	NUMER[R13] = Kp(N2)
	NUMER[R14] = Kp(N3)
	NUMER[R16] = Kp(RBKT)
	NUMER[R21] = ref0("mmEquals") // row 2
	NUMER[R22] = Mt(LALT, N4)
	NUMER[R23] = Mt(LGUI, N5)
	NUMER[R24] = Mt(LSHIFT, N6)
	NUMER[R25] = Kp(COLON)
	NUMER[R26] = ref0("mmQuoteGrave")
	NUMER[R31] = Kp(PLUS) // row 3
	NUMER[R32] = Kp(N7)
	NUMER[R33] = KpKp(RG(COMMA), N8)
	NUMER[R34] = KpKp(RG(DOT), N9)
	NUMER[R35] = Kp(LS(SLASH))
	NUMER[R36] = Kp(PIPE)

	QUICK[L15] = Kp(LG(C_VOL_UP)) // row 1
	QUICK[L16] = Kp(C_VOL_UP)
	QUICK[L25] = Kp(LG(C_VOL_DN)) // row 2
	QUICK[L26] = Kp(C_VOL_DN)
	QUICK[R15] = Kp(PSCRN) // row 1
	QUICK[R16] = Kp(LC(RBKT))
	QUICK[R21] = Kp(HOME) // row 2
	QUICK[R22] = Rmt(LALT, PG_DN)
	QUICK[R23] = Rmt(LGUI, PG_UP)
	QUICK[R24] = Rmt(LSHIFT, END)
	QUICK[R41] = Rmt(LCTRL, F10) // row 4
	QUICK[R42] = Kp(F11)
	QUICK[R43] = Kp(F12)

	SYS[L11] = ref0("bootloader")     // tab
	SYS[L15] = ref0("sys_reset")      // r
	SYS[R12] = ref1("out", "OUT_USB") // u
	SYS[L35] = ref1("out", "OUT_USB") // v - left only backup
	SYS[L36] = ref1("out", "OUT_BLE") // b
	SYS[L25] = ref2("bt", "BT_SEL", "0")
	SYS[L24] = ref2("bt", "BT_SEL", "1")
	SYS[L23] = ref2("bt", "BT_SEL", "2")
	SYS[L22] = ref2("bt", "BT_SEL", "3")
	SYS[L21] = ref2("bt", "BT_SEL", "4")
	SYS[L34] = ref1("bt", "BT_CLR")     // c
	SYS[L33] = ref1("bt", "BT_CLR_ALL") // x
	SYS[R31] = ref1("bt", "BT_CLR_ALL") // n - nuke

	PARENS.Fill(InitToLevelTrans(BASE))
	PARENS[L43] = To(BASE)
	PARENS[L21] = BackspaceDelete()
	PARENS[R24] = XThenLayer(Kp(RIGHT), BASE)

	chain.Add("sdf", Kp(X))
}

func renderKeymap(path string, params Params) {
	tmplPath := path + ".tmpl"
	var funcs = template.FuncMap{"join": strings.Join}
	t := Must(template.New(filepath.Base(path + ".tmpl")).Funcs(funcs).ParseFiles(tmplPath))
	outFile := Must(os.Create(path))
	defer outFile.Close()
	Check(t.Execute(outFile, params))
}

func main() {
	chain.Compile(CHAINS, InitWith(To(BASE)))

	renderKeymap("../config/ergonaut_one.keymap", Params{
		Behaviors: behavior.Render(),
		Macros:    macro.Render(),
		Combos:    combo.Render(),
		Layers:    layer.Render(),
	})
}
