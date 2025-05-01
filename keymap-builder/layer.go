package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/BooleanCat/go-functional/v2/it"
)

var LayerNames []string

type LayerIndexAuto int
type LayerIndex int

var BASE = NewLayerIndex("BASE")
var MOVER = NewLayerIndex("MOVER")
var NUMER = NewLayerIndex("NUMER")
var QUICK = NewLayerIndex("QUICK")
var REPEAT = NewLayerIndex("REPEAT")
var SYS = NewLayerIndex("SYS")
var PARENS = NewLayerIndex("PARENS")
var CHAINS = NewLayerIndex("CHAINS")

func NewLayerIndex(name string) LayerIndex {
	n := len(LayerNames)
	LayerNames = append(LayerNames, name)
	return LayerIndex(n)
}

func (l LayerIndex) String() string {
	return LayerNames[l]
}

type Layer map[RC]Ref

func NewLayer() Layer {
	return make(map[RC]Ref)
}

func (li LayerIndex) Render() string {
	return fmt.Sprintf("%d", li)
}

func (l Layer) Render() []string {
	widths := slices.Repeat([]int{0}, 12)
	cells := make([][]string, 4)
	for i := range cells {
		cells[i] = slices.Repeat([]string{""}, 12)
	}
	for rc := range RCs() {
		b := l[rc]
		rendered := CompileRef(b)
		row := rc.Row - 1
		col := rc.Col - 1
		if rc.Side == Right {
			col += 6
		} else if rc.Row == 4 {
			col += 3
		}
		cells[row][col] = rendered
		widths[col] = max(widths[col], len(rendered))
	}

	for _, row := range cells {
		for j := range row {
			row[j] = fmt.Sprintf("%-*s", widths[j], row[j])
		}
	}

	return slices.Collect(it.Map(slices.Values(cells), func(ss []string) string {
		s := strings.Join(ss, "  ")
		s = strings.TrimRight(s, " ")
		return s
	}))
}

type RenderedLayer struct {
	Index int
	Name  string
	Rows  []string
}

func InitWith(b Ref) Layer {
	return InitBy(func(RC) Ref { return b })
}

func InitBy(f func(RC) Ref) Layer {
	layer := NewLayer()
	for rc := range RCs() {
		layer[rc] = f(rc)
	}

	return layer
}

func InitToLevelTrans(index LayerIndex) Layer {
	base := layers[BASE]
	return InitBy(func(rc RC) Ref {
		name := fmt.Sprintf("to%d%s", index, rc)
		macro := Macro{
			Name:  name,
			Label: fmt.Sprintf("To %d, %s", index, rc.Pretty()),
			Cells: 0,
			Refs:  []Ref{To(index), MacroPress, base[rc], MacroPause, MacroRelease, base[rc]},
		}
		AddMacro(macro)
		return Ref0(name)
	})
}
