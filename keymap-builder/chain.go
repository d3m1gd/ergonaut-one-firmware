package main

import (
	"fmt"
	"strings"

	. "keyboard/key"
	"keyboard/rowcol"
	. "keyboard/util"
)

var chain = NewChain()

type Chain struct {
	Keys   map[Key]Ref
	Chains map[Key]*Chain
}

func NewChain() Chain {
	return Chain{
		Keys:   make(map[Key]Ref),
		Chains: make(map[Key]*Chain),
	}
}

func AddChain(keys string, r Ref) {
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
			tmp := NewChain()
			c.Chains[k] = &tmp
			c = &tmp
		}
	}

	k := KeyFrom(keys[len(keys)-1])
	_, ok := c.Chains[k]
	Panicif(ok)

	c.Keys[k] = r
}

func CompileChains(layers []Layer) []Layer {
	layers = compileChain(chain, "", CHAINS, layers)
	return layers
}

func compileChain(c Chain, prefix string, n LayerIndex, layers []Layer) []Layer {
	for k, ref := range c.Keys {
		layers[n][rowcol.FromKey(k)] = ref
	}

	for k, sub := range SortedMapKV(c.Chains) {
		name := prefix + string(k)
		subn := NewLayerIndex(fmt.Sprintf("%s_%s", CHAINS, name))
		layers = append(layers, InitWith(Trans))
		layers = compileChain(*sub, name, subn, layers)
		layers[n][rowcol.FromKey(k)] = To(subn)
	}

	return layers
}
