package uucode

//go:generate go run ./cmd/uucodegen -out tables_gen.go

type runtimeRow struct {
	gb         uint8
	width      uint8
	wb         uint8
	sb         uint8
	lb         uint8
	eaw        uint8
	gc         uint8
	flags      uint8
	flags2     uint16
	upperDelta int32
	lowerDelta int32
	titleDelta int32
	foldDelta  int32
}

const maxRune = 0x10ffff

const (
	runtimeWidthMask      = 0x03
	runtimeZeroWidthFlag  = 0x04
	runtimeEmojiVSFlag    = 0x01
	runtimeEmojiPresFlag  = 0x02
	runtimeExtPictoFlag   = 0x04
	runtimeWhiteSpaceFlag = 0x08
)

const (
	runtimeASCIIHexDigitFlag uint16 = 1 << iota
	runtimeHexDigitFlag
	runtimeDashFlag
	runtimeDiacriticFlag
	runtimeQuotationMarkFlag
	runtimePatternSyntaxFlag
	runtimePatternWhiteSpaceFlag
	runtimeVariationSelectorFlag
	runtimeNoncharacterFlag
	runtimeUnifiedIdeographFlag
)

const (
	gbOther uint8 = iota
	gbControl
	gbPrepend
	gbCR
	gbLF
	gbRegionalIndicator
	gbSpacingMark
	gbL
	gbV
	gbT
	gbLV
	gbLVT
	gbZWJ
	gbZWNJ
	gbExtendedPictographic
	gbEmojiModifierBase
	gbEmojiModifier
	gbIndicConjunctExtend
	gbIndicConjunctLinker
	gbIndicConjunctConsonant
)

var defaultRuntimeRow = runtimeRow{gb: gbOther, width: 1}

func runtimeLookup(cp rune) runtimeRow {
	if cp < 0 || cp > maxRune {
		return defaultRuntimeRow
	}
	stage2Offset := runtimeStage1[cp>>8]
	rowIndex := runtimeStage2[stage2Offset+uint32(cp&0xff)]
	return runtimeStage3[rowIndex]
}

func (r runtimeRow) graphemeBreak() GraphemeBreak {
	return GraphemeBreak(r.gb)
}

func (r runtimeRow) wcwidthStandalone() int {
	return int(r.width & runtimeWidthMask)
}

func (r runtimeRow) wcwidthZeroInGrapheme() bool {
	return r.width&runtimeZeroWidthFlag != 0
}

func (r runtimeRow) isEmojiVSBase() bool {
	return r.flags&runtimeEmojiVSFlag != 0
}

func (r runtimeRow) isEmojiPresentation() bool {
	return r.flags&runtimeEmojiPresFlag != 0
}

func (r runtimeRow) isExtendedPictographic() bool {
	return r.flags&runtimeExtPictoFlag != 0
}

func (r runtimeRow) isWhiteSpace() bool {
	return r.flags&runtimeWhiteSpaceFlag != 0
}

func (r runtimeRow) isASCIIHexDigit() bool {
	return r.flags2&runtimeASCIIHexDigitFlag != 0
}

func (r runtimeRow) isHexDigit() bool {
	return r.flags2&runtimeHexDigitFlag != 0
}

func (r runtimeRow) isDash() bool {
	return r.flags2&runtimeDashFlag != 0
}

func (r runtimeRow) isDiacritic() bool {
	return r.flags2&runtimeDiacriticFlag != 0
}

func (r runtimeRow) isQuotationMark() bool {
	return r.flags2&runtimeQuotationMarkFlag != 0
}

func (r runtimeRow) isPatternSyntax() bool {
	return r.flags2&runtimePatternSyntaxFlag != 0
}

func (r runtimeRow) isPatternWhiteSpace() bool {
	return r.flags2&runtimePatternWhiteSpaceFlag != 0
}

func (r runtimeRow) isVariationSelector() bool {
	return r.flags2&runtimeVariationSelectorFlag != 0
}

func (r runtimeRow) isNoncharacter() bool {
	return r.flags2&runtimeNoncharacterFlag != 0
}

func (r runtimeRow) isUnifiedIdeograph() bool {
	return r.flags2&runtimeUnifiedIdeographFlag != 0
}

func (r runtimeRow) wordBreak() WordBreakClass {
	return WordBreakClass(r.wb)
}

func (r runtimeRow) sentenceBreak() SentenceBreakClass {
	return SentenceBreakClass(r.sb)
}

func (r runtimeRow) lineBreak() LineBreakClass {
	return LineBreakClass(r.lb)
}

func (r runtimeRow) eastAsianWidth() EastAsianWidthClass {
	return EastAsianWidthClass(r.eaw)
}

func (r runtimeRow) generalCategory() GeneralCategoryClass {
	return GeneralCategoryClass(r.gc)
}
