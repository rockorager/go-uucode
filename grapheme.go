package uucode

type BreakState string

const (
	BreakStateDefault                BreakState = "default"
	BreakStateRegionalIndicator      BreakState = "regional_indicator"
	BreakStateExtendedPictographic   BreakState = "extended_pictographic"
	BreakStateIndicConjunctConsonant BreakState = "indic_conjunct_break_consonant"
	BreakStateIndicConjunctLinker    BreakState = "indic_conjunct_break_linker"
)

type IteratorResult struct {
	CodePoint CodePoint
	IsBreak   bool
}

type Grapheme struct {
	Start int
	End   int
}

type indexedIterator interface {
	Next() (CodePoint, bool)
	Peek() (CodePoint, bool)
	Index() int
	Clone() RuneScanner
}

type GraphemeIterator struct {
	I       int
	State   BreakState
	cpIt    indexedIterator
	nextCP  CodePoint
	hasNext bool
	nextGB  GraphemeBreak
}

func NewGraphemeIterator(cpIt indexedIterator) *GraphemeIterator {
	cp, ok := cpIt.Next()
	gb := GraphemeOther
	if ok {
		gb = GetAll(cp).GraphemeBreak
	}
	return &GraphemeIterator{
		I:       0,
		State:   BreakStateDefault,
		cpIt:    cpIt,
		nextCP:  cp,
		hasNext: ok,
		nextGB:  gb,
	}
}

func NewUTF8GraphemeIterator(s string) *GraphemeIterator {
	return NewGraphemeIterator(NewUTF8Iterator(s))
}

func (it *GraphemeIterator) NextCodePoint() (IteratorResult, bool) {
	if !it.hasNext {
		return IteratorResult{}, false
	}
	cp1 := it.nextCP
	gb1 := it.nextGB
	it.I = it.cpIt.Index()
	cp2, ok := it.cpIt.Next()
	it.nextCP, it.hasNext = cp2, ok
	if ok {
		it.nextGB = GetAll(cp2).GraphemeBreak
		return IteratorResult{CodePoint: cp1, IsBreak: ComputeGraphemeBreak(gb1, it.nextGB, &it.State)}, true
	}
	return IteratorResult{CodePoint: cp1, IsBreak: true}, true
}

func (it GraphemeIterator) PeekCodePoint() (IteratorResult, bool) {
	it.cpIt = it.cpIt.Clone().(indexedIterator)
	return (&it).NextCodePoint()
}

func (it *GraphemeIterator) NextGrapheme() (Grapheme, bool) {
	if !it.hasNext {
		return Grapheme{}, false
	}
	start := it.I
	for {
		res, ok := it.NextCodePoint()
		if !ok {
			return Grapheme{}, false
		}
		if res.IsBreak {
			return Grapheme{Start: start, End: it.I}, true
		}
	}
}

func (it GraphemeIterator) PeekGrapheme() (Grapheme, bool) {
	it.cpIt = it.cpIt.Clone().(indexedIterator)
	return (&it).NextGrapheme()
}

func IsBreak(cp1, cp2 CodePoint, state *BreakState) bool {
	return ComputeGraphemeBreak(GetAll(cp1).GraphemeBreak, GetAll(cp2).GraphemeBreak, state)
}

