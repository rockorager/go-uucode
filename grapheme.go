package uucode

import "unicode/utf8"

// BreakState carries state between adjacent grapheme break decisions.
//
// Most callers should use GraphemeIterator instead of managing BreakState
// directly.
type BreakState string

// Grapheme break state values.
const (
	BreakStateDefault                BreakState = "default"
	BreakStateRegionalIndicator      BreakState = "regional_indicator"
	BreakStateExtendedPictographic   BreakState = "extended_pictographic"
	BreakStateIndicConjunctConsonant BreakState = "indic_conjunct_break_consonant"
	BreakStateIndicConjunctLinker    BreakState = "indic_conjunct_break_linker"
)

const (
	breakStateDefault uint8 = iota
	breakStateRegionalIndicator
	breakStateExtendedPictographic
	breakStateIndicConjunctConsonant
	breakStateIndicConjunctLinker
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
}

// NewGraphemeIterator returns a grapheme cluster iterator for s.
func NewGraphemeIterator(s string) GraphemeIterator {
	cp, size, ok := nextRuneInString(s, 0)
	gb := gbOther
	if ok {
		gb = runtimeLookup(cp).gb
	}
	return GraphemeIterator{
		s:       s,
		nextEnd: size,
		nextCP:  cp,
		hasNext: ok,
		nextGB:  gb,
	}
}

func nextRuneInString(s string, i int) (rune, int, bool) {
	if i >= len(s) {
		return 0, i, false
	}
	if s[i] < utf8.RuneSelf {
		return rune(s[i]), i + 1, true
	}
	r, size := utf8.DecodeRuneInString(s[i:])
	if r == utf8.RuneError && size == 0 {
		return 0, i, false
	}
	return r, i + size, true
}

