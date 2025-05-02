package layer

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strings"

	"keyboard/ref"
	"keyboard/rowcol"
	. "keyboard/util"
)

type T = Layer
type R = Rendered

type Layer map[rowcol.T]ref.T

type Rendered struct {
	Index int
	Name  string
	Rows  []string
}

var layers []Layer
var namerc = rowcol.T{}

func New(name string, init func(Layer)) Layer {
	layer := make(Layer)
	layer.Fill(init)
	layer[namerc] = ref.Ref0(name) // store name in special slot
	layers = append(layers, layer)
	return layer
}

func Get(name string) (Layer, bool) {
	for _, l := range layers {
		if l.Name() == name {
			return l, true
		}
	}

	return Layer{}, false
}

func Name(l Layer) string {
	return l.Name()
}

func (l Layer) Name() string {
	return l[namerc].Name
}

func (l Layer) String() string {
	return l.Name()
}

func (l Layer) Index() int {
	return slices.IndexFunc(layers, func(other Layer) bool {
		return l.Equal(other)
	})
}

func (l Layer) Equal(other Layer) bool {
	return l[namerc].Name == other[namerc].Name
}

func Less(l, other Layer) int {
	return l.Less(other)
}

func (l Layer) Less(other Layer) int {
	return cmp.Compare(l[namerc].Name, other[namerc].Name)
}

func (l Layer) Render() []string {
	widths := slices.Repeat([]int{0}, 12)
	cells := make([][]string, 4)
	for i := range cells {
		cells[i] = slices.Repeat([]string{""}, 12)
	}
	for rc := range rowcol.All() {
		b := l[rc]
		rendered := ref.Compile(b)
		row := rc.Row - 1
		col := rc.Col - 1
		if rc.Side == rowcol.Right {
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

	return Map(cells, func(ss []string) string {
		s := strings.Join(ss, "  ")
		s = strings.TrimRight(s, " ")
		return s
	})
}

func (l Layer) Extend(other Layer) {
	maps.Copy(l, other)
}

func (l Layer) Fill(fn func(Layer)) {
	fn(l)
}

func InitBy(f func(rowcol.T) ref.T) func(Layer) {
	return func(l Layer) {
		for rc := range rowcol.All() {
			l[rc] = f(rc)
		}
	}
}

func Render() []Rendered {
	return MapEnumerated(layers, func(n int, l Layer) Rendered {
		return Rendered{
			Index: n,
			Name:  l.Name(),
			Rows:  l.Render(),
		}
	})
}
