package uucode

// GraphemeBreak is a Unicode grapheme break property value.
type GraphemeBreak string

// Grapheme break property values.
const (
	GraphemeOther                  GraphemeBreak = "other"
	GraphemeControl                GraphemeBreak = "control"
	GraphemePrepend                GraphemeBreak = "prepend"
	GraphemeCR                     GraphemeBreak = "cr"
	GraphemeLF                     GraphemeBreak = "lf"
	GraphemeRegionalIndicator      GraphemeBreak = "regional_indicator"
	GraphemeSpacingMark            GraphemeBreak = "spacing_mark"
	GraphemeL                      GraphemeBreak = "l"
	GraphemeV                      GraphemeBreak = "v"
	GraphemeT                      GraphemeBreak = "t"
	GraphemeLV                     GraphemeBreak = "lv"
	GraphemeLVT                    GraphemeBreak = "lvt"
	GraphemeZWJ                    GraphemeBreak = "zwj"
	GraphemeZWNJ                   GraphemeBreak = "zwnj"
	GraphemeExtendedPictographic   GraphemeBreak = "extended_pictographic"
	GraphemeEmojiModifierBase      GraphemeBreak = "emoji_modifier_base"
	GraphemeEmojiModifier          GraphemeBreak = "emoji_modifier"
	GraphemeIndicConjunctExtend    GraphemeBreak = "indic_conjunct_break_extend"
	GraphemeIndicConjunctLinker    GraphemeBreak = "indic_conjunct_break_linker"
	GraphemeIndicConjunctConsonant GraphemeBreak = "indic_conjunct_break_consonant"
)

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
	return "XX"
}
