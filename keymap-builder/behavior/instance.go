package behavior

import (
	"fmt"
	"slices"
	"strings"

	"keyboard/key"
	"keyboard/layer"
	. "keyboard/layout"
	"keyboard/macro"
	"keyboard/ref"
	"keyboard/rowcol"
	"keyboard/util"
)

type Layer = layer.Layer

var LayerPlaceholder = Layer{
	Name: "MACRO_PLACEHOLDER",
}

var KeyPlaceholder = key.MACRO_PLACEHOLDER

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

	Param11 = ref.Ref0("macro_param_1to1")
	Param12 = ref.Ref0("macro_param_1to2")
	Param21 = ref.Ref0("macro_param_2to1")
	Param22 = ref.Ref0("macro_param_2to2")
)

func Lt(l Layer, tap key.Key) ref.T {
	return ref2("lt", l.Name, tap)
}

func To(l Layer) ref.T {
	return ref1("to", l.Name)
}

func Mo(l Layer) ref.T {
	return ref1("mo", l.Name)
}

func Mt(mod, tap key.Key) ref.T {
	return ref2("mt", mod, tap)
}

func Kp(k key.Key) ref.T {
	return ref1("kp", k)
}

func MKp(tap key.Key) ref.T {
	return ref1("mkp", tap)
}

func Rmt(mod, tap key.Key) ref.T {
	keys := util.Map(slices.Concat(LLLL, []rowcol.T{R22, R23, R24, R41, R42, R43}), rowcol.ToSerial)
	return HoldTapOpts(Kp(mod), Kp(tap), "rmt", Props{
		"hold-trigger-key-positions": keys,
		"hold-trigger-on-release":    true,
		"flavor":                     "tap-preferred",
		"tapping-term-ms":            250,
		"quick-tap-ms":               200,
	})
}

func HoldTap(h, t ref.T) ref.T {
	return HoldTapOpts(h, t, "ht", Props{
		"flavor":          "tap-preferred",
		"tapping-term-ms": 200,
		"quick-tap-ms":    200,
	})
}

func HoldTapFast(h, t ref.T) ref.T {
	return HoldTapOpts(h, t, "htf", Props{
		"flavor":          "tap-preferred",
		"tapping-term-ms": 100,
		"quick-tap-ms":    100,
	})
}

func HoldTapStickyLong(k key.Key) ref.T {
	return HoldTap(Kp(k), Skl(k))
}

func HoldTapOpts(h, t ref.T, name string, properties Props) ref.T {
	return Add(Behavior{
		Name:  name,
		Type:  TypeHoldTap,
		Refs:  []ref.T{h, t},
		Props: properties,
	})
}

func Skl(k key.Key) ref.T {
	return Add(Behavior{
		Name: "skl",
		Type: TypeStickyKey,
		Refs: []ref.T{Kp(k)},
		Props: Props{
			"release-after-ms": 9000,
			"quick-release":    true,
			"lazy":             true,
			"ignore-modifiers": true,
		},
	})
}

func Sll(l Layer) ref.T {
	return Add(Behavior{
		Name: "sll",
		Type: TypeStickyKey,
		Refs: []ref.T{Mo(l)},
		Props: Props{
			"release-after-ms": 2000,
			"quick-release":    true,
		},
	})
}

func Sl(l Layer, duration int) ref.T {
	return Add(Behavior{
		Name: fmt.Sprintf("sll%d", duration),
		Type: TypeStickyKey,
		Refs: []ref.T{Mo(l)},
		Props: Props{
			"release-after-ms": duration,
			"quick-release":    true,
		},
	})
}

func XSl(x ref.T, l Layer, duration int) ref.T {
	sl := Sl(LayerPlaceholder, duration)
	return macro.Add(macro.T{
		Name:  fmt.Sprintf("%s%s", x.Name, sl.Name),
		Cells: 1,
		Refs:  []ref.T{x, Param11, sl},
	}).Invoke(l)
}

func KpSl(k key.Key, l Layer, duration int) ref.T {
	kp := Kp(KeyPlaceholder)
	sl := Sl(LayerPlaceholder, duration)
	return macro.Add(macro.T{
		Name:  fmt.Sprintf("%s%s", kp.Name, sl.Name),
		Cells: 2,
		Refs:  []ref.T{Param11, kp, Param21, sl},
	}).Invoke(k, l)
}

func KpKp(a, b key.Key) ref.T {
	return HoldTap(TapNoRepeat(a), Kp(b))
}

func XKp(r ref.T, k key.Key) ref.T {
	return HoldTapOpts(r, Kp(k), "xkp", Props{
		"flavor":          "tap-preferred",
		"tapping-term-ms": 350,
		"quick-tap-ms":    200,
	})
}

func MoTo(mo, to Layer) ref.T {
	return HoldTapOpts(Mo(mo), To(to), "moto", Props{
		"flavor":          "balanced",
		"tapping-term-ms": 300,
	})
}

func MoX(mo Layer, x ref.T) ref.T {
	return HoldTapOpts(Mo(mo), x, "mo", Props{
		"flavor":          "balanced",
		"tapping-term-ms": 300,
	})
}

func ModX(mod key.Key, x ref.T) ref.T {
	return HoldTapOpts(Kp(mod), x, "m", Props{
		"flavor":          "tap-preferred",
		"tapping-term-ms": 200,
		"quick-tap-ms":    200,
	})
}

