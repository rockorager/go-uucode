package uucode

// GeneralCategory returns the Unicode general category for r.
func GeneralCategory(r rune) GeneralCategoryClass {
	return runtimeLookup(r).generalCategory()
}

// EastAsianWidth returns the Unicode East Asian Width property for r.
func EastAsianWidth(r rune) EastAsianWidthClass {
	return runtimeLookup(r).eastAsianWidth()
}

// WordBreak returns the Unicode word break property for r.
func WordBreak(r rune) WordBreakClass {
	return runtimeLookup(r).wordBreak()
}

// SentenceBreak returns the Unicode sentence break property for r.
func SentenceBreak(r rune) SentenceBreakClass {
	return runtimeLookup(r).sentenceBreak()
}

// LineBreak returns the Unicode line break property for r.
func LineBreak(r rune) LineBreakClass {
	return runtimeLookup(r).lineBreak()
}

// GraphemeBreakProperty returns the Unicode grapheme break property for r.
func GraphemeBreakProperty(r rune) GraphemeBreak {
	return runtimeLookup(r).graphemeBreak()
}

// RuneWidth returns the terminal cell width for r by itself.
func RuneWidth(r rune) int {
	return runtimeLookup(r).wcwidthStandalone()
}

// ToUpper maps r to its simple uppercase mapping.
func ToUpper(r rune) rune {
	return r + rune(runtimeLookup(r).upperDelta)
}

// ToLower maps r to its simple lowercase mapping.
func ToLower(r rune) rune {
	return r + rune(runtimeLookup(r).lowerDelta)
}

// ToTitle maps r to its simple titlecase mapping.
func ToTitle(r rune) rune {
	return r + rune(runtimeLookup(r).titleDelta)
}

// SimpleFold returns the next rune equivalent to r under simple case folding.
func SimpleFold(r rune) rune {
	return r + rune(runtimeLookup(r).foldDelta)
}

// IsUpper reports whether r has general category Lu.
func IsUpper(r rune) bool { return runtimeLookup(r).generalCategory() == GeneralCategoryLu }

// IsLower reports whether r has general category Ll.
func IsLower(r rune) bool { return runtimeLookup(r).generalCategory() == GeneralCategoryLl }

// IsTitle reports whether r has general category Lt.
func IsTitle(r rune) bool { return runtimeLookup(r).generalCategory() == GeneralCategoryLt }

// IsLetter reports whether r has a Unicode letter general category.
func IsLetter(r rune) bool {
	switch runtimeLookup(r).generalCategory() {
	case GeneralCategoryLu, GeneralCategoryLl, GeneralCategoryLt, GeneralCategoryLm, GeneralCategoryLo:
		return true
	default:
		return false
	}
}

// IsNumber reports whether r has a Unicode number general category.
func IsNumber(r rune) bool {
	switch runtimeLookup(r).generalCategory() {
	case GeneralCategoryNd, GeneralCategoryNl, GeneralCategoryNo:
		return true
	default:
		return false
	}
}

// IsDigit reports whether r has general category Nd.
func IsDigit(r rune) bool { return runtimeLookup(r).generalCategory() == GeneralCategoryNd }

// IsMark reports whether r has a Unicode mark general category.
func IsMark(r rune) bool {
	switch runtimeLookup(r).generalCategory() {
	case GeneralCategoryMn, GeneralCategoryMc, GeneralCategoryMe:
		return true
	default:
		return false
	}
}

// IsPunct reports whether r has a Unicode punctuation general category.
func IsPunct(r rune) bool {
	switch runtimeLookup(r).generalCategory() {
	case GeneralCategoryPc, GeneralCategoryPd, GeneralCategoryPs, GeneralCategoryPe, GeneralCategoryPi, GeneralCategoryPf, GeneralCategoryPo:
		return true
	default:
		return false
	}
}

// IsSymbol reports whether r has a Unicode symbol general category.
func IsSymbol(r rune) bool {
	switch runtimeLookup(r).generalCategory() {
	case GeneralCategorySm, GeneralCategorySc, GeneralCategorySk, GeneralCategorySo:
		return true
	default:
		return false
	}
}

