package key

import (
	"cmp"
)

// Key ...
//
// Deprecated: use something else
type Key string

// Mod ...
//
// Deprecated: use something else
type Mod string

func (k Key) Less(other Key) int {
	return cmp.Compare(k, other)
}

const (
	LControlMod Mod = "MOD_LCTL"
	LShiftMod   Mod = "MOD_LSFT"
	LAltMod     Mod = "MOD_LALT"
	LGuiMod     Mod = "MOD_LGUI"
	RControlMod Mod = "MOD_RCTL"
	RShiftMod   Mod = "MOD_RSFT"
	RAltMod     Mod = "MOD_RALT"
	RGuiMod     Mod = "MOD_RGUI"
)

var Shifts = []Mod{LShiftMod, RShiftMod}
var Ctrls = []Mod{LControlMod, RControlMod}
var ShiftsAndCtrls = []Mod{LShiftMod, RShiftMod, LControlMod, RControlMod}
