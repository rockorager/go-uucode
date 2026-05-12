package uucode

// CodePoint is a Unicode scalar value. Values outside Unicode's range are
// accepted by the API but resolve to default property values.
type CodePoint = rune

type GeneralCategory string

const (
	LetterUppercase         GeneralCategory = "letter_uppercase"
	LetterLowercase         GeneralCategory = "letter_lowercase"
	LetterTitlecase         GeneralCategory = "letter_titlecase"
	LetterModifier          GeneralCategory = "letter_modifier"
	LetterOther             GeneralCategory = "letter_other"
	MarkNonspacing          GeneralCategory = "mark_nonspacing"
	MarkSpacingCombining    GeneralCategory = "mark_spacing_combining"
	MarkEnclosing           GeneralCategory = "mark_enclosing"
	NumberDecimalDigit      GeneralCategory = "number_decimal_digit"
	NumberLetter            GeneralCategory = "number_letter"
	NumberOther             GeneralCategory = "number_other"
	PunctuationConnector    GeneralCategory = "punctuation_connector"
	PunctuationDash         GeneralCategory = "punctuation_dash"
	PunctuationOpen         GeneralCategory = "punctuation_open"
	PunctuationClose        GeneralCategory = "punctuation_close"
	PunctuationInitialQuote GeneralCategory = "punctuation_initial_quote"
	PunctuationFinalQuote   GeneralCategory = "punctuation_final_quote"
	PunctuationOther        GeneralCategory = "punctuation_other"
	SymbolMath              GeneralCategory = "symbol_math"
	SymbolCurrency          GeneralCategory = "symbol_currency"
	SymbolModifier          GeneralCategory = "symbol_modifier"
	SymbolOther             GeneralCategory = "symbol_other"
	SeparatorSpace          GeneralCategory = "separator_space"
	SeparatorLine           GeneralCategory = "separator_line"
	SeparatorParagraph      GeneralCategory = "separator_paragraph"
	OtherControl            GeneralCategory = "other_control"
	OtherFormat             GeneralCategory = "other_format"
	OtherSurrogate          GeneralCategory = "other_surrogate"
	OtherPrivateUse         GeneralCategory = "other_private_use"
	OtherNotAssigned        GeneralCategory = "other_not_assigned"
)

type BidiClass string

const (
	LeftToRight              BidiClass = "left_to_right"
	LeftToRightEmbedding     BidiClass = "left_to_right_embedding"
	LeftToRightOverride      BidiClass = "left_to_right_override"
	RightToLeft              BidiClass = "right_to_left"
	RightToLeftArabic        BidiClass = "right_to_left_arabic"
	RightToLeftEmbedding     BidiClass = "right_to_left_embedding"
	RightToLeftOverride      BidiClass = "right_to_left_override"
	PopDirectionalFormat     BidiClass = "pop_directional_format"
	EuropeanNumber           BidiClass = "european_number"
	EuropeanNumberSeparator  BidiClass = "european_number_separator"
	EuropeanNumberTerminator BidiClass = "european_number_terminator"
	ArabicNumber             BidiClass = "arabic_number"
	CommonNumberSeparator    BidiClass = "common_number_separator"
	NonspacingMark           BidiClass = "nonspacing_mark"
	BoundaryNeutral          BidiClass = "boundary_neutral"
	ParagraphSeparator       BidiClass = "paragraph_separator"
	SegmentSeparator         BidiClass = "segment_separator"
	Whitespace               BidiClass = "whitespace"
	OtherNeutrals            BidiClass = "other_neutrals"
	LeftToRightIsolate       BidiClass = "left_to_right_isolate"
	RightToLeftIsolate       BidiClass = "right_to_left_isolate"
	FirstStrongIsolate       BidiClass = "first_strong_isolate"
	PopDirectionalIsolate    BidiClass = "pop_directional_isolate"
)

type DecompositionType string

const (
	DecompositionDefault   DecompositionType = "default"
	DecompositionCanonical DecompositionType = "canonical"
	DecompositionFont      DecompositionType = "font"
	DecompositionNoBreak   DecompositionType = "noBreak"
	DecompositionInitial   DecompositionType = "initial"
	DecompositionMedial    DecompositionType = "medial"
	DecompositionFinal     DecompositionType = "final"
	DecompositionIsolated  DecompositionType = "isolated"
	DecompositionCircle    DecompositionType = "circle"
	DecompositionSuper     DecompositionType = "super"
	DecompositionSub       DecompositionType = "sub"
	DecompositionVertical  DecompositionType = "vertical"
	DecompositionWide      DecompositionType = "wide"
	DecompositionNarrow    DecompositionType = "narrow"
	DecompositionSmall     DecompositionType = "small"
	DecompositionSquare    DecompositionType = "square"
	DecompositionFraction  DecompositionType = "fraction"
	DecompositionCompat    DecompositionType = "compat"
)

