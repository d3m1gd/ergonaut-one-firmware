package rowcol

import (
	"cmp"
	"fmt"
	"iter"
	"strings"

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
	return rc.Serial()
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

func FromByte(b byte) RowCol {
	rc, ok := map[byte]RowCol{
		'T': L(1, 1),
		'q': L(1, 2),
		'w': L(1, 3),
		'e': L(1, 4),
		'r': L(1, 5),
		't': L(1, 6),
		'y': R(1, 1),
		'u': R(1, 2),
		'i': R(1, 3),
		'o': R(1, 4),
		'p': R(1, 5),
		'[': R(1, 6),
		// left
		'B': L(2, 1),
		'a': L(2, 2),
		's': L(2, 3),
		'd': L(2, 4),
		'f': L(2, 5),
		'g': L(2, 6),
		'h': R(2, 1),
		'j': R(2, 2),
		'k': R(2, 3),
		'l': R(2, 4),
		':': R(2, 5),
		'"': R(2, 6),
		'-': L(3, 1),
		'z': L(3, 2),
		'x': L(3, 3),
		'c': L(3, 4),
		'v': L(3, 5),
		'b': L(3, 6),
		'n': R(3, 1),
		'm': R(3, 2),
		',': R(3, 3),
		'.': R(3, 4),
		'/': R(3, 5),
		'|': R(3, 6),
		// row 4
		'<': L(4, 1),
		'Q': L(4, 2),
		'E': L(4, 3),
		'R': R(4, 1),
		'S': R(4, 2),
		' ': R(4, 2),
		'>': R(4, 3),

		// aliases
		'`': L(1, 1),
		'~': L(1, 1),
		'0': R(1, 1),
		'1': R(1, 2),
		'2': R(1, 3),
		'3': R(1, 4),
		'4': R(2, 2),
		'5': R(2, 3),
		'6': R(2, 4),
		'7': R(3, 2),
		'8': R(3, 3),
		'9': R(3, 4),
	}[b]
	Panicif(!ok, "key not handled: %s (%c)", b, b)

	return rc
}
