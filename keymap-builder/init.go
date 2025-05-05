package main

import (
	"keyboard/chain"
	"keyboard/combo"
	. "keyboard/instance"
	"keyboard/key"
	. "keyboard/key/keys"
	"keyboard/layer"
	. "keyboard/layout"
	"keyboard/ref"
	"keyboard/rowcol"
)

const (
	comboBothSidesTimeout = 50
	comboBothSidesIdle    = 150
)

var (
	BASE   = layer.New("BASE", InitWith(Trans))
	MOVER  = layer.New("MOVER", InitWith(Trans))
	NUMER  = layer.New("NUMER", InitWith(Trans))
	QUICK  = layer.New("QUICK", InitWith(Trans))
	REPEAT = layer.New("REPEAT", InitWith(Trans))
	SYS    = layer.New("SYS", InitWith(None))
	PARENS = layer.New("PARENS", InitWith(Trans))
	CHAINS = layer.New(chain.Name(""), InitWith(To(BASE)))
)

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
		L34: XKp(Text("XdgConfig", `"$HOME/.config"/`), C),
		L35: Kp(V),
		L36: Kp(B),
		L41: MoTo(QUICK, CHAINS), // row 4
		L42: MoX(NUMER, ModMorph(To(MOVER), Kp(UNDER), key.Shifts, nil)),
		L43: Mt(LCTRL, ESCAPE),
		// BASE RIGHT
		R11: Kp(Y), // row 1
		R12: Kp(U),
		R13: Kp(I),
		R14: Kp(O),
		R15: Kp(P),
		R16: Kp(LBKT),
		R21: Kp(H), // row 2
		// R22: Rmt(LALT, J),
		// R23: Rmt(LGUI, K),
		// R24: Rmt(LSHIFT, L),
		R22: Mt(LALT, J),
		R23: Mt(LGUI, K),
		R24: Mt(LSHIFT, L),
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
		R43: Wrap(MoTo(QUICK, CHAINS)),
	})

	MOVER.Fill(InitLevelOffTrans(MOVER, BASE))
	MOVER.Extend(layer.T{
		L43: To(BASE),
		R21: Kp(LEFT),
		// R22: Rmt(LALT, DOWN),
		// R23: Rmt(LGUI, UP),
		// R24: Rmt(LSHIFT, RIGHT),
		R22: Mt(LALT, DOWN),
		R23: Mt(LGUI, UP),
		R24: Mt(LSHIFT, RIGHT),
	})

	NUMER.Extend(layer.T{
		L11: Kp(LS(TAB)),
		L16: Kp(TILDE),
		L21: Kp(DELETE),      // row 2
		L31: Mt(LCTRL, PLUS), // row 3
		L35: Kp(LS(INSERT)),
		L42: Kp(UNDER),
		R11: Kp(N0), // row 1
		R12: Kp(N1),
		R13: Kp(N2),
		R14: Kp(N3),
		R16: Kp(RBKT),
		R21: ModMorph(Kp(EQUAL), Kp(EQUAL), key.Shifts, nil), // row 2
		R22: Mt(LALT, N4),
		R23: Mt(LGUI, N5),
		R24: Mt(LSHIFT, N6),
		R25: Kp(COLON),
		R26: ModMorph(Kp(DQT), Kp(GRAVE), key.Shifts, nil),
		R31: Kp(PLUS), // row 3
		R32: Kp(N7),
		R33: KpKp(RG(COMMA), N8),
		R34: KpKp(RG(DOT), N9),
		R35: Kp(LS(SLASH)),
		R36: Kp(PIPE),
	})

	QUICK.Extend(layer.T{
		L15: Kp(LG(C_VOL_UP)), // row 1
		L16: Kp(C_VOL_UP),
		L25: Kp(LG(C_VOL_DN)), // row 2
		L26: Kp(C_VOL_DN),
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

	SYS.Extend(layer.T{
		L11: ref.Ref0("bootloader"),     // tab
		L15: ref.Ref0("sys_reset"),      // r
		R12: ref.Ref1("out", "OUT_USB"), // u
		L35: ref.Ref1("out", "OUT_USB"), // v - left only backup
		L36: ref.Ref1("out", "OUT_BLE"), // b
		L25: ref.Ref2("bt", "BT_SEL", "0"),
		L24: ref.Ref2("bt", "BT_SEL", "1"),
		L23: ref.Ref2("bt", "BT_SEL", "2"),
		L22: ref.Ref2("bt", "BT_SEL", "3"),
		L21: ref.Ref2("bt", "BT_SEL", "4"),
		L34: ref.Ref1("bt", "BT_CLR"),     // c
		L33: ref.Ref1("bt", "BT_CLR_ALL"), // x
		R31: ref.Ref1("bt", "BT_CLR_ALL"), // n - nuke
	})

	PARENS.Fill(InitLevelOffTrans(PARENS, BASE))
	PARENS.Extend(layer.T{
		L43: To(BASE),
		L21: BackspaceDelete(),
		R24: XThenLayer(Mt(LSHIFT, RIGHT), BASE),
		R41: XThenLayer(ShiftEnter(), BASE),
	})

	chain.Add(CHAINS, InitWith(To(BASE)), map[string]ref.T{
		"sdf": Kp(X),
		"gie": Text("GoIfError", CursorAt("go if%fine", "%")),
	})

	combo.Add(combo.T{
		Name:     "System",
		Refs:     []ref.T{Sll(SYS)},
		Keys:     []rowcol.T{L15, L16},
		IdleMs:   500,
		TimoutMs: 100,
	})

	combo.Add(combo.T{
		Name: "LeftEnter",
		Refs: []ref.T{Kp(RETURN)},
		Keys: []rowcol.T{L42, L43},
	})

	combo.Add(combo.T{
		Name: "LeftSpace",
		Refs: []ref.T{Kp(SPACE)},
		Keys: []rowcol.T{L41, L42},
	})

	combo.Add(combo.T{
		Name: "RightCaps",
		Refs: []ref.T{CapsWord},
		Keys: []rowcol.T{R23, R24},
	})

	combo.Add(combo.T{
		Name: "RightHomeMods",
		Refs: []ref.T{Kp(LS(LA(LWIN)))},
		Keys: []rowcol.T{R22, R23, R24},
	})

	combo.Add(combo.T{
		Name:     "MiddleMouse",
		Refs:     []ref.T{MKp(MCLK)},
		Keys:     []rowcol.T{L34, L35},
		IdleMs:   200,
		TimoutMs: 100,
	})

	combo.Add(combo.T{
		Name:     "Curlies",
		Refs:     []ref.T{Curlies(PARENS)},
		Keys:     []rowcol.T{L25, R22},
		TimoutMs: comboBothSidesTimeout,
		IdleMs:   comboBothSidesIdle,
	})

	combo.Add(combo.T{
		Name:     "Parens",
		Refs:     []ref.T{Parens(PARENS)},
		Keys:     []rowcol.T{L24, R23},
		TimoutMs: comboBothSidesTimeout,
		IdleMs:   comboBothSidesIdle,
	})

	combo.Add(combo.T{
		Name:     "Brackets",
		Refs:     []ref.T{Brackets(PARENS)},
		Keys:     []rowcol.T{L23, R24},
		TimoutMs: comboBothSidesTimeout,
		IdleMs:   comboBothSidesIdle,
	})

	combo.Add(combo.T{
		Name:     "DoubleQuotes",
		Refs:     []ref.T{DoubleQuotes(PARENS)},
		Keys:     []rowcol.T{L15, R12},
		TimoutMs: comboBothSidesTimeout,
		IdleMs:   comboBothSidesIdle,
	})

	combo.Add(combo.T{
		Name:     "SingleQuotes",
		Refs:     []ref.T{SingleQuotes(PARENS)},
		Keys:     []rowcol.T{L14, R13},
		TimoutMs: comboBothSidesTimeout,
		IdleMs:   comboBothSidesIdle,
	})

	combo.Add(combo.T{
		Name:     "BackQuotes",
		Refs:     []ref.T{HoldTap(Text("CodeBlock", CursorAt("```%```", "%")), BackQuotes(PARENS))},
		Keys:     []rowcol.T{L13, R14},
		TimoutMs: comboBothSidesTimeout,
		IdleMs:   comboBothSidesIdle,
	})

	// combo.Add(combo.T{
	// 	Name:     "CodeQuotes",
	// 	Refs:     []ref.T{},
	// 	Keys:     []rowcol.T{L16, R11},
	// 	TimoutMs: comboBothSidesTimeout,
	// 	IdleMs:   comboBothSidesIdle,
	// })
}