type NumericType string

const (
	NumericNone    NumericType = "none"
	NumericDecimal NumericType = "decimal"
	NumericDigit   NumericType = "digit"
	NumericNumeric NumericType = "numeric"
)

type EastAsianWidth string

const (
	EastAsianNeutral   EastAsianWidth = "neutral"
	EastAsianFullwidth EastAsianWidth = "fullwidth"
	EastAsianHalfwidth EastAsianWidth = "halfwidth"
	EastAsianWide      EastAsianWidth = "wide"
	EastAsianNarrow    EastAsianWidth = "narrow"
	EastAsianAmbiguous EastAsianWidth = "ambiguous"
)

type GraphemeBreak string

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

type IndicConjunctBreak string

const (
	IndicConjunctNone      IndicConjunctBreak = "none"
	IndicConjunctLinker    IndicConjunctBreak = "linker"
	IndicConjunctConsonant IndicConjunctBreak = "consonant"
	IndicConjunctExtend    IndicConjunctBreak = "extend"
)

type BidiPairedBracketType string

const (
	BidiPairedBracketNone  BidiPairedBracketType = "none"
	BidiPairedBracketOpen  BidiPairedBracketType = "open"
	BidiPairedBracketClose BidiPairedBracketType = "close"
)

type BidiPairedBracket struct {
	Type BidiPairedBracketType
	Rune CodePoint
}

type Properties struct {
	Name                        string
	GeneralCategory             GeneralCategory
	CanonicalCombiningClass     uint8
	BidiClass                   BidiClass
	DecompositionType           DecompositionType
	DecompositionMapping        []CodePoint
	NumericType                 NumericType
	NumericValueDecimal         *uint8
	NumericValueDigit           *uint8
	NumericValueNumeric         string
	IsBidiMirrored              bool
	Unicode1Name                string
	SimpleUppercaseMapping      CodePoint
	SimpleLowercaseMapping      CodePoint
	SimpleTitlecaseMapping      CodePoint
	CaseFoldingSimple           CodePoint
	CaseFoldingFull             []CodePoint
	UppercaseMapping            []CodePoint
	LowercaseMapping            []CodePoint
	TitlecaseMapping            []CodePoint
	IsMath                      bool
	IsAlphabetic                bool
	IsLowercase                 bool
	IsUppercase                 bool
	IsCased                     bool
	IsCaseIgnorable             bool
	ChangesWhenLowercased       bool
	ChangesWhenUppercased       bool
	ChangesWhenTitlecased       bool
	ChangesWhenCasefolded       bool
	ChangesWhenCasemapped       bool
	IsIDStart                   bool
	IsIDContinue                bool
	IsXIDStart                  bool
	IsXIDContinue               bool
	IsDefaultIgnorable          bool
	IsGraphemeExtend            bool
	IsGraphemeBase              bool
	IsGraphemeLink              bool
	IndicConjunctBreak          IndicConjunctBreak
	EastAsianWidth              EastAsianWidth
	OriginalGraphemeBreak       GraphemeBreak
	IsEmoji                     bool
	IsEmojiPresentation         bool
	IsEmojiModifier             bool
	IsEmojiModifierBase         bool
	IsEmojiComponent            bool
	IsExtendedPictographic      bool
	IsEmojiVSBase               bool
	GraphemeBreak               GraphemeBreak
	BidiPairedBracket           BidiPairedBracket
	BidiMirroring               *CodePoint
	Block                       string
	Script                      string
	JoiningType                 string
	JoiningGroup                string
	IsCompositionExclusion      bool
	IndicPositionalCategory     string
	IndicSyllabicCategory       string
	GraphemeBreakNoControl      GraphemeBreak
	WcwidthStandalone           int
	WcwidthZeroInGrapheme       bool
	SpecialCasingConditions     []string
	SpecialLowercaseMapping     []CodePoint
	SpecialTitlecaseMapping     []CodePoint
	SpecialUppercaseMapping     []CodePoint
	SpecialLowercaseConditional []CodePoint
	SpecialTitlecaseConditional []CodePoint
	SpecialUppercaseConditional []CodePoint
}
