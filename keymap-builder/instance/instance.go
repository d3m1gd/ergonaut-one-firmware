package instance

import (
	"fmt"
	"slices"
	"strings"

	"keyboard/behavior"
	. "keyboard/key"
	"keyboard/layer"
	. "keyboard/layout"
	"keyboard/macro"
	"keyboard/ref"
	"keyboard/rowcol"
	. "keyboard/util"
)

type Layer = layer.Layer

var LayerPlaceholder = layer.Layer{
	Name: "MACRO_PLACEHOLDER",
}

var (
	ref0 = ref.Ref0
	ref1 = ref.Ref1
	ref2 = ref.Ref2

	Trans    = ref0("trans")
	None     = ref0("none")
	CapsWord = ref0("caps_word")
	Tap      = ref0("macro_tap")
	Press    = ref0("macro_press")
	Release  = ref0("macro_release")
	Pause    = ref0("macro_pause_for_release")

	Param11 = macroParamBuilder(1, 1)
	Param12 = macroParamBuilder(1, 2)
	Param21 = macroParamBuilder(2, 1)
	Param22 = macroParamBuilder(2, 2)
)

func macroParamBuilder(a, b int) ref.T {
	return ref.Ref0(fmt.Sprintf("macro_param_%dto%d", a, b))
}

func Lt(l Layer, tap Key) ref.T {
	return ref2("lt", l.Name, tap)
}

func To(l Layer) ref.T {
	return ref1("to", l.Name)
}

func Mo(l Layer) ref.T {
	return ref1("mo", l.Name)
}

func Mt(mod, tap Key) ref.T {
	return ref2("mt", mod, tap)
}

func Kp(k Key) ref.T {
	return ref1("kp", k)
}

func MKp(tap Key) ref.T {
	return ref1("mkp", tap)
}

func Rmt(mod, tap Key) ref.T {
	keys := Map(slices.Concat(LLLL, []rowcol.T{R22, R23, R24, R41, R42, R43}), rowcol.ToSerial)
	return HoldTapOpts(Kp(mod), Kp(tap), "rmt", "RightModTap", behavior.Props{
		"hold-trigger-key-positions": keys,
		"hold-trigger-on-release":    true,
		"flavor":                     "tap-preferred",
		"tapping-term-ms":            250,
		"quick-tap-ms":               200,
	})
}

func HoldTap(h, t ref.T) ref.T {
	return HoldTapOpts(h, t, "ht", "HoldTap", behavior.Props{
		"flavor":          "tap-preferred",
		"tapping-term-ms": 200,
		"quick-tap-ms":    200,
	})
}

func HoldTapOpts(h, t ref.T, name, label string, properties behavior.Props) ref.T {
	return behavior.AddY(behavior.T{
		Name:  name,
		Label: label,
		Type:  behavior.TypeHoldTap,
		Refs:  []ref.T{h, t},
		Props: properties,
	})
}

func Sll(l Layer) ref.T {
	return behavior.AddY(behavior.T{
		Name:  "sll",
		Label: "StickyLayerLong",
		Type:  behavior.TypeStickyKey,
		Refs:  []ref.T{Mo(l)},
		Props: behavior.Props{
			"release-after-ms": 2000,
			"quick-release":    true,
		},
	})
}

func Sl(l Layer, duration int) ref.T {
	return behavior.AddY(behavior.T{
		Name:  fmt.Sprintf("sll%d", duration),
		Label: fmt.Sprintf("StickyLayer%d", duration),
		Type:  behavior.TypeStickyKey,
		Refs:  []ref.T{Mo(l)},
		Props: behavior.Props{
			"release-after-ms": duration,
			"quick-release":    true,
		},
	})
}

func KpSl(k Key, l Layer, duration int) ref.T {
	name := fmt.Sprintf("KpSl%s%d", l, duration)

	m := macro.Add(macro.T{
		Name:  name,
		Label: name,
		Cells: 1,
		Refs:  []ref.T{Param11, Kp(MACRO_PLACEHOLDER), Sl(l, duration)},
	})

	return m.Invoke(k)
}

func KpKp(a, b Key) ref.T {
	return HoldTap(TapNoRepeat(a), Kp(b))
}

func XKp(r ref.T, k Key) ref.T {
	return HoldTapOpts(r, Kp(k), "xkp", "XKeyPress", behavior.Props{
		"flavor":          "tap-preferred",
		"tapping-term-ms": 350,
		"quick-tap-ms":    200,
	})
}

func MoTo(mo, to Layer) ref.T {
	return HoldTapOpts(Mo(mo), To(to), "moto", "MomentaryTo", behavior.Props{
		"flavor":          "balanced",
		"tapping-term-ms": 300,
	})
}

func MoX(mo Layer, x ref.T) ref.T {
	return HoldTapOpts(Mo(mo), x, "mo", "Momentary", behavior.Props{
		"flavor":          "balanced",
		"tapping-term-ms": 300,
	})
}

func ModX(mod Key, x ref.T) ref.T {
	return HoldTapOpts(Kp(mod), x, "m", "Mod", behavior.Props{
		"flavor":          "tap-preferred",
		"tapping-term-ms": 200,
		"quick-tap-ms":    200,
	})
}

func ModMorph(a, b ref.T, mods []Mod, keep []Mod) ref.T {
	props := behavior.Props{
		"mods": mods,
	}
	if len(keep) > 0 {
		props["keep-mods"] = keep
	}

	return behavior.AddY(behavior.T{
		Name:  "mm",
		Label: "ModMorph",
		Type:  behavior.TypeModMorph,
		Refs:  []ref.T{a, b},
		Props: props,
	})
}

