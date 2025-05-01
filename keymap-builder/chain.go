package main

import (
	"fmt"
	"strings"
)

var chain = NewChain()

type Chain struct {
	Keys   map[KeyCode]Ref
	Chains map[KeyCode]*Chain
}

func NewChain() Chain {
	return Chain{
		Keys:   make(map[KeyCode]Ref),
		Chains: make(map[KeyCode]*Chain),
	}
}

func AddChain(keys string, r Ref) {
	panicif(len(keys) == 0)

	c := &chain

	keys = strings.ToUpper(keys)

	for i := range len(keys) - 1 {
		k := KeyCodeFrom(keys[i])
		_, ok := c.Keys[k]
		panicif(ok)
		newc, ok := c.Chains[k]
		if ok {
			c = newc
		} else {
			tmp := NewChain()
			c.Chains[k] = &tmp
			c = &tmp
		}
	}

	k := KeyCodeFrom(keys[len(keys)-1])
	_, ok := c.Chains[k]
	panicif(ok)

	c.Keys[k] = r
}

func CompileChains(layers []Layer) []Layer {
	layers = compileChain(chain, "", CHAINS, layers)
	return layers
}

func compileChain(c Chain, prefix string, n LayerIndex, layers []Layer) []Layer {
	for k, ref := range c.Keys {
		layers[n][RCFromKeyCode(k)] = ref
	}

	for k, sub := range SortedMapKV(c.Chains) {
		name := prefix + string(k)
		subn := NewLayerIndex(fmt.Sprintf("%s_%s", CHAINS, name))
		layers = append(layers, InitWith(Trans))
		layers = compileChain(*sub, name, subn, layers)
		layers[n][RCFromKeyCode(k)] = To(subn)
	}

	return layers
}
