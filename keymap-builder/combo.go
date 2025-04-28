package main

import (
	"slices"
	"strings"
)

type Combo struct {
	Name               string
	Layers             []LayerIndex // array A list of layers on which the combo may be triggered. -1 allows all layers.
	Refs               []Ref        // phandle-array A behavior to run when the combo is triggered
	Keys               []RC         // array A list of key position indices for the keys which should trigger the combo
	TimoutMs           int          // int All the keys in key-positions must be pressed within this time in milliseconds to trigger the combo 50
	RequirePriorIdleMs int          // int If any non-modifier key is pressed within require-prior-idle-ms before a key in the combo, the key will not be considered for the combo -1 (disabled)
	SlowRelease        bool         // bool Releases the combo when all keys are released instead of when any key is released false
}

func (c Combo) RenderLayers() string {
	slices.Sort(c.Layers)
	return strings.Join(Map(c.Layers, func(x LayerIndex) string {
		return x.Render()
	}), " ")
}

func (c Combo) RenderBindings() string {
	return strings.Join(Map(c.Refs, CompileRef), " ")
}

func (c Combo) RenderKeys() string {
	return strings.Join(Map(c.Keys, func(x RC) string {
		return x.Render()
	}), " ")
}
