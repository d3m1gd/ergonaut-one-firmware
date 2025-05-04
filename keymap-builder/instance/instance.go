package instance

import (
	"fmt"
	"slices"
	"strings"

	"keyboard/behavior"
	"keyboard/key"
	. "keyboard/key/keys"
	"keyboard/layer"
	. "keyboard/layout"
	"keyboard/macro"
	"keyboard/ref"
	"keyboard/rowcol"
	. "keyboard/util"
)

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

func Lt(l layer.T, tap key.T) ref.T {
	return ref2("lt", l.Name(), tap)
}

func To(l layer.T) ref.T {
	return ref1("to", l.Name())
}

func Mo(l layer.T) ref.T {
	return ref1("mo", l.Name())
}

func Mt(mod, tap key.T) ref.T {
	return ref2("mt", mod, tap)
}

func Kp(k key.T) ref.T {
	return ref1("kp", k)
}

func MKp(tap key.T) ref.T {
	return ref1("mkp", tap)
}

func Rmt(mod, tap key.T) ref.T {
	name := "rmt"
	keys := Map(slices.Concat(LLLL, []rowcol.T{R22, R23, R24, R41, R42, R43}), rowcol.ToSerial)
	behavior.Add(behavior.T{
		Name:  name,
		Label: "RightModTap",
		Cells: behavior.TypeHoldTap.Cells,
		Type:  behavior.TypeHoldTap.Name,
		Refs:  []ref.T{ref0("kp"), ref0("kp")},
		Props: behavior.Props{
			"hold-trigger-key-positions": keys,
			"hold-trigger-on-release":    true,
			"flavor":                     "tap-preferred",
			"tapping-term-ms":            250,
			"quick-tap-ms":               200,
		},
	})

	return ref.Filled(name, behavior.TypeHoldTap.Cells, mod, tap)
}

func Sll(l layer.T) ref.T {
	name := "sll"
	behavior.Add(behavior.T{
		Name:  name,
		Label: "StickyLayerLong",
		Cells: behavior.TypeStickyKey.Cells,
		Type:  behavior.TypeStickyKey.Name,
		Refs:  []ref.T{ref0("mo")},
		Props: behavior.Props{
			"release-after-ms": 2000,
			"quick-release":    true,
		},
	})

	return ref.Filled(name, behavior.TypeStickyKey.Cells, l)
}
func KpKp(a, b key.T) ref.T {
	name := "kpkp"
	tnr := TapNoRepeat(a).Strip() // instantiate macro
	behavior.Add(behavior.T{
		Name:  name,
		Label: "KeyPressKepPress",
		Cells: behavior.TypeHoldTap.Cells,
		Type:  behavior.TypeHoldTap.Name,
		Refs:  []ref.T{tnr, ref0("kp")},
		Props: behavior.Props{
			"flavor":          "tap-preferred",
			"tapping-term-ms": 200,
			"quick-tap-ms":    200,
		},
	})

	return ref.Filled(name, behavior.TypeHoldTap.Cells, a, b)
}

func XKp(r ref.T, k key.T) ref.T {
	name := r.Show() + "Kp"
	behavior.Add(behavior.T{
		Name:  name,
		Label: r.Show() + "KepPress",
		Cells: behavior.TypeHoldTap.Cells,
		Type:  behavior.TypeHoldTap.Name,
		Refs:  []ref.T{r, ref0("kp")},
		Props: behavior.Props{
			"flavor":          "tap-preferred",
			"tapping-term-ms": 350,
			"quick-tap-ms":    200,
		},
	})

	return ref.Filled(name, behavior.TypeHoldTap.Cells, ZERO, k)
}

func MoTo(mo, to layer.T) ref.T {
	refs := []ref.T{ref0("mo"), ref0("to")}
	name := "moto"
	behavior.Add(behavior.T{
		Name:  name,
		Label: "MomentaryTo",
		Type:  behavior.TypeHoldTap.Name,
		Cells: behavior.TypeHoldTap.Cells,
		Refs:  refs,
		Props: behavior.Props{
			"flavor":          "balanced",
			"tapping-term-ms": 300,
		},
	})

	return ref.Filled(name, behavior.TypeHoldTap.Cells, mo, to)
}

func MoX(mo layer.T, x ref.T) ref.T {
	refs := []ref.T{ref0("mo"), x}
	name := "mo" + x.Name
	behavior.Add(behavior.T{
		Name:  name,
		Label: fmt.Sprintf("Mom%s", x.Show()),
		Type:  behavior.TypeHoldTap.Name,
		Cells: behavior.TypeHoldTap.Cells,
		Refs:  refs,
		Props: behavior.Props{
			"flavor":          "balanced",
			"tapping-term-ms": 300,
		},
	})

	return ref.Filled(name, behavior.TypeHoldTap.Cells, mo, ZERO)
}

