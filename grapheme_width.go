package uucode

// GraphemeWithWidth identifies a grapheme cluster by byte offsets and includes
// its terminal cell width.
type GraphemeWithWidth struct {
	// Start is the byte offset of the first byte in the grapheme cluster.
	Start int
	// End is the byte offset just after the grapheme cluster.
	End int
	// Width is the terminal cell width of the grapheme cluster.
	Width int
}

// GraphemeWidthIterator iterates over extended grapheme clusters in a string
// and computes each cluster's terminal cell width in the same pass.
type GraphemeWidthIterator struct {
	i       int
	state   uint8
	s       string
	nextEnd int
	nextCP  rune
	nextRow runtimeTextRow
	hasNext bool
}

// NewGraphemeWidthIterator returns a grapheme cluster iterator for s that also
// computes display width for each cluster.
func NewGraphemeWidthIterator(s string) GraphemeWidthIterator {
	cp, size, ok := nextRuneInString(s, 0)
	row := defaultRuntimeTextRow
	if ok {
		if cp >= 0x20 && cp < 0x7f && (size >= len(s) || s[size] < 0x80) {
			row = defaultRuntimeTextRow
		} else {
			row = runtimeTextLookup(cp)
		}
	}
	return GraphemeWidthIterator{
		s:       s,
		nextEnd: size,
		nextCP:  cp,
		nextRow: row,
		hasNext: ok,
	}
}

// Next returns the next grapheme cluster with its terminal cell width.
//
// The returned GraphemeWithWidth contains byte offsets into the original string.
// ok is false after the iterator is exhausted.
func (it *GraphemeWidthIterator) Next() (GraphemeWithWidth, bool) {
	if !it.hasNext {
		return GraphemeWithWidth{}, false
	}
	start := it.i

	if it.nextCP >= 0x20 && it.nextCP < 0x7f && (it.nextEnd >= len(it.s) || it.s[it.nextEnd] < 0x80) {
		it.i = it.nextEnd
		if it.nextEnd >= len(it.s) {
			it.nextCP, it.hasNext = 0, false
		} else if it.s[it.nextEnd] >= 0x20 && it.s[it.nextEnd] < 0x7f {
			it.nextCP = rune(it.s[it.nextEnd])
			it.nextEnd++
			if it.nextEnd >= len(it.s) || it.s[it.nextEnd] < 0x80 {
				it.nextRow = defaultRuntimeTextRow
			} else {
				it.nextRow = runtimeTextLookup(it.nextCP)
			}
			it.hasNext = true
		} else {
			cp, nextEnd, ok := nextRuneInString(it.s, it.nextEnd)
			it.nextCP, it.nextEnd, it.hasNext = cp, nextEnd, ok
			if ok {
				it.nextRow = runtimeTextLookup(cp)
			}
		}
		it.state = breakStateDefault
		return GraphemeWithWidth{Start: start, End: it.i, Width: 1}, true
	}

	firstCP := it.nextCP
	firstRow := it.nextRow
	gb1 := firstRow.gb
	standalone := firstRow.wcwidthStandalone()

	it.i = it.nextEnd
	cp, nextEnd, ok := nextRuneInString(it.s, it.nextEnd)
	it.nextCP, it.nextEnd, it.hasNext = cp, nextEnd, ok
	if !ok {
		return GraphemeWithWidth{Start: start, End: it.i, Width: standalone}, true
	}
	it.nextRow = runtimeTextLookup(cp)
	isBreak := computeGraphemeBreakRaw(gb1, it.nextRow.gb, &it.state)
	if isBreak {
		return GraphemeWithWidth{Start: start, End: it.i, Width: standalone}, true
	}

	width := standalone
	if firstRow.wcwidthZeroInGrapheme() {
		width = 0
	}
	prevCP := firstCP
	prevRow := firstRow
	prevState := it.state
	for {
		cp := it.nextCP
		row := it.nextRow
		gb1 := row.gb
		it.i = it.nextEnd
		nextCP, nextEnd, ok := nextRuneInString(it.s, it.nextEnd)
		it.nextCP, it.nextEnd, it.hasNext = nextCP, nextEnd, ok
		if ok {
			it.nextRow = runtimeTextLookup(nextCP)
			isBreak = computeGraphemeBreakRaw(gb1, it.nextRow.gb, &it.state)
		} else {
			isBreak = true
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
				cp = it.nextCP
				row = it.nextRow
				gb1 := row.gb
				it.i = it.nextEnd
				nextCP, nextEnd, ok := nextRuneInString(it.s, it.nextEnd)
				it.nextCP, it.nextEnd, it.hasNext = nextCP, nextEnd, ok
				if ok {
					it.nextRow = runtimeTextLookup(nextCP)
					isBreak = computeGraphemeBreakRaw(gb1, it.nextRow.gb, &it.state)
				} else {
					isBreak = true
				}
				prevCP = cp
				prevRow = row
				prevState = it.state
				if isBreak {
					break
				}
				continue
			}
		case 0x1f3fb, 0x1f3fc, 0x1f3fd, 0x1f3fe, 0x1f3ff:
			width = 2
		default:
			if prevCP == 0x200d && prevState == breakStateExtendedPictographic && !isBreak {
				// The pictograph after an emoji ZWJ sequence is part of the existing
				// emoji cluster and should not add another cell width.
			} else if prevState == breakStateRegionalIndicator {
				width = 2
			} else if !row.wcwidthZeroInGrapheme() {
				width += row.wcwidthStandalone()
			}
		}
		if isBreak {
			break
		}
		prevCP = cp
		prevRow = row
		prevState = it.state
	}
	return GraphemeWithWidth{Start: start, End: it.i, Width: width}, true
}

// Peek returns the next grapheme cluster with width without advancing the iterator.
func (it GraphemeWidthIterator) Peek() (GraphemeWithWidth, bool) {
	return (&it).Next()
}
