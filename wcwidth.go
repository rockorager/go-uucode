package uucode

func Wcwidth(it GraphemeIterator) int {
	return WcwidthNext(&it)
}

func WcwidthNext(it *GraphemeIterator) int {
	first, ok := it.NextCodePoint()
	if !ok {
		return 0
	}
	prevCP := first.CodePoint
	standalone := GetAll(prevCP).WcwidthStandalone
	if first.IsBreak {
		return standalone
	}

	width := standalone
	if GetAll(prevCP).WcwidthZeroInGrapheme {
		width = 0
	}
	prevState := it.State
	for {
		result, ok := it.NextCodePoint()
		if !ok {
			break
		}
		switch result.CodePoint {
		case 0xfe0f:
			if GetAll(prevCP).IsEmojiVSBase {
				width = 2
			}
		case 0xfe0e:
			if GetAll(prevCP).IsEmojiVSBase {
				width = 1
			}
		case 0x200d:
			if prevState == BreakStateExtendedPictographic && !result.IsBreak {
				next, ok := it.NextCodePoint()
				if !ok || next.IsBreak {
					return width
				}
				prevCP = next.CodePoint
				prevState = it.State
				continue
			}
		case 0x1f3fb, 0x1f3fc, 0x1f3fd, 0x1f3fe, 0x1f3ff:
			width = 2
		default:
			p := GetAll(result.CodePoint)
			if prevState == BreakStateRegionalIndicator {
				width = 2
			} else if !p.WcwidthZeroInGrapheme {
				width += p.WcwidthStandalone
			}
		}
		if result.IsBreak {
			break
		}
		prevCP = result.CodePoint
		prevState = it.State
	}
	return width
}

func WcwidthRemaining(it *GraphemeIterator) int {
	width := 0
	for it.hasNext {
		width += WcwidthNext(it)
	}
	return width
}

func UTF8Wcwidth(s string) int {
	return WcwidthRemaining(NewUTF8GraphemeIterator(s))
}

func deriveWcwidth(cp CodePoint, p Properties) (int, bool) {
	width := 1
	if p.GeneralCategory == OtherControl ||
		p.GeneralCategory == OtherSurrogate ||
		p.GeneralCategory == SeparatorLine ||
		p.GeneralCategory == SeparatorParagraph {
		width = 0
	} else if cp == 0x00ad {
		width = 1
	} else if p.IsDefaultIgnorable {
		width = 0
	} else if cp == 0x2e3a {
		width = 2
	} else if cp == 0x2e3b {
		width = 3
	} else if p.EastAsianWidth == EastAsianWide || p.EastAsianWidth == EastAsianFullwidth {
		width = 2
	} else if p.GraphemeBreak == GraphemeRegionalIndicator {
		width = 2
	}
	standalone := width
	if cp == 0x20e3 {
		standalone = 2
	}
	zeroInGrapheme := width == 0 ||
		p.IsEmojiModifier ||
		p.GeneralCategory == MarkNonspacing ||
		p.GeneralCategory == MarkEnclosing ||
		p.GraphemeBreak == GraphemeV ||
		p.GraphemeBreak == GraphemeT ||
		p.GraphemeBreak == GraphemePrepend
	return standalone, zeroInGrapheme
}
