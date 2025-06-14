package main

import (
	"testing"

	"keyboard/rowcol"

	"github.com/stretchr/testify/assert"
)

var l = rowcol.L
var r = rowcol.R

func TestRowColSerial(t *testing.T) {
	assert.Equal(t, 0, l(1, 1).Serial())
	assert.Equal(t, 1, l(1, 2).Serial())
	assert.Equal(t, 2, l(1, 3).Serial())
	assert.Equal(t, 3, l(1, 4).Serial())
	assert.Equal(t, 4, l(1, 5).Serial())
	assert.Equal(t, 5, l(1, 6).Serial())
	assert.Equal(t, 6, r(1, 1).Serial())
	assert.Equal(t, 7, r(1, 2).Serial())
	assert.Equal(t, 8, r(1, 3).Serial())
	assert.Equal(t, 9, r(1, 4).Serial())
	assert.Equal(t, 10, r(1, 5).Serial())
	assert.Equal(t, 11, r(1, 6).Serial())
	assert.Equal(t, 12, l(2, 1).Serial())
	assert.Equal(t, 13, l(2, 2).Serial())
	assert.Equal(t, 14, l(2, 3).Serial())
	assert.Equal(t, 15, l(2, 4).Serial())
	assert.Equal(t, 16, l(2, 5).Serial())
	assert.Equal(t, 17, l(2, 6).Serial())
	assert.Equal(t, 18, r(2, 1).Serial())
	assert.Equal(t, 19, r(2, 2).Serial())
	assert.Equal(t, 20, r(2, 3).Serial())
	assert.Equal(t, 21, r(2, 4).Serial())
	assert.Equal(t, 22, r(2, 5).Serial())
	assert.Equal(t, 23, r(2, 6).Serial())
	assert.Equal(t, 24, l(3, 1).Serial())
	assert.Equal(t, 25, l(3, 2).Serial())
	assert.Equal(t, 26, l(3, 3).Serial())
	assert.Equal(t, 27, l(3, 4).Serial())
	assert.Equal(t, 28, l(3, 5).Serial())
	assert.Equal(t, 29, l(3, 6).Serial())
	assert.Equal(t, 30, r(3, 1).Serial())
	assert.Equal(t, 31, r(3, 2).Serial())
	assert.Equal(t, 32, r(3, 3).Serial())
	assert.Equal(t, 33, r(3, 4).Serial())
	assert.Equal(t, 34, r(3, 5).Serial())
	assert.Equal(t, 35, r(3, 6).Serial())
	assert.Equal(t, 36, l(4, 1).Serial())
	assert.Equal(t, 37, l(4, 2).Serial())
	assert.Equal(t, 38, l(4, 3).Serial())
	assert.Equal(t, 39, r(4, 1).Serial())
	assert.Equal(t, 40, r(4, 2).Serial())
	assert.Equal(t, 41, r(4, 3).Serial())
}

func TestSerialRowCol(t *testing.T) {
	assert.Equal(t, l(1, 1), rowcol.FromInt(0))
	assert.Equal(t, l(1, 2), rowcol.FromInt(1))
	assert.Equal(t, l(1, 3), rowcol.FromInt(2))
	assert.Equal(t, l(1, 4), rowcol.FromInt(3))
	assert.Equal(t, l(1, 5), rowcol.FromInt(4))
	assert.Equal(t, l(1, 6), rowcol.FromInt(5))
	assert.Equal(t, r(1, 1), rowcol.FromInt(6))
	assert.Equal(t, r(1, 2), rowcol.FromInt(7))
	assert.Equal(t, r(1, 3), rowcol.FromInt(8))
	assert.Equal(t, r(1, 4), rowcol.FromInt(9))
	assert.Equal(t, r(1, 5), rowcol.FromInt(10))
	assert.Equal(t, r(1, 6), rowcol.FromInt(11))
	assert.Equal(t, l(2, 1), rowcol.FromInt(12))
	assert.Equal(t, l(2, 2), rowcol.FromInt(13))
	assert.Equal(t, l(2, 3), rowcol.FromInt(14))
	assert.Equal(t, l(2, 4), rowcol.FromInt(15))
	assert.Equal(t, l(2, 5), rowcol.FromInt(16))
	assert.Equal(t, l(2, 6), rowcol.FromInt(17))
	assert.Equal(t, r(2, 1), rowcol.FromInt(18))
	assert.Equal(t, r(2, 2), rowcol.FromInt(19))
	assert.Equal(t, r(2, 3), rowcol.FromInt(20))
	assert.Equal(t, r(2, 4), rowcol.FromInt(21))
	assert.Equal(t, r(2, 5), rowcol.FromInt(22))
	assert.Equal(t, r(2, 6), rowcol.FromInt(23))
	assert.Equal(t, l(3, 1), rowcol.FromInt(24))
	assert.Equal(t, l(3, 2), rowcol.FromInt(25))
	assert.Equal(t, l(3, 3), rowcol.FromInt(26))
	assert.Equal(t, l(3, 4), rowcol.FromInt(27))
	assert.Equal(t, l(3, 5), rowcol.FromInt(28))
	assert.Equal(t, l(3, 6), rowcol.FromInt(29))
	assert.Equal(t, r(3, 1), rowcol.FromInt(30))
	assert.Equal(t, r(3, 2), rowcol.FromInt(31))
	assert.Equal(t, r(3, 3), rowcol.FromInt(32))
	assert.Equal(t, r(3, 4), rowcol.FromInt(33))
	assert.Equal(t, r(3, 5), rowcol.FromInt(34))
	assert.Equal(t, r(3, 6), rowcol.FromInt(35))
	assert.Equal(t, l(4, 1), rowcol.FromInt(36))
	assert.Equal(t, l(4, 2), rowcol.FromInt(37))
	assert.Equal(t, l(4, 3), rowcol.FromInt(38))
	assert.Equal(t, r(4, 1), rowcol.FromInt(39))
	assert.Equal(t, r(4, 2), rowcol.FromInt(40))
	assert.Equal(t, r(4, 3), rowcol.FromInt(41))
}
