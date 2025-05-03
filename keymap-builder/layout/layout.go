package layout

import (
	"slices"

	"keyboard/ref"
	"keyboard/rowcol"
)

var (
	ref0 = ref.Ref0
	ref1 = ref.Ref1
	ref2 = ref.Ref2

	L11 = rowcol.L(1, 1)
	L12 = rowcol.L(1, 2)
	L13 = rowcol.L(1, 3)
	L14 = rowcol.L(1, 4)
	L15 = rowcol.L(1, 5)
	L16 = rowcol.L(1, 6)
	L21 = rowcol.L(2, 1)
	L22 = rowcol.L(2, 2)
	L23 = rowcol.L(2, 3)
	L24 = rowcol.L(2, 4)
	L25 = rowcol.L(2, 5)
	L26 = rowcol.L(2, 6)
	L31 = rowcol.L(3, 1)
	L32 = rowcol.L(3, 2)
	L33 = rowcol.L(3, 3)
	L34 = rowcol.L(3, 4)
	L35 = rowcol.L(3, 5)
	L36 = rowcol.L(3, 6)
	L41 = rowcol.L(4, 1)
	L42 = rowcol.L(4, 2)
	L43 = rowcol.L(4, 3)
	R11 = rowcol.R(1, 1)
	R12 = rowcol.R(1, 2)
	R13 = rowcol.R(1, 3)
	R14 = rowcol.R(1, 4)
	R15 = rowcol.R(1, 5)
	R16 = rowcol.R(1, 6)
	R21 = rowcol.R(2, 1)
	R22 = rowcol.R(2, 2)
	R23 = rowcol.R(2, 3)
	R24 = rowcol.R(2, 4)
	R25 = rowcol.R(2, 5)
	R26 = rowcol.R(2, 6)
	R31 = rowcol.R(3, 1)
	R32 = rowcol.R(3, 2)
	R33 = rowcol.R(3, 3)
	R34 = rowcol.R(3, 4)
	R35 = rowcol.R(3, 5)
	R36 = rowcol.R(3, 6)
	R41 = rowcol.R(4, 1)
	R42 = rowcol.R(4, 2)
	R43 = rowcol.R(4, 3)

	L1 = []rowcol.T{L11, L12, L13, L14, L15, L16}
	L2 = []rowcol.T{L21, L22, L23, L24, L25, L26}
	L3 = []rowcol.T{L31, L32, L33, L34, L35, L36}
	L4 = []rowcol.T{L41, L42, L43}

	R1 = []rowcol.T{R11, R12, R13, R14, R15, R16}
	R2 = []rowcol.T{R21, R22, R23, R24, R25, R26}
	R3 = []rowcol.T{R31, R32, R33, R34, R35, R36}
	R4 = []rowcol.T{R41, R42, R43}

	LLL  = slices.Concat(L1, L2, L3)
	LLLL = slices.Concat(L1, L2, L3, L4)

	RRR  = slices.Concat(R1, R2, R3)
	RRRR = slices.Concat(R1, R2, R3, R4)
)