// Next returns the next grapheme cluster.
//
// The returned Grapheme contains byte offsets into the original string. ok is
// false after the iterator is exhausted.
func (it *GraphemeIterator) Next() (Grapheme, bool) {
	if !it.hasNext {
		return Grapheme{}, false
	}
	start := it.i
	for {
		gb1 := it.nextGB
		it.i = it.nextEnd
		cp, nextEnd, ok := nextRuneInString(it.s, it.nextEnd)
		it.nextCP, it.nextEnd, it.hasNext = cp, nextEnd, ok
		if !ok {
			return Grapheme{Start: start, End: it.i}, true
		}
		it.nextGB = runtimeLookup(cp).gb
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
	switch *state {
	case breakStateRegionalIndicator:
		if gb1 != gbRegionalIndicator || gb2 != gbRegionalIndicator {
			*state = breakStateDefault
		}
	case breakStateExtendedPictographic:
		if !possiblyEmojiSequencePartRaw(gb1) || !possiblyEmojiSequencePartRaw(gb2) {
			*state = breakStateDefault
		}
	case breakStateIndicConjunctConsonant, breakStateIndicConjunctLinker:
		if !possiblyIndicSequencePartRaw(gb1) || !possiblyIndicSequencePartRaw(gb2) {
			*state = breakStateDefault
		}
	}

	if gb1 == gbCR && gb2 == gbLF {
		return false
	}
	if gb1 == gbControl || gb1 == gbCR || gb1 == gbLF {
		return true
	}
	if gb2 == gbControl || gb2 == gbCR || gb2 == gbLF {
		return true
	}
	if gb1 == gbL && (gb2 == gbL || gb2 == gbV || gb2 == gbLV || gb2 == gbLVT) {
		return false
	}
	if (gb1 == gbLV || gb1 == gbV) && (gb2 == gbV || gb2 == gbT) {
		return false
	}
	if (gb1 == gbLVT || gb1 == gbT) && gb2 == gbT {
		return false
	}
	if gb2 == gbSpacingMark {
		return false
	}
	if gb1 == gbPrepend {
		return false
	}

	if gb1 == gbIndicConjunctConsonant {
		if isIndicConjunctBreakExtendRaw(gb2) {
			*state = breakStateIndicConjunctConsonant
			return false
		}
		if gb2 == gbIndicConjunctLinker {
			*state = breakStateIndicConjunctLinker
			return false
		}
	} else if *state == breakStateIndicConjunctConsonant {
		if gb2 == gbIndicConjunctLinker {
			*state = breakStateIndicConjunctLinker
			return false
		}
		if isIndicConjunctBreakExtendRaw(gb2) {
			return false
		}
		*state = breakStateDefault
	} else if *state == breakStateIndicConjunctLinker {
		if gb2 == gbIndicConjunctLinker || isIndicConjunctBreakExtendRaw(gb2) {
			return false
		}
		if gb2 == gbIndicConjunctConsonant {
			*state = breakStateDefault
			return false
		}
		*state = breakStateDefault
	}

	if isExtendedPictographicRaw(gb1) {
		if isExtendRaw(gb2) || gb2 == gbZWJ {
			*state = breakStateExtendedPictographic
			return false
		}
		if gb1 == gbEmojiModifierBase && gb2 == gbEmojiModifier {
			*state = breakStateExtendedPictographic
			return false
		}
	} else if *state == breakStateExtendedPictographic {
		if (isExtendRaw(gb1) || gb1 == gbEmojiModifier) && (isExtendRaw(gb2) || gb2 == gbZWJ) {
			return false
		}
		if gb1 == gbZWJ && isExtendedPictographicRaw(gb2) {
			*state = breakStateDefault
			return false
		}
		*state = breakStateDefault
	}

	if gb1 == gbRegionalIndicator && gb2 == gbRegionalIndicator {
		if *state == breakStateDefault {
			*state = breakStateRegionalIndicator
			return false
		}
		*state = breakStateDefault
		return true
	}
	if isExtendRaw(gb2) || gb2 == gbZWJ {
		return false
	}
	return true
}

func breakStateFromRaw(state uint8) BreakState {
	switch state {
	case breakStateRegionalIndicator:
		return BreakStateRegionalIndicator
	case breakStateExtendedPictographic:
		return BreakStateExtendedPictographic
	case breakStateIndicConjunctConsonant:
		return BreakStateIndicConjunctConsonant
	case breakStateIndicConjunctLinker:
		return BreakStateIndicConjunctLinker
	default:
		return BreakStateDefault
	}
}

func rawBreakState(state BreakState) uint8 {
	switch state {
	case BreakStateRegionalIndicator:
		return breakStateRegionalIndicator
	case BreakStateExtendedPictographic:
		return breakStateExtendedPictographic
	case BreakStateIndicConjunctConsonant:
		return breakStateIndicConjunctConsonant
	case BreakStateIndicConjunctLinker:
		return breakStateIndicConjunctLinker
	default:
		return breakStateDefault
	}
}

func rawGraphemeBreak(gb GraphemeBreak) uint8 {
	switch gb {
	case GraphemeControl:
		return gbControl
	case GraphemePrepend:
		return gbPrepend
	case GraphemeCR:
		return gbCR
	case GraphemeLF:
		return gbLF
	case GraphemeRegionalIndicator:
		return gbRegionalIndicator
	case GraphemeSpacingMark:
		return gbSpacingMark
	case GraphemeL:
		return gbL
	case GraphemeV:
		return gbV
	case GraphemeT:
		return gbT
	case GraphemeLV:
		return gbLV
	case GraphemeLVT:
		return gbLVT
	case GraphemeZWJ:
		return gbZWJ
	case GraphemeZWNJ:
		return gbZWNJ
	case GraphemeExtendedPictographic:
		return gbExtendedPictographic
	case GraphemeEmojiModifierBase:
		return gbEmojiModifierBase
	case GraphemeEmojiModifier:
		return gbEmojiModifier
	case GraphemeIndicConjunctExtend:
		return gbIndicConjunctExtend
	case GraphemeIndicConjunctLinker:
		return gbIndicConjunctLinker
	case GraphemeIndicConjunctConsonant:
		return gbIndicConjunctConsonant
	default:
		return gbOther
	}
}

func isIndicConjunctBreakExtendRaw(gb uint8) bool {
	return gb == gbIndicConjunctExtend || gb == gbZWJ
}

func isExtendRaw(gb uint8) bool {
	return gb == gbZWNJ || gb == gbIndicConjunctExtend || gb == gbIndicConjunctLinker
}

func isExtendedPictographicRaw(gb uint8) bool {
	return gb == gbExtendedPictographic || gb == gbEmojiModifierBase
}

func possiblyEmojiSequencePartRaw(gb uint8) bool {
	return isExtendRaw(gb) || gb == gbZWJ || gb == gbExtendedPictographic || gb == gbEmojiModifierBase || gb == gbEmojiModifier
}

func possiblyIndicSequencePartRaw(gb uint8) bool {
	return gb == gbIndicConjunctConsonant || gb == gbIndicConjunctLinker || gb == gbIndicConjunctExtend || gb == gbZWJ
}
