package uucode

import "unicode/utf8"

const (
	breakStateDefault uint8 = iota
	breakStateRegionalIndicator
	breakStateExtendedPictographic
	breakStateIndicConjunctConsonant
	breakStateIndicConjunctLinker
	breakStateCount
)

const (
	graphemeBreakPropertyCount = int(gbIndicConjunctConsonant + 1)
	graphemeBreakTableSize     = int(breakStateCount) * graphemeBreakPropertyCount * graphemeBreakPropertyCount
	graphemeBreakActionBreak   = 0x80
)

// Grapheme identifies a grapheme cluster by byte offsets into the original
// string.
type Grapheme struct {
	// Start is the byte offset of the first byte in the grapheme cluster.
	Start int
	// End is the byte offset just after the grapheme cluster.
	End int
}

// GraphemeIterator iterates over extended grapheme clusters in a string.
type GraphemeIterator struct {
	i       int
	state   uint8
	s       string
	nextEnd int
	nextCP  rune
	hasNext bool
	nextGB  uint8
	synced  bool
}

// NewGraphemeIterator returns a grapheme cluster iterator for s.
func NewGraphemeIterator(s string) GraphemeIterator {
	cp, size, gb, ok := nextGraphemeRuneInString(s, 0)
	return GraphemeIterator{
		s:       s,
		nextEnd: size,
		nextCP:  cp,
		hasNext: ok,
		nextGB:  gb,
		synced:  true,
	}
}

func nextRuneInString(s string, i int) (rune, int, bool) {
	if i >= len(s) {
		return 0, i, false
	}
	b0 := s[i]
	if b0 < utf8.RuneSelf {
		return rune(b0), i + 1, true
	}
	if b0 < 0xc2 {
		return utf8.RuneError, i + 1, true
	}
	if b0 < 0xe0 {
		if i+1 < len(s) {
			b1 := s[i+1]
			if b1&0xc0 == 0x80 {
				return rune(b0&0x1f)<<6 | rune(b1&0x3f), i + 2, true
			}
		}
		return utf8.RuneError, i + 1, true
	}
	if b0 < 0xf0 {
		if i+2 < len(s) {
			b1 := s[i+1]
			b2 := s[i+2]
			if b1&0xc0 == 0x80 && b2&0xc0 == 0x80 && (b0 != 0xe0 || b1 >= 0xa0) && (b0 != 0xed || b1 < 0xa0) {
				return rune(b0&0x0f)<<12 | rune(b1&0x3f)<<6 | rune(b2&0x3f), i + 3, true
			}
		}
		return utf8.RuneError, i + 1, true
	}
	if b0 < 0xf5 {
		if i+3 < len(s) {
			b1 := s[i+1]
			b2 := s[i+2]
			b3 := s[i+3]
			if b1&0xc0 == 0x80 && b2&0xc0 == 0x80 && b3&0xc0 == 0x80 && (b0 != 0xf0 || b1 >= 0x90) && (b0 != 0xf4 || b1 < 0x90) {
				return rune(b0&0x07)<<18 | rune(b1&0x3f)<<12 | rune(b2&0x3f)<<6 | rune(b3&0x3f), i + 4, true
			}
		}
		return utf8.RuneError, i + 1, true
	}
	return utf8.RuneError, i + 1, true
}

func nextGraphemeRuneInString(s string, i int) (rune, int, uint8, bool) {
	if i >= len(s) {
		return 0, i, gbOther, false
	}
	if s[i] < utf8.RuneSelf {
		return rune(s[i]), i + 1, asciiGraphemeBreak(s[i]), true
	}
	r, end, ok := nextRuneInString(s, i)
	if !ok {
		return 0, i, gbOther, false
	}
	return r, end, runtimeLookup(r).gb, true
}

func asciiGraphemeBreak(b byte) uint8 {
	switch b {
	case '\r':
		return gbCR
	case '\n':
		return gbLF
	}
	if b < 0x20 || b == 0x7f {
		return gbControl
	}
	return gbOther
}

