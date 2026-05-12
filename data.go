package uucode

import (
	"embed"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

//go:embed ucd/*.txt ucd/auxiliary/*.txt ucd/extracted/*.txt ucd/emoji/*.txt
var ucdFS embed.FS

const maxRune = 0x10ffff

type codeRange struct {
	start CodePoint
	end   CodePoint
}

type stringRange struct {
	codeRange
	value string
}

type ucdData struct {
	unicodeData map[CodePoint]*Properties
	bools       map[string][]codeRange
	strings     map[string][]stringRange
	caseFold    map[CodePoint][]CodePoint
	caseSimple  map[CodePoint]CodePoint
	bidiMirror  map[CodePoint]CodePoint
	bidiBracket map[CodePoint]BidiPairedBracket
}

var (
	loadOnce sync.Once
	loaded   *ucdData
	loadErr  error
)

func data() *ucdData {
	loadOnce.Do(func() {
		loaded, loadErr = loadUCD()
	})
	if loadErr != nil {
		panic(loadErr)
	}
	return loaded
}

func GetAll(cp CodePoint) Properties {
	d := data()
	p := defaultProperties(cp)
	if row := d.unicodeData[cp]; row != nil {
		mergeProperties(&p, row)
	}
	if v, ok := d.caseSimple[cp]; ok {
		p.CaseFoldingSimple = v
	}
	if v, ok := d.caseFold[cp]; ok {
		p.CaseFoldingFull = cloneRunes(v)
	}
	if v, ok := d.bidiMirror[cp]; ok {
		vv := v
		p.BidiMirroring = &vv
	}
	if v, ok := d.bidiBracket[cp]; ok {
		p.BidiPairedBracket = v
	}
	p.BidiClass = BidiClass(d.lookupString("bidi_class", cp, string(p.BidiClass)))
	p.EastAsianWidth = EastAsianWidth(d.lookupString("east_asian_width", cp, string(p.EastAsianWidth)))
	p.OriginalGraphemeBreak = GraphemeBreak(d.lookupString("original_grapheme_break", cp, string(p.OriginalGraphemeBreak)))
	p.Block = d.lookupString("block", cp, p.Block)
	p.Script = d.lookupString("script", cp, p.Script)
	p.JoiningType = d.lookupString("joining_type", cp, p.JoiningType)
	p.JoiningGroup = d.lookupString("joining_group", cp, p.JoiningGroup)
	p.IndicPositionalCategory = d.lookupString("indic_positional_category", cp, p.IndicPositionalCategory)
	p.IndicSyllabicCategory = d.lookupString("indic_syllabic_category", cp, p.IndicSyllabicCategory)
	p.IndicConjunctBreak = IndicConjunctBreak(d.lookupString("indic_conjunct_break", cp, string(p.IndicConjunctBreak)))

	for prop := range d.bools {
		if d.lookupBool(prop, cp) {
			setBoolProperty(&p, prop, true)
		}
	}

	p.GraphemeBreak = deriveGraphemeBreak(cp, p)
	p.GraphemeBreakNoControl = p.GraphemeBreak
	if p.GraphemeBreakNoControl == GraphemeControl || p.GraphemeBreakNoControl == GraphemeCR || p.GraphemeBreakNoControl == GraphemeLF {
		p.GraphemeBreakNoControl = GraphemeOther
	}
	p.WcwidthStandalone, p.WcwidthZeroInGrapheme = deriveWcwidth(cp, p)
	return p
}

func Get(field string, cp CodePoint) (any, bool) {
	p := GetAll(cp)
	switch field {
	case "name":
		return p.Name, true
	case "general_category":
		return p.GeneralCategory, true
	case "canonical_combining_class":
		return p.CanonicalCombiningClass, true
	case "bidi_class":
		return p.BidiClass, true
	case "decomposition_type":
		return p.DecompositionType, true
	case "decomposition_mapping":
		return cloneRunes(p.DecompositionMapping), true
	case "numeric_type":
		return p.NumericType, true
	case "numeric_value_decimal":
		return p.NumericValueDecimal, true
	case "numeric_value_digit":
		return p.NumericValueDigit, true
	case "numeric_value_numeric":
		return p.NumericValueNumeric, true
	case "is_bidi_mirrored":
		return p.IsBidiMirrored, true
	case "unicode_1_name":
		return p.Unicode1Name, true
	case "simple_uppercase_mapping":
		return p.SimpleUppercaseMapping, true
	case "simple_lowercase_mapping":
		return p.SimpleLowercaseMapping, true
	case "simple_titlecase_mapping":
		return p.SimpleTitlecaseMapping, true
	case "case_folding_simple":
		return p.CaseFoldingSimple, true
	case "case_folding_full":
		return cloneRunes(p.CaseFoldingFull), true
	case "lowercase_mapping":
		return cloneRunes(p.LowercaseMapping), true
	case "titlecase_mapping":
		return cloneRunes(p.TitlecaseMapping), true
	case "uppercase_mapping":
		return cloneRunes(p.UppercaseMapping), true
	case "east_asian_width":
		return p.EastAsianWidth, true
	case "original_grapheme_break":
		return p.OriginalGraphemeBreak, true
	case "grapheme_break":
		return p.GraphemeBreak, true
	case "bidi_paired_bracket":
		return p.BidiPairedBracket, true
	case "bidi_mirroring":
		return p.BidiMirroring, true
	case "block":
		return p.Block, true
	case "script":
		return p.Script, true
	case "joining_type":
		return p.JoiningType, true
	case "joining_group":
		return p.JoiningGroup, true
	case "indic_positional_category":
		return p.IndicPositionalCategory, true
	case "indic_syllabic_category":
		return p.IndicSyllabicCategory, true
	case "grapheme_break_no_control":
		return p.GraphemeBreakNoControl, true
	case "wcwidth_standalone":
		return p.WcwidthStandalone, true
	case "wcwidth_zero_in_grapheme":
		return p.WcwidthZeroInGrapheme, true
	default:
		return getBoolField(p, field)
	}
}

func Name(cp CodePoint) string                       { return GetAll(cp).Name }
func GeneralCategoryOf(cp CodePoint) GeneralCategory { return GetAll(cp).GeneralCategory }
func SimpleUppercaseMapping(cp CodePoint) CodePoint  { return GetAll(cp).SimpleUppercaseMapping }
func SimpleLowercaseMapping(cp CodePoint) CodePoint  { return GetAll(cp).SimpleLowercaseMapping }
func SimpleTitlecaseMapping(cp CodePoint) CodePoint  { return GetAll(cp).SimpleTitlecaseMapping }
func UppercaseMapping(cp CodePoint) []CodePoint      { return GetAll(cp).UppercaseMapping }
func LowercaseMapping(cp CodePoint) []CodePoint      { return GetAll(cp).LowercaseMapping }
func TitlecaseMapping(cp CodePoint) []CodePoint      { return GetAll(cp).TitlecaseMapping }
func CaseFoldingSimple(cp CodePoint) CodePoint       { return GetAll(cp).CaseFoldingSimple }
func CaseFoldingFull(cp CodePoint) []CodePoint       { return GetAll(cp).CaseFoldingFull }

func defaultProperties(cp CodePoint) Properties {
	p := Properties{
		GeneralCategory:         OtherNotAssigned,
		BidiClass:               LeftToRight,
		DecompositionType:       DecompositionDefault,
		DecompositionMapping:    []CodePoint{cp},
		NumericType:             NumericNone,
		SimpleUppercaseMapping:  cp,
		SimpleLowercaseMapping:  cp,
		SimpleTitlecaseMapping:  cp,
		CaseFoldingSimple:       cp,
		CaseFoldingFull:         []CodePoint{cp},
		UppercaseMapping:        []CodePoint{cp},
		LowercaseMapping:        []CodePoint{cp},
		TitlecaseMapping:        []CodePoint{cp},
		IndicConjunctBreak:      IndicConjunctNone,
		EastAsianWidth:          EastAsianNeutral,
		OriginalGraphemeBreak:   GraphemeOther,
		GraphemeBreak:           GraphemeOther,
		BidiPairedBracket:       BidiPairedBracket{Type: BidiPairedBracketNone},
		Block:                   "no_block",
		Script:                  "unknown",
		JoiningType:             "non_joining",
		JoiningGroup:            "no_joining_group",
		IndicPositionalCategory: "not_applicable",
		IndicSyllabicCategory:   "other",
	}
	return p
}

func mergeProperties(dst *Properties, src *Properties) {
	if src.Name != "" {
		dst.Name = src.Name
	}
	dst.GeneralCategory = src.GeneralCategory
	dst.CanonicalCombiningClass = src.CanonicalCombiningClass
	dst.BidiClass = src.BidiClass
	dst.DecompositionType = src.DecompositionType
	dst.DecompositionMapping = cloneRunes(src.DecompositionMapping)
	dst.NumericType = src.NumericType
	dst.NumericValueDecimal = cloneBytePtr(src.NumericValueDecimal)
	dst.NumericValueDigit = cloneBytePtr(src.NumericValueDigit)
	dst.NumericValueNumeric = src.NumericValueNumeric
	dst.IsBidiMirrored = src.IsBidiMirrored
	dst.Unicode1Name = src.Unicode1Name
	dst.SimpleUppercaseMapping = src.SimpleUppercaseMapping
	dst.SimpleLowercaseMapping = src.SimpleLowercaseMapping
	dst.SimpleTitlecaseMapping = src.SimpleTitlecaseMapping
	dst.UppercaseMapping = cloneRunes(src.UppercaseMapping)
	dst.LowercaseMapping = cloneRunes(src.LowercaseMapping)
	dst.TitlecaseMapping = cloneRunes(src.TitlecaseMapping)
}

func loadUCD() (*ucdData, error) {
	d := &ucdData{
		unicodeData: map[CodePoint]*Properties{},
		bools:       map[string][]codeRange{},
		strings:     map[string][]stringRange{},
		caseFold:    map[CodePoint][]CodePoint{},
		caseSimple:  map[CodePoint]CodePoint{},
		bidiMirror:  map[CodePoint]CodePoint{},
		bidiBracket: map[CodePoint]BidiPairedBracket{},
	}
	loaders := []func(*ucdData) error{
		loadUnicodeData,
		loadCaseFolding,
		loadSpecialCasing,
		loadDerivedCoreProperties,
		loadDerivedBidiClass,
		loadEastAsianWidth,
		loadGraphemeBreak,
		loadEmojiData,
		loadEmojiVariationSequences,
		loadBidiBrackets,
		loadBidiMirroring,
		loadBlocks,
		loadScripts,
		loadJoiningType,
		loadJoiningGroup,
		loadCompositionExclusions,
		loadIndicPositionalCategory,
		loadIndicSyllabicCategory,
	}
	for _, loader := range loaders {
		if err := loader(d); err != nil {
			return nil, err
		}
	}
	return d, nil
}

func loadUnicodeData(d *ucdData) error {
	return eachDataLine("ucd/UnicodeData.txt", func(line string) error {
		f := strings.Split(line, ";")
		if len(f) < 15 {
			return fmt.Errorf("bad UnicodeData line: %q", line)
		}
		cp, err := parseCP(f[0])
		if err != nil {
			return err
		}
		p := defaultProperties(cp)
		p.Name = f[1]
		p.GeneralCategory = generalCategory(f[2])
		ccc, _ := strconv.ParseUint(f[3], 10, 8)
		p.CanonicalCombiningClass = uint8(ccc)
		p.BidiClass = bidiClass(f[4])
		p.DecompositionType, p.DecompositionMapping = parseDecomposition(f[5], cp)
		p.NumericValueDecimal = parseBytePtr(f[6])
		p.NumericValueDigit = parseBytePtr(f[7])
		if f[8] != "" {
			p.NumericType = NumericNumeric
			p.NumericValueNumeric = f[8]
		} else if p.NumericValueDigit != nil {
			p.NumericType = NumericDigit
		} else if p.NumericValueDecimal != nil {
			p.NumericType = NumericDecimal
		}
		p.IsBidiMirrored = f[9] == "Y"
		p.Unicode1Name = f[10]
		p.SimpleUppercaseMapping = parseMappingCP(f[12], cp)
		p.SimpleLowercaseMapping = parseMappingCP(f[13], cp)
		p.SimpleTitlecaseMapping = parseMappingCP(f[14], cp)
		p.UppercaseMapping = []CodePoint{p.SimpleUppercaseMapping}
		p.LowercaseMapping = []CodePoint{p.SimpleLowercaseMapping}
		p.TitlecaseMapping = []CodePoint{p.SimpleTitlecaseMapping}
		d.unicodeData[cp] = &p
		return nil
	})
}

func loadCaseFolding(d *ucdData) error {
	return eachDataLine("ucd/CaseFolding.txt", func(line string) error {
		f := splitSemi(line)
		if len(f) < 3 {
			return nil
		}
		cp, err := parseCP(f[0])
		if err != nil {
			return err
		}
		mapping, err := parseCPList(f[2])
		if err != nil {
			return err
		}
		switch f[1] {
		case "C":
			d.caseSimple[cp] = mapping[0]
			d.caseFold[cp] = mapping
		case "S":
			d.caseSimple[cp] = mapping[0]
		case "F":
			d.caseFold[cp] = mapping
		}
		return nil
	})
}

func loadSpecialCasing(d *ucdData) error {
	return eachDataLine("ucd/SpecialCasing.txt", func(line string) error {
		f := splitSemi(line)
		if len(f) < 5 {
			return nil
		}
		cp, err := parseCP(f[0])
		if err != nil {
			return err
		}
		lower, err := parseCPList(f[1])
		if err != nil {
			return err
		}
		title, err := parseCPList(f[2])
		if err != nil {
			return err
		}
		upper, err := parseCPList(f[3])
		if err != nil {
			return err
		}
		row := d.unicodeData[cp]
		if row == nil {
			p := defaultProperties(cp)
			row = &p
			d.unicodeData[cp] = row
		}
		if f[4] == "" {
			row.LowercaseMapping = lower
			row.TitlecaseMapping = title
			row.UppercaseMapping = upper
			row.SpecialLowercaseMapping = lower
			row.SpecialTitlecaseMapping = title
			row.SpecialUppercaseMapping = upper
		} else {
			row.SpecialCasingConditions = strings.Fields(f[4])
			row.SpecialLowercaseConditional = lower
			row.SpecialTitlecaseConditional = title
			row.SpecialUppercaseConditional = upper
		}
		return nil
	})
}

func loadDerivedCoreProperties(d *ucdData) error {
	return eachDataLine("ucd/DerivedCoreProperties.txt", func(line string) error {
		f := splitSemi(line)
		if len(f) < 2 {
			return nil
		}
		r, err := parseRange(f[0])
		if err != nil {
			return err
		}
		if f[1] == "InCB" && len(f) >= 3 {
			d.addString("indic_conjunct_break", r, indicConjunctBreakName(f[2]))
			return nil
		}
		if prop := derivedCoreName(f[1]); prop != "" {
			d.addBool(prop, r)
		}
		return nil
	})
}

func loadDerivedBidiClass(d *ucdData) error {
	return eachDataLine("ucd/extracted/DerivedBidiClass.txt", func(line string) error {
		f := splitSemi(line)
		if len(f) < 2 {
			return nil
		}
		r, err := parseRange(f[0])
		if err != nil {
			return err
		}
		d.addString("bidi_class", r, string(bidiClass(f[1])))
		return nil
	})
}

func loadEastAsianWidth(d *ucdData) error {
	return loadRangeString(d, "ucd/extracted/DerivedEastAsianWidth.txt", "east_asian_width", eastAsianWidthName)
}

func loadGraphemeBreak(d *ucdData) error {
	return loadRangeString(d, "ucd/auxiliary/GraphemeBreakProperty.txt", "original_grapheme_break", graphemeBreakName)
}

func loadEmojiData(d *ucdData) error {
	return eachDataLine("ucd/emoji/emoji-data.txt", func(line string) error {
		f := splitSemi(line)
		if len(f) < 2 {
			return nil
		}
		r, err := parseRange(f[0])
		if err != nil {
			return err
		}
		if prop := emojiPropertyName(f[1]); prop != "" {
			d.addBool(prop, r)
		}
		return nil
	})
}

func loadEmojiVariationSequences(d *ucdData) error {
	return eachDataLine("ucd/emoji/emoji-variation-sequences.txt", func(line string) error {
		fields := strings.Fields(line)
		if len(fields) < 2 || fields[1] != "FE0E" {
			return nil
		}
		cp, err := parseCP(fields[0])
		if err != nil {
			return err
		}
		d.addBool("is_emoji_vs_base", codeRange{cp, cp})
		return nil
	})
}

func loadBidiBrackets(d *ucdData) error {
	return eachDataLine("ucd/BidiBrackets.txt", func(line string) error {
		f := splitSemi(line)
		if len(f) < 3 {
			return nil
		}
		cp, err := parseCP(f[0])
		if err != nil {
			return err
		}
		pair, err := parseCP(f[1])
		if err != nil {
			return err
		}
		typ := BidiPairedBracketNone
		if f[2] == "o" {
			typ = BidiPairedBracketOpen
		} else if f[2] == "c" {
			typ = BidiPairedBracketClose
		}
		d.bidiBracket[cp] = BidiPairedBracket{Type: typ, Rune: pair}
		return nil
	})
}

func loadBidiMirroring(d *ucdData) error {
	return eachDataLine("ucd/BidiMirroring.txt", func(line string) error {
		f := splitSemi(line)
		if len(f) < 2 {
			return nil
		}
		cp, err := parseCP(f[0])
		if err != nil {
			return err
		}
		mirror, err := parseCP(f[1])
		if err != nil {
			return err
		}
		d.bidiMirror[cp] = mirror
		return nil
	})
}

func loadBlocks(d *ucdData) error {
	return loadRangeString(d, "ucd/Blocks.txt", "block", normalizeName)
}

func loadScripts(d *ucdData) error {
	return loadRangeString(d, "ucd/Scripts.txt", "script", normalizeName)
}

func loadJoiningType(d *ucdData) error {
	return loadRangeString(d, "ucd/extracted/DerivedJoiningType.txt", "joining_type", joiningTypeName)
}

func loadJoiningGroup(d *ucdData) error {
	return loadRangeString(d, "ucd/extracted/DerivedJoiningGroup.txt", "joining_group", normalizeName)
}

func loadCompositionExclusions(d *ucdData) error {
	return eachDataLine("ucd/CompositionExclusions.txt", func(line string) error {
		r, err := parseRange(strings.Fields(line)[0])
		if err != nil {
			return err
		}
		d.addBool("is_composition_exclusion", r)
		return nil
	})
}

func loadIndicPositionalCategory(d *ucdData) error {
	return loadRangeString(d, "ucd/IndicPositionalCategory.txt", "indic_positional_category", normalizeName)
}

func loadIndicSyllabicCategory(d *ucdData) error {
	return loadRangeString(d, "ucd/IndicSyllabicCategory.txt", "indic_syllabic_category", normalizeName)
}

func loadRangeString(d *ucdData, path, prop string, conv func(string) string) error {
	return eachDataLine(path, func(line string) error {
		f := splitSemi(line)
		if len(f) < 2 {
			return nil
		}
		r, err := parseRange(f[0])
		if err != nil {
			return err
		}
		d.addString(prop, r, conv(f[1]))
		return nil
	})
}

func eachDataLine(path string, fn func(string) error) error {
	b, err := ucdFS.ReadFile(path)
	if err != nil {
		return err
	}
	for _, raw := range strings.Split(string(b), "\n") {
		line := trimComment(raw)
		if line == "" {
			continue
		}
		if err := fn(line); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
	}
	return nil
}

func trimComment(s string) string {
	if i := strings.IndexByte(s, '#'); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}

func splitSemi(s string) []string {
	parts := strings.Split(s, ";")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func parseRange(s string) (codeRange, error) {
	if i := strings.Index(s, ".."); i >= 0 {
		start, err := parseCP(s[:i])
		if err != nil {
			return codeRange{}, err
		}
		end, err := parseCP(s[i+2:])
		if err != nil {
			return codeRange{}, err
		}
		return codeRange{start: start, end: end}, nil
	}
	cp, err := parseCP(s)
	return codeRange{start: cp, end: cp}, err
}

func parseCP(s string) (CodePoint, error) {
	n, err := strconv.ParseInt(strings.TrimSpace(s), 16, 32)
	return CodePoint(n), err
}

func parseCPList(s string) ([]CodePoint, error) {
	fields := strings.Fields(s)
	out := make([]CodePoint, 0, len(fields))
	for _, f := range fields {
		cp, err := parseCP(f)
		if err != nil {
			return nil, err
		}
		out = append(out, cp)
	}
	return out, nil
}

func parseMappingCP(s string, same CodePoint) CodePoint {
	if strings.TrimSpace(s) == "" {
		return same
	}
	cp, err := parseCP(s)
	if err != nil {
		return same
	}
	return cp
}

func parseBytePtr(s string) *uint8 {
	if s == "" {
		return nil
	}
	n, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return nil
	}
	v := uint8(n)
	return &v
}

func parseDecomposition(s string, same CodePoint) (DecompositionType, []CodePoint) {
	if s == "" {
		return DecompositionDefault, []CodePoint{same}
	}
	typ := DecompositionCanonical
	if strings.HasPrefix(s, "<") {
		end := strings.IndexByte(s, '>')
		if end >= 0 {
			typ = decompositionType(s[1:end])
			s = strings.TrimSpace(s[end+1:])
		}
	}
	m, err := parseCPList(s)
	if err != nil || len(m) == 0 {
		return typ, []CodePoint{same}
	}
	return typ, m
}

func (d *ucdData) addBool(prop string, r codeRange) {
	d.bools[prop] = append(d.bools[prop], r)
}

func (d *ucdData) addString(prop string, r codeRange, value string) {
	d.strings[prop] = append(d.strings[prop], stringRange{codeRange: r, value: value})
}

func (d *ucdData) lookupBool(prop string, cp CodePoint) bool {
	for _, r := range d.bools[prop] {
		if cp >= r.start && cp <= r.end {
			return true
		}
	}
	return false
}

func (d *ucdData) lookupString(prop string, cp CodePoint, def string) string {
	for _, r := range d.strings[prop] {
		if cp >= r.start && cp <= r.end {
			return r.value
		}
	}
	return def
}

func cloneRunes(in []CodePoint) []CodePoint {
	if in == nil {
		return nil
	}
	out := make([]CodePoint, len(in))
	copy(out, in)
	return out
}

func cloneBytePtr(in *uint8) *uint8 {
	if in == nil {
		return nil
	}
	v := *in
	return &v
}

func normalizeName(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	return strings.ToLower(s)
}

func generalCategory(s string) GeneralCategory {
	switch s {
	case "Lu":
		return LetterUppercase
	case "Ll":
		return LetterLowercase
	case "Lt":
		return LetterTitlecase
	case "Lm":
		return LetterModifier
	case "Lo":
		return LetterOther
	case "Mn":
		return MarkNonspacing
	case "Mc":
		return MarkSpacingCombining
	case "Me":
		return MarkEnclosing
	case "Nd":
		return NumberDecimalDigit
	case "Nl":
		return NumberLetter
	case "No":
		return NumberOther
	case "Pc":
		return PunctuationConnector
	case "Pd":
		return PunctuationDash
	case "Ps":
		return PunctuationOpen
	case "Pe":
		return PunctuationClose
	case "Pi":
		return PunctuationInitialQuote
	case "Pf":
		return PunctuationFinalQuote
	case "Po":
		return PunctuationOther
	case "Sm":
		return SymbolMath
	case "Sc":
		return SymbolCurrency
	case "Sk":
		return SymbolModifier
	case "So":
		return SymbolOther
	case "Zs":
		return SeparatorSpace
	case "Zl":
		return SeparatorLine
	case "Zp":
		return SeparatorParagraph
	case "Cc":
		return OtherControl
	case "Cf":
		return OtherFormat
	case "Cs":
		return OtherSurrogate
	case "Co":
		return OtherPrivateUse
	default:
		return OtherNotAssigned
	}
}

func bidiClass(s string) BidiClass {
	switch s {
	case "L", "Left_To_Right":
		return LeftToRight
	case "LRE":
		return LeftToRightEmbedding
	case "LRO":
		return LeftToRightOverride
	case "R", "Right_To_Left":
		return RightToLeft
	case "AL", "Arabic_Letter":
		return RightToLeftArabic
	case "RLE":
		return RightToLeftEmbedding
	case "RLO":
		return RightToLeftOverride
	case "PDF":
		return PopDirectionalFormat
	case "EN":
		return EuropeanNumber
	case "ES":
		return EuropeanNumberSeparator
	case "ET", "European_Terminator":
		return EuropeanNumberTerminator
	case "AN":
		return ArabicNumber
	case "CS":
		return CommonNumberSeparator
	case "NSM":
		return NonspacingMark
	case "BN":
		return BoundaryNeutral
	case "B":
		return ParagraphSeparator
	case "S":
		return SegmentSeparator
	case "WS":
		return Whitespace
	case "ON":
		return OtherNeutrals
	case "LRI":
		return LeftToRightIsolate
	case "RLI":
		return RightToLeftIsolate
	case "FSI":
		return FirstStrongIsolate
	case "PDI":
		return PopDirectionalIsolate
	default:
		return LeftToRight
	}
}

func decompositionType(s string) DecompositionType {
	switch s {
	case "font":
		return DecompositionFont
	case "noBreak":
		return DecompositionNoBreak
	case "initial":
		return DecompositionInitial
	case "medial":
		return DecompositionMedial
	case "final":
		return DecompositionFinal
	case "isolated":
		return DecompositionIsolated
	case "circle":
		return DecompositionCircle
	case "super":
		return DecompositionSuper
	case "sub":
		return DecompositionSub
	case "vertical":
		return DecompositionVertical
	case "wide":
		return DecompositionWide
	case "narrow":
		return DecompositionNarrow
	case "small":
		return DecompositionSmall
	case "square":
		return DecompositionSquare
	case "fraction":
		return DecompositionFraction
	case "compat":
		return DecompositionCompat
	default:
		return DecompositionCanonical
	}
}

func derivedCoreName(s string) string {
	return map[string]string{
		"Math":                         "is_math",
		"Alphabetic":                   "is_alphabetic",
		"Lowercase":                    "is_lowercase",
		"Uppercase":                    "is_uppercase",
		"Cased":                        "is_cased",
		"Case_Ignorable":               "is_case_ignorable",
		"Changes_When_Lowercased":      "changes_when_lowercased",
		"Changes_When_Uppercased":      "changes_when_uppercased",
		"Changes_When_Titlecased":      "changes_when_titlecased",
		"Changes_When_Casefolded":      "changes_when_casefolded",
		"Changes_When_Casemapped":      "changes_when_casemapped",
		"ID_Start":                     "is_id_start",
		"ID_Continue":                  "is_id_continue",
		"XID_Start":                    "is_xid_start",
		"XID_Continue":                 "is_xid_continue",
		"Default_Ignorable_Code_Point": "is_default_ignorable",
		"Grapheme_Extend":              "is_grapheme_extend",
		"Grapheme_Base":                "is_grapheme_base",
		"Grapheme_Link":                "is_grapheme_link",
	}[s]
}

func eastAsianWidthName(s string) string {
	switch s {
	case "F":
		return string(EastAsianFullwidth)
	case "H":
		return string(EastAsianHalfwidth)
	case "W":
		return string(EastAsianWide)
	case "Na":
		return string(EastAsianNarrow)
	case "A":
		return string(EastAsianAmbiguous)
	default:
		return string(EastAsianNeutral)
	}
}

func graphemeBreakName(s string) string {
	return map[string]string{
		"Prepend":            string(GraphemePrepend),
		"CR":                 string(GraphemeCR),
		"LF":                 string(GraphemeLF),
		"Control":            string(GraphemeControl),
		"Extend":             "extend",
		"Regional_Indicator": string(GraphemeRegionalIndicator),
		"SpacingMark":        string(GraphemeSpacingMark),
		"L":                  string(GraphemeL),
		"V":                  string(GraphemeV),
		"T":                  string(GraphemeT),
		"LV":                 string(GraphemeLV),
		"LVT":                string(GraphemeLVT),
		"ZWJ":                string(GraphemeZWJ),
	}[s]
}

func emojiPropertyName(s string) string {
	return map[string]string{
		"Emoji":                 "is_emoji",
		"Emoji_Presentation":    "is_emoji_presentation",
		"Emoji_Modifier":        "is_emoji_modifier",
		"Emoji_Modifier_Base":   "is_emoji_modifier_base",
		"Emoji_Component":       "is_emoji_component",
		"Extended_Pictographic": "is_extended_pictographic",
	}[s]
}

func indicConjunctBreakName(s string) string {
	switch s {
	case "Linker":
		return string(IndicConjunctLinker)
	case "Consonant":
		return string(IndicConjunctConsonant)
	case "Extend":
		return string(IndicConjunctExtend)
	default:
		return string(IndicConjunctNone)
	}
}

func joiningTypeName(s string) string {
	return map[string]string{
		"R": "right_joining",
		"L": "left_joining",
		"D": "dual_joining",
		"C": "join_causing",
		"U": "non_joining",
		"T": "transparent",
	}[s]
}

func setBoolProperty(p *Properties, prop string, value bool) {
	switch prop {
	case "is_math":
		p.IsMath = value
	case "is_alphabetic":
		p.IsAlphabetic = value
	case "is_lowercase":
		p.IsLowercase = value
	case "is_uppercase":
		p.IsUppercase = value
	case "is_cased":
		p.IsCased = value
	case "is_case_ignorable":
		p.IsCaseIgnorable = value
	case "changes_when_lowercased":
		p.ChangesWhenLowercased = value
	case "changes_when_uppercased":
		p.ChangesWhenUppercased = value
	case "changes_when_titlecased":
		p.ChangesWhenTitlecased = value
	case "changes_when_casefolded":
		p.ChangesWhenCasefolded = value
	case "changes_when_casemapped":
		p.ChangesWhenCasemapped = value
	case "is_id_start":
		p.IsIDStart = value
	case "is_id_continue":
		p.IsIDContinue = value
	case "is_xid_start":
		p.IsXIDStart = value
	case "is_xid_continue":
		p.IsXIDContinue = value
	case "is_default_ignorable":
		p.IsDefaultIgnorable = value
	case "is_grapheme_extend":
		p.IsGraphemeExtend = value
	case "is_grapheme_base":
		p.IsGraphemeBase = value
	case "is_grapheme_link":
		p.IsGraphemeLink = value
	case "is_emoji":
		p.IsEmoji = value
	case "is_emoji_presentation":
		p.IsEmojiPresentation = value
	case "is_emoji_modifier":
		p.IsEmojiModifier = value
	case "is_emoji_modifier_base":
		p.IsEmojiModifierBase = value
	case "is_emoji_component":
		p.IsEmojiComponent = value
	case "is_extended_pictographic":
		p.IsExtendedPictographic = value
	case "is_emoji_vs_base":
		p.IsEmojiVSBase = value
	case "is_composition_exclusion":
		p.IsCompositionExclusion = value
	}
}

func getBoolField(p Properties, prop string) (any, bool) {
	switch prop {
	case "is_math":
		return p.IsMath, true
	case "is_alphabetic":
		return p.IsAlphabetic, true
	case "is_lowercase":
		return p.IsLowercase, true
	case "is_uppercase":
		return p.IsUppercase, true
	case "is_cased":
		return p.IsCased, true
	case "is_case_ignorable":
		return p.IsCaseIgnorable, true
	case "changes_when_lowercased":
		return p.ChangesWhenLowercased, true
	case "changes_when_uppercased":
		return p.ChangesWhenUppercased, true
	case "changes_when_titlecased":
		return p.ChangesWhenTitlecased, true
	case "changes_when_casefolded":
		return p.ChangesWhenCasefolded, true
	case "changes_when_casemapped":
		return p.ChangesWhenCasemapped, true
	case "is_id_start":
		return p.IsIDStart, true
	case "is_id_continue":
		return p.IsIDContinue, true
	case "is_xid_start":
		return p.IsXIDStart, true
	case "is_xid_continue":
		return p.IsXIDContinue, true
	case "is_default_ignorable":
		return p.IsDefaultIgnorable, true
	case "is_grapheme_extend":
		return p.IsGraphemeExtend, true
	case "is_grapheme_base":
		return p.IsGraphemeBase, true
	case "is_grapheme_link":
		return p.IsGraphemeLink, true
	case "is_emoji":
		return p.IsEmoji, true
	case "is_emoji_presentation":
		return p.IsEmojiPresentation, true
	case "is_emoji_modifier":
		return p.IsEmojiModifier, true
	case "is_emoji_modifier_base":
		return p.IsEmojiModifierBase, true
	case "is_emoji_component":
		return p.IsEmojiComponent, true
	case "is_extended_pictographic":
		return p.IsExtendedPictographic, true
	case "is_emoji_vs_base":
		return p.IsEmojiVSBase, true
	case "is_composition_exclusion":
		return p.IsCompositionExclusion, true
	default:
		return nil, false
	}
}
