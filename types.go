package uucode

// GraphemeBreak is a Unicode grapheme break property value.
type GraphemeBreak uint8

// Grapheme break property values.
const (
	GraphemeOther GraphemeBreak = iota
	GraphemeControl
	GraphemePrepend
	GraphemeCR
	GraphemeLF
	GraphemeRegionalIndicator
	GraphemeSpacingMark
	GraphemeL
	GraphemeV
	GraphemeT
	GraphemeLV
	GraphemeLVT
	GraphemeZWJ
	GraphemeZWNJ
	GraphemeExtendedPictographic
	GraphemeEmojiModifierBase
	GraphemeEmojiModifier
	GraphemeIndicConjunctExtend
	GraphemeIndicConjunctLinker
	GraphemeIndicConjunctConsonant
)

var graphemeBreakNames = [...]string{
	"other",
	"control",
	"prepend",
	"cr",
	"lf",
	"regional_indicator",
	"spacing_mark",
	"l",
	"v",
	"t",
	"lv",
	"lvt",
	"zwj",
	"zwnj",
	"extended_pictographic",
	"emoji_modifier_base",
	"emoji_modifier",
	"indic_conjunct_break_extend",
	"indic_conjunct_break_linker",
	"indic_conjunct_break_consonant",
}

// String returns the Unicode grapheme break property name for gb.
func (gb GraphemeBreak) String() string {
	if int(gb) < len(graphemeBreakNames) {
		return graphemeBreakNames[gb]
	}
	return graphemeBreakNames[GraphemeOther]
}

// BreakState carries state between adjacent grapheme break decisions.
//
// Most callers should use GraphemeIterator instead of managing BreakState
// directly.
type BreakState uint8

// Grapheme break state values.
const (
	BreakStateDefault BreakState = iota
	BreakStateRegionalIndicator
	BreakStateExtendedPictographic
	BreakStateIndicConjunctConsonant
	BreakStateIndicConjunctLinker
)

var breakStateNames = [...]string{
	"default",
	"regional_indicator",
	"extended_pictographic",
	"indic_conjunct_break_consonant",
	"indic_conjunct_break_linker",
}

// String returns the state name for state.
func (state BreakState) String() string {
	if int(state) < len(breakStateNames) {
		return breakStateNames[state]
	}
	return breakStateNames[BreakStateDefault]
}

// WordBreakClass is a Unicode Word_Break property value.
type WordBreakClass uint8

// Unicode Word_Break property values.
const (
	WordBreakOther WordBreakClass = iota
	WordBreakLF
	WordBreakNewline
	WordBreakCR
	WordBreakWSegSpace
	WordBreakDoubleQuote
	WordBreakSingleQuote
	WordBreakMidNum
	WordBreakMidNumLet
	WordBreakNumeric
	WordBreakMidLetter
	WordBreakALetter
	WordBreakExtendNumLet
	WordBreakFormat
	WordBreakExtend
	WordBreakHebrewLetter
	WordBreakZWJ
	WordBreakKatakana
	WordBreakRegionalIndicator
)

// String returns the Unicode Word_Break property name for wb.
func (wb WordBreakClass) String() string {
	if int(wb) < len(runtimeWordBreakNames) {
		return runtimeWordBreakNames[wb]
	}
	return runtimeWordBreakNames[WordBreakOther]
}

// SentenceBreakClass is a Unicode Sentence_Break property value.
type SentenceBreakClass uint8

// Unicode Sentence_Break property values.
const (
	SentenceBreakOther SentenceBreakClass = iota
	SentenceBreakSp
	SentenceBreakLF
	SentenceBreakCR
	SentenceBreakSTerm
	SentenceBreakClose
	SentenceBreakSContinue
	SentenceBreakATerm
	SentenceBreakNumeric
	SentenceBreakUpper
	SentenceBreakLower
	SentenceBreakSep
	SentenceBreakFormat
	SentenceBreakOLetter
	SentenceBreakExtend
)

