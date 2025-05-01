package rowcol

import (
	"cmp"
	"fmt"
	"iter"
	"strings"

	key "keyboard/key"
	. "keyboard/util"
)

type T = RowCol

type RowCol struct {
	Side Side
	Row  int
	Col  int
}

const Left Side = "left"
const Right Side = "right"

type Side string

func (s Side) Short() string {
	switch s {
	case Left:
		return "l"
	case Right:
		return "r"
	}
	panic("unhandled side: " + s)
}

func FromInt(n int) RowCol {
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

	return RowCol{side, row, col}
}

func R(a, b int) RowCol {
	return RowCol{Right, a, b}
}

func L(a, b int) RowCol {
	return RowCol{Left, a, b}
}

func (rc RowCol) Render() string {
	return fmt.Sprintf("%d", rc.Serial())
}

func (rc RowCol) String() string {
	return fmt.Sprintf("%s%d%d", rc.Side.Short(), rc.Row, rc.Col)
}

func (rc RowCol) Pretty() string {
	return fmt.Sprintf("%s%d%d", strings.ToUpper(rc.Side.Short()), rc.Row, rc.Col)
}

func (rc RowCol) Offset() int {
	offset := 0
	if rc.Side == Right {
		offset = 6
		if rc.Row == 4 {
			offset = 3
		}
	}

	return offset
}

func ToSerial(rc RowCol) int {
	return rc.Offset() + 12*(rc.Row-1) + rc.Col - 1
}

func (rc RowCol) Serial() int {
	return rc.Offset() + 12*(rc.Row-1) + rc.Col - 1
}

func (rc RowCol) Less(other RowCol) int {
	return cmp.Compare(rc.Serial(), other.Serial())
}

func All() iter.Seq[RowCol] {
	return func(yield func(RowCol) bool) {
		for i := range 42 {
			if !yield(FromInt(i)) {
				return
			}
		}
	}
}

func FromKey(k key.Key) RowCol {
	rc, ok := map[key.Key]RowCol{
		key.TAB:   L(1, 1),
		key.Q:     L(1, 2),
		key.W:     L(1, 3),
		key.E:     L(1, 4),
		key.R:     L(1, 5),
		key.T:     L(1, 6),
		key.Y:     R(1, 1),
		key.U:     R(1, 2),
		key.I:     R(1, 3),
		key.O:     R(1, 4),
		key.P:     R(1, 5),
		key.RBRC:  R(1, 6),
		key.BSPC:  L(2, 1),
		key.A:     L(2, 2),
		key.S:     L(2, 3),
		key.D:     L(2, 4),
		key.F:     L(2, 5),
		key.G:     L(2, 6),
		key.H:     R(2, 1),
		key.J:     R(2, 2),
		key.K:     R(2, 3),
		key.L:     R(2, 4),
		key.COLON: R(2, 5),
		key.DQT:   R(2, 6),
		key.MINUS: L(3, 1),
		key.Z:     L(3, 2),
		key.X:     L(3, 3),
		key.C:     L(3, 4),
		key.V:     L(3, 5),
		key.B:     L(3, 6),
		key.N:     R(3, 1),
		key.M:     R(3, 2),
		key.LT:    R(3, 3),
		key.GT:    R(3, 4),
		key.QMARK: R(3, 5),
		key.BSLH:  R(3, 6),
		// row 4
		// C:     l(4, 1),
		// V:     l(4, 2),
		key.ESC:   L(4, 3),
		key.RET:   R(4, 1),
		key.SPACE: R(4, 2),
		// LT:    r(4, 3),
	}[k]
	Panicif(!ok)

	return rc
}
