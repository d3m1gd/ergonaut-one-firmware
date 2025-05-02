package chain

import (
	"strings"

	"keyboard/instance"
	"keyboard/key"
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

	return name + "_" + keys
}

func Add(l layer.T, init func(layer.T), keyrefs map[string]ref.T) {
	for keys, r := range keyrefs {
		Panicif(len(keys) < 2)

		for _, p := range prefixes {
			Panicif(strings.HasPrefix(keys, p))
			Panicif(strings.HasPrefix(p, keys))
		}

		prefixes = append(prefixes, keys)

		keys = strings.ToUpper(keys)

		for i := range len(keys) - 1 {
			rc := rowcol.FromKey(key.From(keys[i]))
			name := Name(keys[:i+1])
			newl, ok := layer.Get(name)
			if !ok {
				newl = layer.New(name, init)
			}
			l[rc] = instance.To(newl)
			l = newl
		}

		rc := rowcol.FromKey(key.From(keys[len(keys)-1]))
		l[rc] = r
	}
}
