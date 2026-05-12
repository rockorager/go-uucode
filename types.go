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
