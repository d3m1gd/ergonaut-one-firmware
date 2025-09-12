package combo

import (
	"slices"
	"strings"

	"keyboard/layer"
	"keyboard/ref"
	"keyboard/rowcol"
	. "keyboard/util"
	"keyboard/util/indenter"
)

type T = Combo

var combos []Combo

type Combo struct {
	Name   string
	Layers []layer.T  // array   A list of layers on which the combo may be triggered. -1 allows all layers.
	Ref    ref.T      // phandle-array A behavior to run when the combo is triggered
	Keys   []rowcol.T // array   A list of key position indices for the keys which should trigger the combo
	Timout int        // int ms  All the keys in key-positions must be pressed within this time in milliseconds to trigger the combo 50
	Idle   int        // int ms  If any non-modifier key is pressed within require-prior-idle-ms before a key in the combo, the key will not be considered for the combo -1 (disabled)
	Slow   bool       // bool    Releases the combo when all keys are released instead of when any key is released false
}

func (c Combo) RenderLayers() string {
	slices.SortFunc(c.Layers, layer.Less)
	return strings.Join(Map(c.Layers, layer.Name), " ")
}

func (c Combo) RenderKeys() string {
	return strings.Join(Map(c.Keys, func(x rowcol.T) string {
		return x.Render()
	}), " ")
}

func Add(c Combo) {
	combos = append(combos, c)
}

func Render() []Combo {
	return combos
}

func (c Combo) Compile(indent, level int) string {
	ir := indenter.New(indent)

	ir.Sprintf(0, "\n")
	ir.Sprintf(level, "%s {\n", c.Name)
	ir.Sprintf(level+1, "bindings = <%s>;\n", c.Ref.Compile())
	ir.Sprintf(level+1, "key-positions = <%s>;\n", c.RenderKeys())
	if len(c.Layers) > 0 {
		ir.Sprintf(level+1, "layers = <%s>\n", c.RenderLayers())
	}
	if c.Timout > 0 {
		ir.Sprintf(level+1, "timeout-ms = <%d>;\n", c.Timout)
	}
	if c.Idle > 0 {
		ir.Sprintf(level+1, "require-prior-idle-ms = <%d>;\n", c.Idle)
	}
	if c.Slow {
		ir.Sprintf(level+1, "slow-release;\n")
	}
	ir.Sprintf(level, "};\n")
	return ir.String()
}
