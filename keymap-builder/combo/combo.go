package combo

import (
	"slices"
	"strings"

	"keyboard/layer"
	"keyboard/ref"
	"keyboard/rowcol"
	. "keyboard/util"
)

type T = Combo

var combos []Combo

type Combo struct {
	Name        string
	Layers      []layer.T  // array A list of layers on which the combo may be triggered. -1 allows all layers.
	Refs        []ref.T    // phandle-array A behavior to run when the combo is triggered
	Keys        []rowcol.T // array A list of key position indices for the keys which should trigger the combo
	TimoutMs    int        // int All the keys in key-positions must be pressed within this time in milliseconds to trigger the combo 50
	IdleMs      int        // int If any non-modifier key is pressed within require-prior-idle-ms before a key in the combo, the key will not be considered for the combo -1 (disabled)
	SlowRelease bool       // bool Releases the combo when all keys are released instead of when any key is released false
}

func (c Combo) RenderLayers() string {
	slices.SortFunc(c.Layers, layer.Less)
	return strings.Join(Map(c.Layers, layer.Name), " ")
}

func (c Combo) RenderBindings() string {
	return strings.Join(Map(c.Refs, ref.CompileRef), " ")
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
