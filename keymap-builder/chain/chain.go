package chain

import (
	"strings"

	"keyboard/instance"
	"keyboard/layer"
	"keyboard/ref"
	"keyboard/rowcol"
	. "keyboard/util"
)

var prefixes []string

const name = "CHAINS"

func Name(keys string) string {
	if len(keys) == 0 {
		return name
	}

	for _, k := range keys {
		// todo map to printable
		Panicif(k < 'a' || k > 'z')
	}

	return name + "_" + strings.ToUpper(keys)
}

func Add(start layer.T, init func(layer.T), keyrefs map[string]ref.T) {
	for keys, r := range SortedMap(keyrefs) {
		Panicif(len(keys) < 2)
		Panicif(strings.ToLower(keys) != keys)

		for _, p := range prefixes {
			Panicif(strings.HasPrefix(keys, p))
			Panicif(strings.HasPrefix(p, keys))
		}

		prefixes = append(prefixes, keys)

		l := start
		for i := range len(keys) - 1 {
			rc := rowcol.FromByte(keys[i])
			name := Name(keys[:i+1])
			newl, ok := layer.Get(name)
			if !ok {
				newl = layer.New(name, init)
			}
			l.Cells[rc] = instance.To(newl)
			l = newl
		}

		rc := rowcol.FromByte(keys[len(keys)-1])
		l.Cells[rc] = instance.OffX(l, r)
	}
}
