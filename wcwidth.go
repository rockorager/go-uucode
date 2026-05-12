package uucode

import "unicode/utf8"

// StringWidth returns the grapheme-aware terminal cell width of s.
//
// The result follows the width rules generated into this package's runtime
// tables and treats extended grapheme clusters as display units.
func StringWidth(s string) int {
	width := 0
	for i := 0; i < len(s); {
		if s[i] >= utf8.RuneSelf {
			return width + stringWidthFrom(s, i)
		}
		if i+1 < len(s) && s[i+1] >= utf8.RuneSelf {
			return width + stringWidthFrom(s, i)
		}
		if s[i] >= 0x20 && s[i] != 0x7f {
			width++
		}
		i++
	}
	return width
}

func stringWidthFrom(s string, start int) int {
	it := newUTF8WidthIterator(s, start)
	width := 0
	for it.hasNext {
		width += utf8WidthNext(&it)
	}
	return width
}

type utf8WidthIterator struct {
	s       string
	state   uint8
	nextEnd int
	nextCP  rune
	nextRow runtimeRow
	hasNext bool
}

func newUTF8WidthIterator(s string, start int) utf8WidthIterator {
	cp, size, ok := nextRuneInString(s, start)
	row := defaultRuntimeRow
	if ok {
		row = runtimeLookup(cp)
	}
	return utf8WidthIterator{
		s:       s,
		nextEnd: size,
		nextCP:  cp,
		nextRow: row,
		hasNext: ok,
	}
}

func (it *utf8WidthIterator) next() (rune, runtimeRow, bool, bool) {
	if !it.hasNext {
		return 0, defaultRuntimeRow, false, false
	}
	cp1 := it.nextCP
	row1 := it.nextRow
	cp2, nextEnd, ok := nextRuneInString(it.s, it.nextEnd)
	it.nextCP, it.nextEnd, it.hasNext = cp2, nextEnd, ok
	if ok {
		it.nextRow = runtimeLookup(cp2)
		return cp1, row1, computeGraphemeBreakRaw(row1.gb, it.nextRow.gb, &it.state), true
	}
	return cp1, row1, true, true
}

func utf8WidthNext(it *utf8WidthIterator) int {
	_, firstRow, firstBreak, ok := it.next()
	if !ok {
		return 0
	}
	standalone := firstRow.wcwidthStandalone()
	if firstBreak {
		return standalone
	}

	width := standalone
	if firstRow.wcwidthZeroInGrapheme() {
		width = 0
	}
	prevRow := firstRow
	prevState := it.state
	for {
		cp, row, isBreak, ok := it.next()
		if !ok {
			break
		}
		switch cp {
		case 0xfe0f:
			if prevRow.isEmojiVSBase() {
				width = 2
			}
		case 0xfe0e:
			if prevRow.isEmojiVSBase() {
				width = 1
			}
		case 0x200d:
			if prevState == breakStateExtendedPictographic && !isBreak {
				_, nextRow, nextBreak, ok := it.next()
				if !ok || nextBreak {
					return width
				}
				prevRow = nextRow
				prevState = it.state
				continue
			}
		case 0x1f3fb, 0x1f3fc, 0x1f3fd, 0x1f3fe, 0x1f3ff:
			width = 2
		default:
			if prevState == breakStateRegionalIndicator {
				width = 2
			} else if !row.wcwidthZeroInGrapheme() {
				width += row.wcwidthStandalone()
			}
		}
		if isBreak {
			break
		}
		prevRow = row
		prevState = it.state
	}
	return width
}