func ModMorph(a, b ref.T, mods []key.Mod, keep []key.Mod) ref.T {
	props := Props{
		"mods": mods,
	}
	if len(keep) > 0 {
		props["keep-mods"] = keep
	}

	return Add(Behavior{
		Name:  "mm",
		Type:  TypeModMorph,
		Refs:  []ref.T{a, b},
		Props: props,
	})
}

func Off(l Layer) ref.T {
	name := "off"
	Add(Behavior{
		Name: name,
		Type: TypeToggleLayer,
		Props: Props{
			"display-name": "Layer Off",
			"toggle-mode":  "off",
		},
	})
	return ref.Ref1(name, l)
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
	return macro.Add(macro.T{
		Name: "BackspaceDelete",
		Refs: []ref.T{Kp(key.DEL), Kp(key.BSPC)},
	}).Invoke()
}

func Parens(l Layer) ref.T {
	return OpenCloseMacro("parens", key.LPAR, key.RPAR, l)
}

func Brackets(l Layer) ref.T {
	return OpenCloseMacro("brackets", key.LBKT, key.RBKT, l)
}

func Curlies(l Layer) ref.T {
	return OpenCloseMacro("curlies", key.LBRC, key.RBRC, l)
}

func DoubleQuotes(l Layer) ref.T {
	return OpenCloseMacro("dquotes", key.DQT, key.DQT, l)
}

func SingleQuotes(l Layer) ref.T {
	return OpenCloseMacro("squotes", key.SQT, key.SQT, l)
}

func BackQuotes(l Layer) ref.T {
	return OpenCloseMacro("bquotes", key.GRAVE, key.GRAVE, l)
}

func OpenCloseMacro(name string, left, right key.Key, l Layer) ref.T {
	return macro.Add(macro.T{
		Name: name,
		Refs: []ref.T{Kp(left), Kp(right), Kp(key.LEFT), Sll(l)},
	}).Invoke()
}

func ReRet() ref.T {
	return macro.Add(macro.T{
		Name: "ReRet",
		Refs: []ref.T{Kp(key.RET), Kp(key.UP), Kp(key.END), Kp(key.RET)},
	}).Invoke()
}

func TapNoRepeat(k key.Key) ref.T {
	return macro.Add(macro.T{
		Name:  "TapNoRepeat",
		Cells: 1,
		Refs:  []ref.T{Param11, Kp(KeyPlaceholder), Pause},
	}).Invoke(k)
}

func TapReliableNoRepeat(k key.Key, mods ...key.Key) ref.T {
	return macro.Add(macro.T{
		Name:  "TapReliableNoRepeat",
		Cells: 1,
		Refs:  []ref.T{Param11, Reliable(KeyPlaceholder, mods...), Pause},
	}).Invoke(k)
}

func InitOffTrans(l Layer, base Layer) func(Layer) {
	return layer.InitBy(func(rc rowcol.T) ref.T {
		key := base.Cells[rc]
		return macro.Add(macro.T{
			Name:  fmt.Sprintf("off%s", rc),
			Cells: 1,
			Refs:  []ref.T{Press, key, Pause, Release, key, Tap, Param11, Off(LayerPlaceholder)},
		}).Invoke(l)
	})
}

func OffX(l Layer, r ref.T) ref.T {
	if r.Name == "kp" {
		return OffKey(l, r.Fields[0].(key.Key))
	}
	return macro.Add(macro.T{
		Name:  "Off" + r.Show(),
		Cells: 1,
		Refs:  []ref.T{Press, r, Pause, Release, r, Tap, Param11, Off(LayerPlaceholder)},
	}).Invoke(l)
}

func OffKey(l Layer, k key.Key) ref.T {
	return macro.Add(macro.T{
		Name:  "OffKey",
		Cells: 2,
		Refs:  []ref.T{Press, Param21, Kp(KeyPlaceholder), Pause, Release, Param21, Kp(KeyPlaceholder), Tap, Param11, Off(LayerPlaceholder)},
	}).Invoke(l, k)
}

func Text(name string, keys string) ref.T {
	return macro.Add(macro.T{
		Name:  name,
		Cells: 0,
		Refs:  util.MapString(keys, func(b byte) ref.T { return Kp(key.From(b)) }),
	}).Invoke()
}

func Reliable(k key.Key, mods ...key.Key) ref.T {
	name := fmt.Sprintf("Reliable%s", strings.Join(util.Map(mods, func(k key.Key) string {
		return string(k)
	}), ""))

	refs := []ref.T{Press}

	for _, m := range mods {
		refs = append(refs, Kp(m))
	}

	refs = append(refs, Tap)
	refs = append(refs, Param11)
	refs = append(refs, Kp(KeyPlaceholder))
	refs = append(refs, Release)

	for _, m := range mods {
		refs = append(refs, Kp(m))
	}

	return macro.Add(macro.T{
		Name:  name,
		Cells: 1,
		Refs:  refs,
	}).Invoke(k)
}

func CursorAt(str, marker string) string {
	subs := strings.Split(str, marker)
	util.Panicif(len(subs) != 2, "want exactly one marker (%s): %s", marker, str)
	left := subs[0]
	right := strings.Split(subs[1], "")
	slices.Reverse(right)

	return left + strings.Join(right, "\b") + "\b"
}