func ComputeGraphemeBreak(gb1, gb2 GraphemeBreak, state *BreakState) bool {
	if state == nil {
		s := BreakStateDefault
		state = &s
	}
	switch *state {
	case BreakStateRegionalIndicator:
		if gb1 != GraphemeRegionalIndicator || gb2 != GraphemeRegionalIndicator {
			*state = BreakStateDefault
		}
	case BreakStateExtendedPictographic:
		if !possiblyEmojiSequencePart(gb1) || !possiblyEmojiSequencePart(gb2) {
			*state = BreakStateDefault
		}
	case BreakStateIndicConjunctConsonant, BreakStateIndicConjunctLinker:
		if !possiblyIndicSequencePart(gb1) || !possiblyIndicSequencePart(gb2) {
			*state = BreakStateDefault
		}
	}

	if gb1 == GraphemeCR && gb2 == GraphemeLF {
		return false
	}
	if gb1 == GraphemeControl || gb1 == GraphemeCR || gb1 == GraphemeLF {
		return true
	}
	if gb2 == GraphemeControl || gb2 == GraphemeCR || gb2 == GraphemeLF {
		return true
	}
	if gb1 == GraphemeL && (gb2 == GraphemeL || gb2 == GraphemeV || gb2 == GraphemeLV || gb2 == GraphemeLVT) {
		return false
	}
	if (gb1 == GraphemeLV || gb1 == GraphemeV) && (gb2 == GraphemeV || gb2 == GraphemeT) {
		return false
	}
	if (gb1 == GraphemeLVT || gb1 == GraphemeT) && gb2 == GraphemeT {
		return false
	}
	if gb2 == GraphemeSpacingMark {
		return false
	}
	if gb1 == GraphemePrepend {
		return false
	}

	if gb1 == GraphemeIndicConjunctConsonant {
		if isIndicConjunctBreakExtend(gb2) {
			*state = BreakStateIndicConjunctConsonant
			return false
		}
		if gb2 == GraphemeIndicConjunctLinker {
			*state = BreakStateIndicConjunctLinker
			return false
		}
	} else if *state == BreakStateIndicConjunctConsonant {
		if gb2 == GraphemeIndicConjunctLinker {
			*state = BreakStateIndicConjunctLinker
			return false
		}
		if isIndicConjunctBreakExtend(gb2) {
			return false
		}
		*state = BreakStateDefault
	} else if *state == BreakStateIndicConjunctLinker {
		if gb2 == GraphemeIndicConjunctLinker || isIndicConjunctBreakExtend(gb2) {
			return false
		}
		if gb2 == GraphemeIndicConjunctConsonant {
			*state = BreakStateDefault
			return false
		}
		*state = BreakStateDefault
	}

	if isExtendedPictographic(gb1) {
		if isExtend(gb2) || gb2 == GraphemeZWJ {
			*state = BreakStateExtendedPictographic
			return false
		}
		if gb1 == GraphemeEmojiModifierBase && gb2 == GraphemeEmojiModifier {
			*state = BreakStateExtendedPictographic
			return false
		}
	} else if *state == BreakStateExtendedPictographic {
		if (isExtend(gb1) || gb1 == GraphemeEmojiModifier) && (isExtend(gb2) || gb2 == GraphemeZWJ) {
			return false
		}
		if gb1 == GraphemeZWJ && isExtendedPictographic(gb2) {
			*state = BreakStateDefault
			return false
		}
		*state = BreakStateDefault
	}

	if gb1 == GraphemeRegionalIndicator && gb2 == GraphemeRegionalIndicator {
		if *state == BreakStateDefault {
			*state = BreakStateRegionalIndicator
			return false
		}
		*state = BreakStateDefault
		return true
	}
	if isExtend(gb2) || gb2 == GraphemeZWJ {
		return false
	}
	return true
}

func deriveGraphemeBreak(cp CodePoint, p Properties) GraphemeBreak {
	if p.IsEmojiModifier {
		return GraphemeEmojiModifier
	}
	if p.IsEmojiModifierBase {
		return GraphemeEmojiModifierBase
	}
	if p.IsExtendedPictographic {
		return GraphemeExtendedPictographic
	}
	switch p.IndicConjunctBreak {
	case IndicConjunctExtend:
		if cp == 0x200d {
			return GraphemeZWJ
		}
		return GraphemeIndicConjunctExtend
	case IndicConjunctLinker:
		return GraphemeIndicConjunctLinker
	case IndicConjunctConsonant:
		return GraphemeIndicConjunctConsonant
	}
	if p.OriginalGraphemeBreak == "extend" {
		if cp == 0x200c {
			return GraphemeZWNJ
		}
		return GraphemeIndicConjunctExtend
	}
	return p.OriginalGraphemeBreak
}

func isIndicConjunctBreakExtend(gb GraphemeBreak) bool {
	return gb == GraphemeIndicConjunctExtend || gb == GraphemeZWJ
}

func isExtend(gb GraphemeBreak) bool {
	return gb == GraphemeZWNJ || gb == GraphemeIndicConjunctExtend || gb == GraphemeIndicConjunctLinker
}

func isExtendedPictographic(gb GraphemeBreak) bool {
	return gb == GraphemeExtendedPictographic || gb == GraphemeEmojiModifierBase
}

func possiblyEmojiSequencePart(gb GraphemeBreak) bool {
	return isExtend(gb) || gb == GraphemeZWJ || gb == GraphemeExtendedPictographic || gb == GraphemeEmojiModifierBase || gb == GraphemeEmojiModifier
}

func possiblyIndicSequencePart(gb GraphemeBreak) bool {
	return gb == GraphemeIndicConjunctConsonant || gb == GraphemeIndicConjunctLinker || gb == GraphemeIndicConjunctExtend || gb == GraphemeZWJ
}