func ModX(mod key.T, x ref.T) ref.T {
	name := "m" + x.Show()
	behavior.Add(behavior.T{
		Name:  name,
		Label: "Mod" + x.Show(),
		Type:  behavior.TypeHoldTap.Name,
		Cells: behavior.TypeHoldTap.Cells,
		Refs:  []ref.T{ref0("kp"), x},
		Props: behavior.Props{
			"flavor":          "tap-preferred",
			"tapping-term-ms": 200,
			"quick-tap-ms":    200,
		},
	})

	return ref.Filled(name, behavior.TypeHoldTap.Cells, mod, ZERO)
}

func ModMorph(a, b ref.T, mods []key.Mod, keep []key.Mod) ref.T {
	refs := []ref.T{a, b}
	name := "MM" + a.Show() + b.Show()
	props := behavior.Props{
		"mods": mods,
	}
	if len(keep) > 0 {
		props["keep-mods"] = keep
	}

	behavior.Add(behavior.T{
		Name:  name,
		Label: "ModMorph" + a.Show() + b.Show(),
		Cells: behavior.TypeModMorph.Cells,
		Type:  behavior.TypeModMorph.Name,
		Refs:  refs,
		Props: props,
	})

	return ref.Filled(name, behavior.TypeModMorph.Cells)
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

func Wrap(r ref.T) ref.T {
	name := fmt.Sprintf("W%s", r.Name)
	params := macroParams(len(r.Args()))
	refs := []ref.T{}
	refs = append(refs, Press)
	refs = append(refs, params...)
	refs = append(refs, macro.Placeholder(r))
	refs = append(refs, Pause)
	refs = append(refs, Release)
	refs = append(refs, params...)
	refs = append(refs, macro.Placeholder(r))
	macro.Add(macro.T{
		Name:  name,
		Label: fmt.Sprintf("Wrap%s", r.Name),
		Cells: len(r.Args()),
		Refs:  refs,
	})

	return ref.RefN(name, Map(r.Args(), ToAny))
}

func BackspaceDelete() ref.T {
	name := "bspcdel"
	macro.Add(macro.T{
		Name:  name,
		Label: "BackspaceDelete",
		Cells: 0,
		Refs:  []ref.T{Kp(BSPC), Kp(DEL)},
	})

	return ref0(name)
}

func Parens(l layer.T) ref.T {
	return OpenCloseMacro("parens", LPAR, RPAR, l)
}

func Brackets(l layer.T) ref.T {
	return OpenCloseMacro("brackets", LBKT, RBKT, l)
}

func Curlies(l layer.T) ref.T {
	return OpenCloseMacro("curlies", LBRC, RBRC, l)
}

func DoubleQuotes(l layer.T) ref.T {
	return OpenCloseMacro("dquotes", DQT, DQT, l)
}

func SingleQuotes(l layer.T) ref.T {
	return OpenCloseMacro("squotes", SQT, SQT, l)
}

func BackQuotes(l layer.T) ref.T {
	return OpenCloseMacro("bquotes", GRAVE, GRAVE, l)
}

func OpenCloseMacro(name string, left, right key.T, l layer.T) ref.T {
	macro.Add(macro.T{
		Name:  name,
		Label: fmt.Sprintf("OpenClose_%s", name),
		Cells: 0,
		// Refs:  []ref.T{Kp(left), Kp(right), Kp(LEFT)},
		Refs: []ref.T{Kp(left), Kp(right), Kp(LEFT), To(l)}, // todo
	})

	return ref0(name)
}

func XThenTrans(r ref.T, l layer.T, rc rowcol.T) ref.T {
	name := fmt.Sprintf("%sTrans%s", r.Show(), l)
	macro.Add(macro.T{
		Name:  name,
		Label: fmt.Sprintf("%s_Trans_%s", r.Name, l),
		Cells: 1,
		Refs:  []ref.T{Param11, macro.Placeholder(r), To(l), l[rc]},
	})

	return ref.RefN(name, Map(r.Args(), ToAny))
}

func XThenLayer(r ref.T, l layer.T) ref.T {
	name := fmt.Sprintf("%s%s", r.Name, l)
	macro.Add(macro.T{
		Name:  name,
		Label: name,
		Cells: 1,
		Refs:  []ref.T{Param11, macro.Placeholder(r), To(l)},
	})

	return ref.RefN(name, Map(r.Args(), ToAny))
}

func TapNoRepeat(k key.T) ref.T {
	name := "TapNoRepeat"
	macro.Add(macro.T{
		Name:  name,
		Label: name,
		Cells: 1,
		Refs:  []ref.T{Param11, macro.Placeholder(Kp(k)), Pause},
	})

	return ref1(name, Kp(k))
}

func InitWith(b ref.T) func(layer.T) {
	return layer.InitBy(func(rowcol.T) ref.T { return b })
}

func InitToLevelTrans(l layer.T) func(layer.T) {
	return layer.InitBy(func(rc rowcol.T) ref.T {
		name := fmt.Sprintf("to%d%s", l.Index(), rc)
		macro.Add(macro.T{
			Name:  name,
			Label: fmt.Sprintf("To%s%s", l.Name(), rc.Pretty()),
			Cells: 0,
			Refs:  []ref.T{Press, l[rc], Pause, Release, l[rc], Tap, To(l)},
		})
		return ref0(name)
	})
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
