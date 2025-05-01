package main

import (
	"keyboard/combo"
	. "keyboard/instance"
	. "keyboard/key/keys"
	"keyboard/ref"
	"keyboard/rowcol"
)

func init() {
	combo.Add(combo.T{
		Name:   "System",
		Refs:   []ref.T{ref1("sll", SYS)}, // todo
		Keys:   []rowcol.T{L15, L16},
		IdleMs: 500,
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
		Name:     "MiddleMouse",
		Refs:     []ref.T{MKp(MCLK)},
		Keys:     []rowcol.T{L34, L35},
		IdleMs:   200,
		TimoutMs: 100,
	})

	combo.Add(combo.T{
		Name:     "Curlies",
		Refs:     []ref.T{Curlies()},
		Keys:     []rowcol.T{L25, R22},
		IdleMs:   250,
		TimoutMs: 50,
	})

	combo.Add(combo.T{
		Name:     "Parens",
		Refs:     []ref.T{Parens()},
		Keys:     []rowcol.T{L24, R23},
		IdleMs:   250,
		TimoutMs: 50,
	})

	combo.Add(combo.T{
		Name:     "Brackets",
		Refs:     []ref.T{Brackets()},
		Keys:     []rowcol.T{L23, R24},
		IdleMs:   250,
		TimoutMs: 50,
	})

	combo.Add(combo.T{
		Name:     "DoubleQuotes",
		Refs:     []ref.T{DoubleQuotes()},
		Keys:     []rowcol.T{L15, R12},
		IdleMs:   250,
		TimoutMs: 50,
	})

	combo.Add(combo.T{
		Name:     "SingleQuotes",
		Refs:     []ref.T{SingleQuotes()},
		Keys:     []rowcol.T{L14, R13},
		IdleMs:   250,
		TimoutMs: 50,
	})

	combo.Add(combo.T{
		Name:     "BackQuotes",
		Refs:     []ref.T{BackQuotes()},
		Keys:     []rowcol.T{L13, R14},
		IdleMs:   250,
		TimoutMs: 50,
	})

	combo.Add(combo.T{
		Name:     "CodeQuotes",
		Refs:     []ref.T{ref0("mdCode")}, // todo
		Keys:     []rowcol.T{L16, R11},
		IdleMs:   250,
		TimoutMs: 50,
	})
}