func Off(l Layer) ref.T {
	name := "off"
	behavior.Add(behavior.T{
		Name:  name,
		Label: "Off",
		Type:  behavior.TypeToggleLayer,
		Props: behavior.Props{
			"display-name": "Layer Off",
			"toggle-mode":  "off",
		},
	})
	return ref.Ref1(name, l)
}

func macroParams(n int) []ref.T {
	switch n {
	case 0:
		return []ref.T{}
	case 1:
		return []ref.T{Param11}
	case 2:
		return []ref.T{Param11, Param22}
	}
	panic(fmt.Sprintf("bad n: %d", n))
}

// func Wrap(r ref.T) ref.T {
// 	name := fmt.Sprintf("W%s", r.Name)
// 	params := macroParams(len(r.Args()))
// 	refs := []ref.T{}
// 	refs = append(refs, Press)
// 	refs = append(refs, params...)
// 	refs = append(refs, macro.Placeholder(r))
// 	refs = append(refs, Pause)
// 	refs = append(refs, Release)
// 	refs = append(refs, params...)
// 	refs = append(refs, macro.Placeholder(r))
// 	macro.Add(macro.T{
// 		Name:  name,
// 		Label: fmt.Sprintf("Wrap%s", r.Name),
// 		Cells: len(r.Args()),
// 		Refs:  refs,
// 	})
//
// 	return ref.RefN(name, Map(r.Args(), ToAny))
// }

func BackspaceDelete() ref.T {
	name := "bspcdel"
	macro.Add(macro.T{
		Name:  name,
		Label: "BackspaceDelete",
		Cells: 0,
		Refs:  []ref.T{Kp(DEL), Kp(BSPC)},
	})

	return ref0(name)
}

func Parens(l Layer) ref.T {
	return OpenCloseMacro("parens", LPAR, RPAR, l)
}

func Brackets(l Layer) ref.T {
	return OpenCloseMacro("brackets", LBKT, RBKT, l)
}

func Curlies(l Layer) ref.T {
	return OpenCloseMacro("curlies", LBRC, RBRC, l)
}

func DoubleQuotes(l Layer) ref.T {
	return OpenCloseMacro("dquotes", DQT, DQT, l)
}

func SingleQuotes(l Layer) ref.T {
	return OpenCloseMacro("squotes", SQT, SQT, l)
}

func BackQuotes(l Layer) ref.T {
	return OpenCloseMacro("bquotes", GRAVE, GRAVE, l)
}

func OpenCloseMacro(name string, left, right Key, l Layer) ref.T {
	macro.Add(macro.T{
		Name:  name,
		Label: fmt.Sprintf("OpenClose_%s", name),
		Cells: 0,
		Refs:  []ref.T{Kp(left), Kp(right), Kp(LEFT), Sll(l)},
	})

	return ref0(name)
}

// todo delme
// func MSll(mod Key, l Layer) ref.T {
// 	name := "ModSll"
// 	macro.Add(macro.T{
// 		Name:  name,
// 		Label: name,
// 		Cells: 0,
// 		Refs:  []ref.T{Param11, macro.Placeholder(Kp(mod)), Sll(l)},
// 	})
//
// 	return ref2(name, mod, l)
// }

func ReRet() ref.T {
	return macro.Add(macro.T{
		Name:  "ReRet",
		Cells: 0,
		Refs:  []ref.T{Kp(RET), Kp(UP), Kp(END), Kp(RET)},
	}).Invoke()
}

func TapNoRepeat(k Key) ref.T {
	name := "TapNoRepeat"
	return macro.Add(macro.T{
		Name:  name,
		Label: name,
		Cells: 1,
		Refs:  []ref.T{Param11, Kp(MACRO_PLACEHOLDER), Pause},
	}).Invoke(k)
}

func InitWith(b ref.T) func(Layer) {
	return layer.InitBy(func(rowcol.T) ref.T { return b })
}

func InitOffTrans(l Layer, base Layer) func(Layer) {
	return layer.InitBy(func(rc rowcol.T) ref.T {
		name := fmt.Sprintf("off%s", rc)
		key := base.Cells[rc]
		macro.Add(macro.T{
			Name:  name,
			Label: fmt.Sprintf("Off%s", rc.Pretty()),
			Cells: 1,
			Refs:  []ref.T{Press, key, Pause, Release, key, Tap, Param11, Off(LayerPlaceholder)}, // todo macro strip
		})
		return ref1(name, l)
	})
}

func OffX(l Layer, r ref.T) ref.T {
	if r.Name == "kp" {
		return OffKey(l, r.Fields[0].(Key))
	}
	name := "Off" + r.Show()
	macro.Add(macro.T{
		Name:  name,
		Label: name,
		Cells: 1,
		Refs:  []ref.T{Press, r, Pause, Release, r, Tap, Param11, Off(LayerPlaceholder)},
	})

	return ref1(name, l)
}

func OffKey(l Layer, k Key) ref.T {
	name := "OffKey"
	macro.Add(macro.T{
		Name:  name,
		Label: name,
		Cells: 2,
		Refs:  []ref.T{Press, Param21, Kp(MACRO_PLACEHOLDER), Pause, Release, Param21, Kp(MACRO_PLACEHOLDER), Tap, Param11, Off(LayerPlaceholder)},
	})

	return ref2(name, l, k)
}

func Text(name string, keys string) ref.T {
	macro.Add(macro.T{
		Name:  name,
		Label: name,
		Cells: 0,
		Refs:  MapString(keys, func(b byte) ref.T { return Kp(From(b)) }),
	})

	return ref0(name)
}

func CursorAt(str, marker string) string {
	subs := strings.Split(str, marker)
	Panicif(len(subs) != 2, "want exactly one marker (%s): %s", marker, str)
	left := subs[0]
	right := strings.Split(subs[1], "")
	slices.Reverse(right)

	return left + strings.Join(right, "\b") + "\b"
}