// String returns the Unicode Sentence_Break property name for sb.
func (sb SentenceBreakClass) String() string {
	if int(sb) < len(runtimeSentenceBreakNames) {
		return runtimeSentenceBreakNames[sb]
	}
	return runtimeSentenceBreakNames[SentenceBreakOther]
}

// LineBreakClass is a Unicode Line_Break property value.
type LineBreakClass uint8

// Unicode Line_Break property values.
const (
	LineBreakXX LineBreakClass = iota
	LineBreakCM
	LineBreakBA
	LineBreakLF
	LineBreakBK
	LineBreakCR
	LineBreakSP
	LineBreakEX
	LineBreakQU
	LineBreakAL
	LineBreakPR
	LineBreakPO
	LineBreakOP
	LineBreakCP
	LineBreakIS
	LineBreakHY
	LineBreakSY
	LineBreakNU
	LineBreakCL
	LineBreakNL
	LineBreakGL
	LineBreakAI
	LineBreakBB
	LineBreakHH
	LineBreakHL
	LineBreakSA
	LineBreakJL
	LineBreakJV
	LineBreakJT
	LineBreakNS
	LineBreakAK
	LineBreakVI
	LineBreakAS
	LineBreakID
	LineBreakVF
	LineBreakZW
	LineBreakZWJ
	LineBreakB2
	LineBreakIN
	LineBreakWJ
	LineBreakEB
	LineBreakCJ
	LineBreakH2
	LineBreakH3
	LineBreakSG
	LineBreakCB
	LineBreakAP
	LineBreakRI
	LineBreakEM
)

// String returns the Unicode Line_Break abbreviation for lb.
func (lb LineBreakClass) String() string {
	if int(lb) < len(runtimeLineBreakNames) {
		return runtimeLineBreakNames[lb]
	}
	return runtimeLineBreakNames[LineBreakXX]
}

// EastAsianWidthClass is a Unicode East_Asian_Width property value.
type EastAsianWidthClass uint8

// Unicode East_Asian_Width property values.
const (
	EastAsianWidthN EastAsianWidthClass = iota
	EastAsianWidthNa
	EastAsianWidthA
	EastAsianWidthW
	EastAsianWidthH
	EastAsianWidthF
)

// String returns the Unicode East_Asian_Width abbreviation for eaw.
func (eaw EastAsianWidthClass) String() string {
	if int(eaw) < len(runtimeEastAsianWidthNames) {
		return runtimeEastAsianWidthNames[eaw]
	}
	return runtimeEastAsianWidthNames[EastAsianWidthN]
}

// GeneralCategoryClass is a Unicode General_Category property value.
type GeneralCategoryClass uint8

// Unicode General_Category property values.
const (
	GeneralCategoryCn GeneralCategoryClass = iota
	GeneralCategoryCc
	GeneralCategoryZs
	GeneralCategoryPo
	GeneralCategorySc
	GeneralCategoryPs
	GeneralCategoryPe
	GeneralCategorySm
	GeneralCategoryPd
	GeneralCategoryNd
	GeneralCategoryLu
	GeneralCategorySk
	GeneralCategoryPc
	GeneralCategoryLl
	GeneralCategorySo
	GeneralCategoryLo
	GeneralCategoryPi
	GeneralCategoryCf
	GeneralCategoryNo
	GeneralCategoryPf
	GeneralCategoryLt
	GeneralCategoryLm
	GeneralCategoryMn
	GeneralCategoryMe
	GeneralCategoryMc
	GeneralCategoryNl
	GeneralCategoryZl
	GeneralCategoryZp
	GeneralCategoryCs
	GeneralCategoryCo
)

// String returns the Unicode General_Category abbreviation for gc.
func (gc GeneralCategoryClass) String() string {
	if int(gc) < len(runtimeGeneralCategoryNames) {
		return runtimeGeneralCategoryNames[gc]
	}
	return runtimeGeneralCategoryNames[GeneralCategoryCn]
}
