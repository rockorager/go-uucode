package uucode

import "unicode/utf8"

type CodePointIterator struct {
	I          int
	CodePoints []CodePoint
}

func NewCodePointIterator(codePoints []CodePoint) *CodePointIterator {
	return &CodePointIterator{CodePoints: codePoints}
}

func (it *CodePointIterator) Next() (CodePoint, bool) {
	if it.I >= len(it.CodePoints) {
		return 0, false
	}
	cp := it.CodePoints[it.I]
	it.I++
	return cp, true
}

func (it CodePointIterator) Peek() (CodePoint, bool) {
	return (&it).Next()
}

type UTF8Iterator struct {
	I     int
	Bytes []byte
}

func NewUTF8Iterator(s string) *UTF8Iterator {
	return &UTF8Iterator{Bytes: []byte(s)}
}

func NewUTF8BytesIterator(b []byte) *UTF8Iterator {
	return &UTF8Iterator{Bytes: b}
}

func (it *UTF8Iterator) Next() (CodePoint, bool) {
	if it.I >= len(it.Bytes) {
		return 0, false
	}
	r, size := utf8.DecodeRune(it.Bytes[it.I:])
	if r == utf8.RuneError && size == 0 {
		return 0, false
	}
	it.I += size
	return r, true
}

func (it UTF8Iterator) Peek() (CodePoint, bool) {
	return (&it).Next()
}

type RuneScanner interface {
	Next() (CodePoint, bool)
	Peek() (CodePoint, bool)
	Index() int
	Clone() RuneScanner
}

func (it *UTF8Iterator) Index() int      { return it.I }
func (it *CodePointIterator) Index() int { return it.I }

func (it *UTF8Iterator) Clone() RuneScanner {
	cp := *it
	return &cp
}

func (it *CodePointIterator) Clone() RuneScanner {
	cp := *it
	return &cp
}
