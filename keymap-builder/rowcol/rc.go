package rowcol

import (
	"cmp"
	"fmt"
	"iter"
	"strings"

	"keyboard/key"
	"keyboard/key/keys"
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

func FromKey(k key.T) RowCol {
	rc, ok := map[key.T]RowCol{
		keys.TAB:   L(1, 1),
		keys.Q:     L(1, 2),
		keys.W:     L(1, 3),
		keys.E:     L(1, 4),
		keys.R:     L(1, 5),
		keys.T:     L(1, 6),
		keys.Y:     R(1, 1),
		keys.U:     R(1, 2),
		keys.I:     R(1, 3),
		keys.O:     R(1, 4),
		keys.P:     R(1, 5),
		keys.RBRC:  R(1, 6),
		keys.BSPC:  L(2, 1),
		keys.A:     L(2, 2),
		keys.S:     L(2, 3),
		keys.D:     L(2, 4),
		keys.F:     L(2, 5),
		keys.G:     L(2, 6),
		keys.H:     R(2, 1),
		keys.J:     R(2, 2),
		keys.K:     R(2, 3),
		keys.L:     R(2, 4),
		keys.COLON: R(2, 5),
		keys.DQT:   R(2, 6),
		keys.MINUS: L(3, 1),
		keys.Z:     L(3, 2),
		keys.X:     L(3, 3),
		keys.C:     L(3, 4),
		keys.V:     L(3, 5),
		keys.B:     L(3, 6),
		keys.N:     R(3, 1),
		keys.M:     R(3, 2),
		keys.LT:    R(3, 3),
		keys.GT:    R(3, 4),
		keys.QMARK: R(3, 5),
		keys.BSLH:  R(3, 6),
		// row 4
		// C:     l(4, 1),
		// V:     l(4, 2),
		keys.ESC:   L(4, 3),
		keys.RET:   R(4, 1),
		keys.SPACE: R(4, 2),
		// LT:    r(4, 3),
	}[k]
	Panicif(!ok)

	return rc
}
