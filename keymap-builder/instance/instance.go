package instance

import (
	"fmt"
	"slices"
	"strings"

	"keyboard/behavior"
	"keyboard/key"
	. "keyboard/key/keys"
	"keyboard/layer"
	"keyboard/macro"
	"keyboard/ref"
	"keyboard/rowcol"
	. "keyboard/util"
)

var ref0 = ref.Ref0
var ref1 = ref.Ref1
var ref2 = ref.Ref2

var Trans = ref0("trans")
var None = ref0("none")
var CapsWord = ref0("caps_word")

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
	return ref2("rmt", mod, tap)
}

func KpKp(a, b key.T) ref.T {
	name := "kpkp"
	_ = TapNoRepeat(a) // instantiate macro
	behavior.Add(behavior.T{
		Name:  name,
		Label: "kpkp",
		Cells: behavior.TypeHoldTap.Cells,
		Type:  behavior.TypeHoldTap.Name,
		Refs:  []ref.T{ref0("TapNoRepeat"), ref0("kp")},
		Props: behavior.Props{
			"flavor":          "tap-preferred",
			"tapping-term-ms": 200,
			"quick-tap-ms":    200,
		},
	})

	return ref.Filled(name, behavior.TypeHoldTap.Cells, a, b)
}

func MoTo(mo, to layer.T) ref.T {
	refs := []ref.T{ref0("mo"), ref0("to")}
	name := "moto"
	behavior.Add(behavior.T{
		Name:  name,
		Label: "Momentary/To",
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
	name := "mo" + x.Show()
	behavior.Add(behavior.T{
		Name:  name,
		Label: "Momentary " + name,
		Type:  behavior.TypeHoldTap.Name,
		Cells: behavior.TypeHoldTap.Cells,
		Refs:  refs,
		Props: behavior.Props{
			"flavor":          "balanced",
			"tapping-term-ms": 300,
		},
	})

	return ref.Filled(name, behavior.TypeHoldTap.Cells, mo)
}

func ModX(mod key.T, x ref.T) ref.T {
	name := "m" + x.Show()
	behavior.Add(behavior.T{
		Name:  name,
		Label: "Mod " + x.Name,
		Type:  behavior.TypeHoldTap.Name,
		Cells: behavior.TypeHoldTap.Cells,
		Refs:  []ref.T{ref0("kp"), x},
		Props: behavior.Props{
			"flavor":          "tap-preferred",
			"tapping-term-ms": 200,
			"quick-tap-ms":    200,
		},
	})

	return ref.Filled(name, behavior.TypeHoldTap.Cells, mod)
}

func ModMorph(a, b ref.T, mods []key.Mod, keep []key.Mod) ref.T {
	refs := []ref.T{a, b}
	name := "mm" + a.Show() + b.Show()
	props := behavior.Props{
		"mods": mods,
	}
	if len(keep) > 0 {
		props["keep-mods"] = keep
	}

	behavior.Add(behavior.T{
		Name:  name,
		Label: "ModMorph " + a.String() + " " + b.String(),
		Cells: behavior.TypeModMorph.Cells,
		Type:  behavior.TypeModMorph.Name,
		Refs:  refs,
		Props: props,
	})

	return ref.Filled(name, behavior.TypeModMorph.Cells)
}

func Wrap(r ref.T) ref.T {
	name := fmt.Sprintf("W%s", r.Name)
	params := macro.MapParams(len(r.Args()))
	refs := []ref.T{}
	refs = append(refs, macro.Press)
	refs = append(refs, params...)
	refs = append(refs, macro.Placeholder(r))
	refs = append(refs, macro.Pause)
	refs = append(refs, macro.Release)
	refs = append(refs, params...)
	refs = append(refs, macro.Placeholder(r))
	macro.Add(macro.T{
		Name:  name,
		Label: fmt.Sprintf("Wrap %s", r.Name),
		Cells: len(r.Args()),
		Refs:  refs,
	})

	return ref.RefN(name, MapToAny(r.Args()))
}

func BackspaceDelete() ref.T {
	name := "bspcdel"
	macro.Add(macro.T{
		Name:  name,
		Label: "Backspace Delete",
		Cells: 0,
		Refs:  []ref.T{Kp(BSPC), Kp(DEL)},
	})

	return ref0(name)
}

func Parens() ref.T {
	return OpenCloseMacro("parens", LPAR, RPAR)
}

func Brackets() ref.T {
	return OpenCloseMacro("brackets", LBKT, RBKT)
}

func Curlies() ref.T {
	return OpenCloseMacro("curlies", LBRC, RBRC)
}

func DoubleQuotes() ref.T {
	return OpenCloseMacro("dquotes", DQT, DQT)
}

func SingleQuotes() ref.T {
	return OpenCloseMacro("squotes", SQT, SQT)
}

func BackQuotes() ref.T {
	return OpenCloseMacro("bquotes", GRAVE, GRAVE)
}

func OpenCloseMacro(name string, left, right key.T) ref.T {
	macro.Add(macro.T{
		Name:  name,
		Label: fmt.Sprintf("OpenClose %s", name),
		Cells: 0,
		Refs:  []ref.T{Kp(left), Kp(right), Kp(LEFT)},
		// Refs:  []ref.T{Kp(left), Kp(right), Kp(LEFT), To(PARENS)}, // todo
	})

	return ref0(name)
}

func XThenTrans(r ref.T, l layer.T, rc rowcol.T) ref.T {
	name := fmt.Sprintf("xThenTrans%d", l.Index())
	macro.Add(macro.T{
		Name:  name,
		Label: fmt.Sprintf("X Then Trans %s", l),
		Cells: 1,
		Refs:  []ref.T{macro.Param11, macro.Placeholder(r), To(l), l[rc]},
	})

	return ref.RefN(name, MapToAny(r.Args()))
}

func XThenLayer(r ref.T, l layer.T) ref.T {
	name := fmt.Sprintf("xThenLayer%d", l.Index())
	macro.Add(macro.T{
		Name:  name,
		Label: fmt.Sprintf("X Then Layer %s", l),
		Cells: 1,
		Refs:  []ref.T{macro.Param11, macro.Placeholder(r), To(l)},
	})

	return ref.RefN(name, MapToAny(r.Args()))
}

func TapNoRepeat(k key.T) ref.T {
	name := "TapNoRepeat"
	macro.Add(macro.T{
		Name:  name,
		Label: "Tap No Repeat",
		Cells: 1,
		Refs:  []ref.T{macro.Param11, macro.Placeholder(Kp(k)), macro.Pause},
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
			Label: fmt.Sprintf("To %s, %s", l.Name(), rc.Pretty()),
			Cells: 0,
			Refs:  []ref.T{To(l), macro.Press, l[rc], macro.Pause, macro.Release, l[rc]},
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