// Next returns the next grapheme cluster.
//
// The returned Grapheme contains byte offsets into the original string. ok is
// false after the iterator is exhausted.
func (it *GraphemeIterator) Next() (Grapheme, bool) {
	if it.i >= len(it.s) {
		it.hasNext = false
		return Grapheme{}, false
	}
	start := it.i
	if b := it.s[it.i]; b >= 0x20 && b < 0x7f && (it.i+1 >= len(it.s) || it.s[it.i+1] < utf8.RuneSelf) {
		it.i++
		it.state = breakStateDefault
		it.synced = false
		return Grapheme{Start: start, End: it.i}, true
	}
	if !it.synced {
		cp, nextEnd, gb, ok := nextGraphemeRuneInString(it.s, it.i)
		it.nextCP, it.nextEnd, it.nextGB, it.hasNext = cp, nextEnd, gb, ok
		it.synced = true
	}
	if !it.hasNext {
		return Grapheme{}, false
	}
	for {
		gb1 := it.nextGB
		it.i = it.nextEnd
		cp, nextEnd, gb, ok := nextGraphemeRuneInString(it.s, it.nextEnd)
		it.nextCP, it.nextEnd, it.nextGB, it.hasNext = cp, nextEnd, gb, ok
		if !ok {
			return Grapheme{Start: start, End: it.i}, true
		}
		if computeGraphemeBreakRaw(gb1, it.nextGB, &it.state) {
			return Grapheme{Start: start, End: it.i}, true
		}
	}
}

// Peek returns the next grapheme cluster without advancing the iterator.
func (it GraphemeIterator) Peek() (Grapheme, bool) {
	return (&it).Next()
}

// IsBreak reports whether there is a grapheme cluster boundary between cp1 and
// cp2, updating state for rules that depend on previous code points.
func IsBreak(cp1, cp2 rune, state *BreakState) bool {
	return ComputeGraphemeBreak(runtimeLookup(cp1).graphemeBreak(), runtimeLookup(cp2).graphemeBreak(), state)
}

// ComputeGraphemeBreak reports whether there is a grapheme cluster boundary
// between two grapheme break properties, updating state for rules that depend
// on previous properties.
func ComputeGraphemeBreak(gb1, gb2 GraphemeBreak, state *BreakState) bool {
	rawState := breakStateDefault
	if state != nil {
		rawState = rawBreakState(*state)
	}
	ok := computeGraphemeBreakRaw(rawGraphemeBreak(gb1), rawGraphemeBreak(gb2), &rawState)
	if state != nil {
		*state = breakStateFromRaw(rawState)
	}
	return ok
}

func computeGraphemeBreakRaw(gb1, gb2 uint8, state *uint8) bool {
	if state == nil {
		s := breakStateDefault
		state = &s
	}
	if gb1 > gbIndicConjunctConsonant {
		gb1 = gbOther
	}
	if gb2 > gbIndicConjunctConsonant {
		gb2 = gbOther
	}
	if *state >= breakStateCount {
		*state = breakStateDefault
	}
	action := graphemeBreakTable[(uint32(*state)*uint32(graphemeBreakPropertyCount)+uint32(gb1))*uint32(graphemeBreakPropertyCount)+uint32(gb2)]
	*state = action &^ graphemeBreakActionBreak
	return action&graphemeBreakActionBreak != 0
}

func breakStateFromRaw(state uint8) BreakState {
	if state <= breakStateIndicConjunctLinker {
		return BreakState(state)
	}
	return BreakStateDefault
}

func rawBreakState(state BreakState) uint8 {
	if state <= BreakStateIndicConjunctLinker {
		return uint8(state)
	}
	return breakStateDefault
}

func rawGraphemeBreak(gb GraphemeBreak) uint8 {
	if gb <= GraphemeIndicConjunctConsonant {
		return uint8(gb)
	}
	return gbOther
}
