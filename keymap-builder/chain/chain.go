package chain

import (
	"strings"

	"keyboard/instance"
	. "keyboard/key"
	"keyboard/layer"
	"keyboard/ref"
	"keyboard/rowcol"
	. "keyboard/util"
)

var chain = New()

type Chain struct {
	Keys   map[Key]ref.T
	Chains map[Key]*Chain
}

func New() Chain {
	return Chain{
		Keys:   make(map[Key]ref.T),
		Chains: make(map[Key]*Chain),
	}
}

func Add(keys string, r ref.T) {
	Panicif(len(keys) == 0)

	c := &chain

	keys = strings.ToUpper(keys)

	for i := range len(keys) - 1 {
		k := KeyFrom(keys[i])
		_, ok := c.Keys[k]
		Panicif(ok)
		newc, ok := c.Chains[k]
		if ok {
			c = newc
		} else {
			tmp := New()
			c.Chains[k] = &tmp
			c = &tmp
		}
	}

	k := KeyFrom(keys[len(keys)-1])
	_, ok := c.Chains[k]
	Panicif(ok)

	c.Keys[k] = r
}

func Compile(l layer.T, init func(layer.T)) {
	compile(chain, l.Name()+"_", l, init)
}

func compile(c Chain, prefix string, l layer.T, init func(layer.T)) {
	for k, ref := range c.Keys {
		l[rowcol.FromKey(k)] = ref
	}

	for k, sub := range SortedMapKV(c.Chains) {
		prefix := prefix + string(k)
		subl := layer.New(prefix, init)
		compile(*sub, prefix, subl, init)
		l[rowcol.FromKey(k)] = instance.To(subl)
	}
}
