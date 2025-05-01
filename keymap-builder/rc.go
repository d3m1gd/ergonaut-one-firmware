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
	row := n/12 + 1
	col := n%12 + 1
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

func ToSerial(rc RC) int {
	return rc.Offset() + 12*(rc.Row-1) + rc.Col - 1
}

func (rc RC) Serial() int {
	return rc.Offset() + 12*(rc.Row-1) + rc.Col - 1
}

func (rc RC) Less(other RC) int {
	return cmp.Compare(rc.Serial(), other.Serial())
}

func RCs() iter.Seq[RC] {
	return func(yield func(RC) bool) {
		for i := range 42 {
			if !yield(RCFrom(i)) {
				return
			}
		}
	}
}

func RCFromKeyCode(k KeyCode) RC {
	rc, ok := map[KeyCode]RC{
		TAB:   l(1, 1),
		Q:     l(1, 2),
		W:     l(1, 3),
		E:     l(1, 4),
		R:     l(1, 5),
		T:     l(1, 6),
		Y:     r(1, 1),
		U:     r(1, 2),
		I:     r(1, 3),
		O:     r(1, 4),
		P:     r(1, 5),
		RBRC:  r(1, 6),
		BSPC:  l(2, 1),
		A:     l(2, 2),
		S:     l(2, 3),
		D:     l(2, 4),
		F:     l(2, 5),
		G:     l(2, 6),
		H:     r(2, 1),
		J:     r(2, 2),
		K:     r(2, 3),
		L:     r(2, 4),
		COLON: r(2, 5),
		DQT:   r(2, 6),
		MINUS: l(3, 1),
		Z:     l(3, 2),
		X:     l(3, 3),
		C:     l(3, 4),
		V:     l(3, 5),
		B:     l(3, 6),
		N:     r(3, 1),
		M:     r(3, 2),
		LT:    r(3, 3),
		GT:    r(3, 4),
		QMARK: r(3, 5),
		BSLH:  r(3, 6),
		// row 4
		// C:     l(4, 1),
		// V:     l(4, 2),
		ESC:   l(4, 3),
		RET:   r(4, 1),
		SPACE: r(4, 2),
		// LT:    r(4, 3),
	}[k]
	panicif(!ok)

	return rc
}