// IsGraphic reports whether r is defined as a Graphic by Go's unicode package.
func IsGraphic(r rune) bool {
	switch runtimeLookup(r).generalCategory() {
	case GeneralCategoryLu, GeneralCategoryLl, GeneralCategoryLt, GeneralCategoryLm, GeneralCategoryLo,
		GeneralCategoryMn, GeneralCategoryMc, GeneralCategoryMe,
		GeneralCategoryNd, GeneralCategoryNl, GeneralCategoryNo,
		GeneralCategoryPc, GeneralCategoryPd, GeneralCategoryPs, GeneralCategoryPe, GeneralCategoryPi, GeneralCategoryPf, GeneralCategoryPo,
		GeneralCategorySm, GeneralCategorySc, GeneralCategorySk, GeneralCategorySo,
		GeneralCategoryZs:
		return true
	default:
		return false
	}
}

// IsPrint reports whether r is defined as printable by Go's unicode package.
func IsPrint(r rune) bool {
	if r == ' ' {
		return true
	}
	switch runtimeLookup(r).generalCategory() {
	case GeneralCategoryLu, GeneralCategoryLl, GeneralCategoryLt, GeneralCategoryLm, GeneralCategoryLo,
		GeneralCategoryMn, GeneralCategoryMc, GeneralCategoryMe,
		GeneralCategoryNd, GeneralCategoryNl, GeneralCategoryNo,
		GeneralCategoryPc, GeneralCategoryPd, GeneralCategoryPs, GeneralCategoryPe, GeneralCategoryPi, GeneralCategoryPf, GeneralCategoryPo,
		GeneralCategorySm, GeneralCategorySc, GeneralCategorySk, GeneralCategorySo:
		return true
	default:
		return false
	}
}

// IsControl reports whether r has general category Cc.
func IsControl(r rune) bool {
	return (r >= 0 && r <= 0x1f) || (r >= 0x7f && r <= 0x9f)
}

// IsSpace reports whether r has the Unicode White_Space property.
func IsSpace(r rune) bool { return runtimeLookup(r).isWhiteSpace() }

// IsASCIIHexDigit reports whether r has the Unicode ASCII_Hex_Digit property.
func IsASCIIHexDigit(r rune) bool { return runtimeLookup(r).isASCIIHexDigit() }

// IsHexDigit reports whether r has the Unicode Hex_Digit property.
func IsHexDigit(r rune) bool { return runtimeLookup(r).isHexDigit() }

// IsDash reports whether r has the Unicode Dash property.
func IsDash(r rune) bool { return runtimeLookup(r).isDash() }

// IsDiacritic reports whether r has the Unicode Diacritic property.
func IsDiacritic(r rune) bool { return runtimeLookup(r).isDiacritic() }

// IsQuotationMark reports whether r has the Unicode Quotation_Mark property.
func IsQuotationMark(r rune) bool { return runtimeLookup(r).isQuotationMark() }

// IsPatternSyntax reports whether r has the Unicode Pattern_Syntax property.
func IsPatternSyntax(r rune) bool { return runtimeLookup(r).isPatternSyntax() }

// IsPatternWhiteSpace reports whether r has the Unicode Pattern_White_Space property.
func IsPatternWhiteSpace(r rune) bool { return runtimeLookup(r).isPatternWhiteSpace() }

// IsVariationSelector reports whether r has the Unicode Variation_Selector property.
func IsVariationSelector(r rune) bool { return runtimeLookup(r).isVariationSelector() }

// IsNoncharacter reports whether r has the Unicode Noncharacter_Code_Point property.
func IsNoncharacter(r rune) bool { return runtimeLookup(r).isNoncharacter() }

// IsUnifiedIdeograph reports whether r has the Unicode Unified_Ideograph property.
func IsUnifiedIdeograph(r rune) bool { return runtimeLookup(r).isUnifiedIdeograph() }

// IsEmojiPresentation reports whether r has emoji presentation by default.
func IsEmojiPresentation(r rune) bool {
	return runtimeLookup(r).isEmojiPresentation()
}

// IsExtendedPictographic reports whether r has the Extended_Pictographic property.
func IsExtendedPictographic(r rune) bool {
	return runtimeLookup(r).isExtendedPictographic()
}
