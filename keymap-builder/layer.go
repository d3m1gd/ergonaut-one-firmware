package main

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/BooleanCat/go-functional/v2/it"
)

type Layer map[RC]Reference
type LayerName string
type LayerIndex int

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
		rendered := CompileReference(b)
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

func (ln LayerName) Less(other LayerName) int {
	return cmp.Compare(ln, other)
}

type RenderedLayer struct {
	Index int
	Name  string
	Rows  []string
}

func InitWith(b Reference) Layer {
	return InitBy(func(RC) Reference { return b })
}

func InitBy(f func(RC) Reference) Layer {
	layer := Layer{}
	for rc := range RCs() {
		layer[rc] = f(rc)
	}

	return layer
}

func InitToLevelTrans(index LayerIndex) Layer {
	base := layers[BASE]
	return InitBy(func(rc RC) Reference {
		name := fmt.Sprintf("to%d%s", index, rc)
		macro := Macro{
			Name:  name,
			Label: fmt.Sprintf("To %d, %s", index, rc.Pretty()),
			Cells: 0,
			Refs:  []Reference{To{index}, MacroPress, base[rc], MacroWait, MacroRelease, base[rc]},
		}
		AddMacro(macro)
		return Custom0(name)
	})
}
