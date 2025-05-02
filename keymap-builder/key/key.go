package key

import (
	"cmp"
	"fmt"
)

type T = Key

type Key string
type Mod string

func (k Key) Less(other Key) int {
	return cmp.Compare(k, other)
}

func From(b byte) Key {
	if b >= 'A' && b <= 'Z' {
		return Key(b)
	}

	switch b {
	// case '\a':
	// case '\f':
	// case '\r':
	// case '\v':
	case ' ':
		return Key("SPACE")
	case '\b':
		return Key("LEFT")
	case '\n':
		return Key("RETURN")
	case '\t':
		return Key("TAB")
	}

	panic(fmt.Sprintf("unhandled key: '%c'", b))
}

const MOD_LCTL Mod = "MOD_LCTL"
const MOD_LSFT Mod = "MOD_LSFT"
const MOD_LALT Mod = "MOD_LALT"
const MOD_LGUI Mod = "MOD_LGUI"
const MOD_RCTL Mod = "MOD_RCTL"
const MOD_RSFT Mod = "MOD_RSFT"
const MOD_RALT Mod = "MOD_RALT"
const MOD_RGUI Mod = "MOD_RGUI"

var Shifts = []Mod{MOD_LSFT, MOD_RSFT}
var Ctrls = []Mod{MOD_LCTL, MOD_RCTL}
var ShiftsCtrls = []Mod{MOD_LSFT, MOD_RSFT, MOD_LCTL, MOD_RCTL}
