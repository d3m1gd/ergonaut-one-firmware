package main

import (
	"cmp"
	"fmt"
	"iter"
	"strings"
)

type RC struct {
	Side Side
	Row  int
	Col  int
}

func RCFrom(n int) RC {
	row := (n-1)/12 + 1
	col := (n-1)%12 + 1
	side := Left
	if col > 6 {
		side = Right
		col -= 6
	}
	if row == 4 {
		if col > 3 {
			side = Right
			col -= 3
		}
	}

	return RC{side, row, col}
}

func r(a, b int) RC {
	return RC{Right, a, b}
}

func l(a, b int) RC {
	return RC{Left, a, b}
}

func (rc RC) Render() string {
	return fmt.Sprintf("%d", rc.Serial())
}

func (rc RC) String() string {
	return fmt.Sprintf("%s%d%d", rc.Side.Short(), rc.Row, rc.Col)
}

func (rc RC) Pretty() string {
	return fmt.Sprintf("%s%d%d", strings.ToUpper(rc.Side.Short()), rc.Row, rc.Col)
}

func (rc RC) Offset() int {
	offset := 0
	if rc.Side == Right {
		offset = 6
		if rc.Row == 4 {
			offset = 3
		}
	}

	return offset
}

func toSerial(rc RC) int {
	return rc.Offset() + 12*(rc.Row-1) + rc.Col
}

func (rc RC) Serial() int {
	return rc.Offset() + 12*(rc.Row-1) + rc.Col
}

func (rc RC) Less(other RC) int {
	return cmp.Compare(rc.Serial(), other.Serial())
}

func RCs() iter.Seq[RC] {
	return func(yield func(RC) bool) {
		for i := range 42 {
			if !yield(RCFrom(i + 1)) {
				return
			}
		}
	}
}
