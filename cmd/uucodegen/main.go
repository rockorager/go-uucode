package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	maxRune   = 0x10ffff
	blockSize = 256
	numBlocks = (maxRune + 1) / blockSize
)

type row struct {
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

type props struct {
	generalCategory       string
	eastAsianWidth        string
	wordBreak             string
	sentenceBreak         string
	lineBreak             string
	originalGraphemeBreak string
	indicConjunctBreak    string
	defaultIgnorable      bool
	emojiModifier         bool
	emojiModifierBase     bool
	emojiComponent        bool
	emojiPresentation     bool
	extendedPictographic  bool
	emojiVSBase           bool
	whiteSpace            bool
	asciiHexDigit         bool
	hexDigit              bool
	dash                  bool
	diacritic             bool
	quotationMark         bool
	patternSyntax         bool
	patternWhiteSpace     bool
	variationSelector     bool
	noncharacter          bool
	unifiedIdeograph      bool
	upperDelta            int32
	lowerDelta            int32
	titleDelta            int32
	foldDelta             int32
}

const (
	widthMask      = 0x03
	zeroWidthFlag  = 0x04
	emojiVSFlag    = 0x01
	emojiPresFlag  = 0x02
	extPictoFlag   = 0x04
	whiteSpaceFlag = 0x08
)

const (
	asciiHexDigitFlag uint16 = 1 << iota
	hexDigitFlag
	dashFlag
	diacriticFlag
	quotationMarkFlag
	patternSyntaxFlag
	patternWhiteSpaceFlag
	variationSelectorFlag
	noncharacterFlag
	unifiedIdeographFlag
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

func main() {
	out := flag.String("out", "tables_gen.go", "output file")
	ucd := flag.String("ucd", "ucd", "Unicode Character Database directory")
	flag.Parse()

	rows, err := buildRows(*ucd)
	if err != nil {
		log.Fatal(err)
	}
	src, err := generate(rows)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*out, src, 0o644); err != nil {
		log.Fatal(err)
	}
}

func buildRows(ucd string) ([]row, error) {
	allProps := make([]props, maxRune+1)
	for i := range allProps {
		allProps[i] = props{
			generalCategory:       "Cn",
			eastAsianWidth:        "N",
			wordBreak:             "Other",
			sentenceBreak:         "Other",
			lineBreak:             "XX",
			originalGraphemeBreak: "Other",
			indicConjunctBreak:    "None",
		}
	}
	if err := loadUnicodeData(filepath.Join(ucd, "UnicodeData.txt"), allProps); err != nil {
		return nil, err
	}
	if err := loadCaseFolding(filepath.Join(ucd, "CaseFolding.txt"), allProps); err != nil {
		return nil, err
	}
	if err := loadRangeProperty(filepath.Join(ucd, "extracted", "DerivedGeneralCategory.txt"), func(r codeRange, value string) {
		for cp := r.start; cp <= r.end; cp++ {
			allProps[cp].generalCategory = value
		}
	}); err != nil {
		return nil, err
	}
	if err := loadDerivedCoreProperties(filepath.Join(ucd, "DerivedCoreProperties.txt"), allProps); err != nil {
		return nil, err
	}
	if err := loadRangeProperty(filepath.Join(ucd, "extracted", "DerivedEastAsianWidth.txt"), func(r codeRange, value string) {
		for cp := r.start; cp <= r.end; cp++ {
			allProps[cp].eastAsianWidth = value
		}
	}); err != nil {
		return nil, err
	}
	if err := loadRangeProperty(filepath.Join(ucd, "auxiliary", "WordBreakProperty.txt"), func(r codeRange, value string) {
		for cp := r.start; cp <= r.end; cp++ {
			allProps[cp].wordBreak = value
		}
	}); err != nil {
		return nil, err
	}
	if err := loadRangeProperty(filepath.Join(ucd, "auxiliary", "SentenceBreakProperty.txt"), func(r codeRange, value string) {
		for cp := r.start; cp <= r.end; cp++ {
			allProps[cp].sentenceBreak = value
		}
	}); err != nil {
		return nil, err
	}
	if err := loadRangeProperty(filepath.Join(ucd, "LineBreak.txt"), func(r codeRange, value string) {
		for cp := r.start; cp <= r.end; cp++ {
			allProps[cp].lineBreak = value
		}
	}); err != nil {
		return nil, err
	}
	if err := loadRangeProperty(filepath.Join(ucd, "auxiliary", "GraphemeBreakProperty.txt"), func(r codeRange, value string) {
		for cp := r.start; cp <= r.end; cp++ {
			allProps[cp].originalGraphemeBreak = value
		}
	}); err != nil {
		return nil, err
	}
	if err := loadRangeProperty(filepath.Join(ucd, "PropList.txt"), func(r codeRange, value string) {
		for cp := r.start; cp <= r.end; cp++ {
			switch value {
			case "White_Space":
				allProps[cp].whiteSpace = true
			case "ASCII_Hex_Digit":
				allProps[cp].asciiHexDigit = true
			case "Hex_Digit":
				allProps[cp].hexDigit = true
			case "Dash":
				allProps[cp].dash = true
			case "Diacritic":
				allProps[cp].diacritic = true
			case "Quotation_Mark":
				allProps[cp].quotationMark = true
			case "Pattern_Syntax":
				allProps[cp].patternSyntax = true
			case "Pattern_White_Space":
				allProps[cp].patternWhiteSpace = true
			case "Variation_Selector":
				allProps[cp].variationSelector = true
			case "Noncharacter_Code_Point":
				allProps[cp].noncharacter = true
			case "Unified_Ideograph":
				allProps[cp].unifiedIdeograph = true
			}
		}
	}); err != nil {
		return nil, err
	}
	if err := loadEmojiData(filepath.Join(ucd, "emoji", "emoji-data.txt"), allProps); err != nil {
		return nil, err
	}
	if err := loadEmojiVariationSequences(filepath.Join(ucd, "emoji", "emoji-variation-sequences.txt"), allProps); err != nil {
		return nil, err
	}

	rows := make([]row, maxRune+1)
	wordIDs := newInterner("Other")
	sentenceIDs := newInterner("Other")
	lineIDs := newInterner("XX")
	eawIDs := newInterner("N")
	gcIDs := newInterner("Cn")
	for cp, p := range allProps {
		gb := deriveGraphemeBreak(cp, p)
		width, zero := deriveWidth(cp, p, gb)
		packedWidth := uint8(width) & widthMask
		if zero {
			packedWidth |= zeroWidthFlag
		}
		var flags uint8
		if p.emojiVSBase {
			flags |= emojiVSFlag
		}
		if p.emojiPresentation {
			flags |= emojiPresFlag
		}
		if p.extendedPictographic {
			flags |= extPictoFlag
		}
		if p.whiteSpace {
			flags |= whiteSpaceFlag
		}
		var flags2 uint16
		if p.asciiHexDigit {
			flags2 |= asciiHexDigitFlag
		}
		if p.hexDigit {
			flags2 |= hexDigitFlag
		}
		if p.dash {
			flags2 |= dashFlag
		}
		if p.diacritic {
			flags2 |= diacriticFlag
		}
		if p.quotationMark {
			flags2 |= quotationMarkFlag
		}
		if p.patternSyntax {
			flags2 |= patternSyntaxFlag
		}
		if p.patternWhiteSpace {
			flags2 |= patternWhiteSpaceFlag
		}
		if p.variationSelector {
			flags2 |= variationSelectorFlag
		}
		if p.noncharacter {
			flags2 |= noncharacterFlag
		}
		if p.unifiedIdeograph {
			flags2 |= unifiedIdeographFlag
		}
		rows[cp] = row{
			gb:         gb,
			width:      packedWidth,
			wb:         wordIDs.id(p.wordBreak),
			sb:         sentenceIDs.id(p.sentenceBreak),
			lb:         lineIDs.id(p.lineBreak),
			eaw:        eawIDs.id(p.eastAsianWidth),
			gc:         gcIDs.id(p.generalCategory),
			flags:      flags,
			flags2:     flags2,
			upperDelta: p.upperDelta,
			lowerDelta: p.lowerDelta,
			titleDelta: p.titleDelta,
			foldDelta:  p.foldDelta,
		}
	}
	generatedNames = map[string][]string{
		"runtimeWordBreakNames":       wordIDs.names,
		"runtimeSentenceBreakNames":   sentenceIDs.names,
		"runtimeLineBreakNames":       lineIDs.names,
		"runtimeEastAsianWidthNames":  eawIDs.names,
		"runtimeGeneralCategoryNames": gcIDs.names,
	}
	return rows, nil
}

var generatedNames map[string][]string

type interner struct {
	ids   map[string]uint8
	names []string
}

func newInterner(defaultValue string) *interner {
	i := &interner{ids: map[string]uint8{}, names: []string{}}
	i.id(defaultValue)
	return i
}

func (i *interner) id(value string) uint8 {
	if id, ok := i.ids[value]; ok {
		return id
	}
	if len(i.names) > 0xff {
		panic("too many property values")
	}
	id := uint8(len(i.names))
	i.ids[value] = id
	i.names = append(i.names, value)
	return id
}

type unionFind struct {
	parent map[int]int
}

func newUnionFind() *unionFind {
	return &unionFind{parent: map[int]int{}}
}

func (u *unionFind) find(x int) int {
	parent, ok := u.parent[x]
	if !ok {
		u.parent[x] = x
		return x
	}
	if parent != x {
		parent = u.find(parent)
		u.parent[x] = parent
	}
	return parent
}

func (u *unionFind) union(a, b int) {
	ra := u.find(a)
	rb := u.find(b)
	if ra == rb {
		return
	}
	if ra < rb {
		u.parent[rb] = ra
	} else {
		u.parent[ra] = rb
	}
}

func (u *unionFind) groups() [][]int {
	byRoot := map[int][]int{}
	for cp := range u.parent {
		root := u.find(cp)
		byRoot[root] = append(byRoot[root], cp)
	}
	groups := make([][]int, 0, len(byRoot))
	for _, group := range byRoot {
		if len(group) < 2 {
			continue
		}
		sort.Ints(group)
		groups = append(groups, group)
	}
	return groups
}

type codeRange struct {
	start int
	end   int
}

func loadUnicodeData(path string, props []props) error {
	return eachDataLine(path, func(line string) error {
		f := strings.Split(line, ";")
		if len(f) < 15 {
			return fmt.Errorf("bad UnicodeData line: %q", line)
		}
		cp, err := parseCP(f[0])
		if err != nil {
			return err
		}
		props[cp].generalCategory = f[2]
		if f[12] != "" {
			mapped, err := parseCP(f[12])
			if err != nil {
				return err
			}
			props[cp].upperDelta = int32(mapped - cp)
		}
		if f[13] != "" {
			mapped, err := parseCP(f[13])
			if err != nil {
				return err
			}
			props[cp].lowerDelta = int32(mapped - cp)
		}
		if f[14] != "" {
			mapped, err := parseCP(f[14])
			if err != nil {
				return err
			}
			props[cp].titleDelta = int32(mapped - cp)
		}
		return nil
	})
}

func loadCaseFolding(path string, props []props) error {
	uf := newUnionFind()
	err := eachDataLine(path, func(line string) error {
		f := splitSemi(line)
		if len(f) < 3 || (f[1] != "C" && f[1] != "S") {
			return nil
		}
		cp, err := parseCP(f[0])
		if err != nil {
			return err
		}
		mapped, err := parseCP(f[2])
		if err != nil {
			return err
		}
		uf.union(cp, mapped)
		return nil
	})
	if err != nil {
		return err
	}
	for _, group := range uf.groups() {
		for i, cp := range group {
			next := group[(i+1)%len(group)]
			props[cp].foldDelta = int32(next - cp)
		}
	}
	return nil
}

func loadDerivedCoreProperties(path string, props []props) error {
	return eachDataLine(path, func(line string) error {
		f := splitSemi(line)
		if len(f) < 2 {
			return nil
		}
		r, err := parseRange(f[0])
		if err != nil {
			return err
		}
		if f[1] == "InCB" {
			if len(f) < 3 {
				return nil
			}
			for cp := r.start; cp <= r.end; cp++ {
				props[cp].indicConjunctBreak = f[2]
			}
			return nil
		}
		if f[1] == "Default_Ignorable_Code_Point" {
			for cp := r.start; cp <= r.end; cp++ {
				props[cp].defaultIgnorable = true
			}
		}
		return nil
	})
}

func loadEmojiData(path string, props []props) error {
	return loadRangeProperty(path, func(r codeRange, value string) {
		for cp := r.start; cp <= r.end; cp++ {
			switch value {
			case "Emoji_Modifier":
				props[cp].emojiModifier = true
			case "Emoji_Modifier_Base":
				props[cp].emojiModifierBase = true
			case "Emoji_Component":
				props[cp].emojiComponent = true
			case "Emoji_Presentation":
				props[cp].emojiPresentation = true
			case "Extended_Pictographic":
				props[cp].extendedPictographic = true
			}
		}
	})
}

func loadEmojiVariationSequences(path string, props []props) error {
	return eachDataLine(path, func(line string) error {
		fields := strings.Fields(line)
		if len(fields) < 2 || fields[1] != "FE0E" {
			return nil
		}
		cp, err := parseCP(fields[0])
		if err != nil {
			return err
		}
		props[cp].emojiVSBase = true
		return nil
	})
}

func loadRangeProperty(path string, fn func(codeRange, string)) error {
	return eachDataLine(path, func(line string) error {
		f := splitSemi(line)
		if len(f) < 2 {
			return nil
		}
		r, err := parseRange(f[0])
		if err != nil {
			return err
		}
		fn(r, f[1])
		return nil
	})
}

func eachDataLine(path string, fn func(string) error) error {
	b, err := os.ReadFile(path)
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

func parseCP(s string) (int, error) {
	n, err := strconv.ParseInt(strings.TrimSpace(s), 16, 32)
	return int(n), err
}

func deriveGraphemeBreak(cp int, p props) uint8 {
	if p.emojiModifier {
		return gbEmojiModifier
	}
	if p.emojiModifierBase {
		return gbEmojiModifierBase
	}
	if p.extendedPictographic {
		return gbExtendedPictographic
	}
	switch p.indicConjunctBreak {
	case "Extend":
		if cp == 0x200d {
			return gbZWJ
		}
		return gbIndicConjunctExtend
	case "Linker":
		return gbIndicConjunctLinker
	case "Consonant":
		return gbIndicConjunctConsonant
	}
	if p.originalGraphemeBreak == "Extend" {
		if cp == 0x200c {
			return gbZWNJ
		}
		return gbIndicConjunctExtend
	}
	switch p.originalGraphemeBreak {
	case "Control":
		return gbControl
	case "Prepend":
		return gbPrepend
	case "CR":
		return gbCR
	case "LF":
		return gbLF
	case "Regional_Indicator":
		return gbRegionalIndicator
	case "SpacingMark":
		return gbSpacingMark
	case "L":
		return gbL
	case "V":
		return gbV
	case "T":
		return gbT
	case "LV":
		return gbLV
	case "LVT":
		return gbLVT
	case "ZWJ":
		return gbZWJ
	default:
		return gbOther
	}
}

func deriveWidth(cp int, p props, gb uint8) (int, bool) {
	width := 1
	if p.generalCategory == "Cc" ||
		p.generalCategory == "Cs" ||
		p.generalCategory == "Zl" ||
		p.generalCategory == "Zp" {
		width = 0
	} else if cp == 0x00ad {
		width = 1
	} else if p.defaultIgnorable {
		width = 0
	} else if cp == 0x2e3a {
		width = 2
	} else if cp == 0x2e3b {
		width = 3
	} else if p.eastAsianWidth == "W" || p.eastAsianWidth == "F" {
		width = 2
	} else if gb == gbRegionalIndicator {
		width = 2
	}
	standalone := width
	if cp == 0x20e3 {
		standalone = 2
	}
	zero := width == 0 ||
		p.emojiModifier ||
		p.generalCategory == "Mn" ||
		p.generalCategory == "Me" ||
		gb == gbV ||
		gb == gbT ||
		gb == gbPrepend
	return standalone, zero
}

func generate(rows []row) ([]byte, error) {
	rowIDs := map[row]uint32{}
	var stage3 []row
	stage2Blocks := map[string]uint32{}
	var stage1 []uint32
	var stage2 []uint32

	for blockIndex := 0; blockIndex < numBlocks; blockIndex++ {
		var block [blockSize]uint32
		for i := 0; i < blockSize; i++ {
			r := rows[blockIndex*blockSize+i]
			id, ok := rowIDs[r]
			if !ok {
				id = uint32(len(stage3))
				rowIDs[r] = id
				stage3 = append(stage3, r)
			}
			block[i] = id
		}

		key := blockKey(block)
		offset, ok := stage2Blocks[key]
		if !ok {
			offset = uint32(len(stage2))
			stage2Blocks[key] = offset
			stage2 = append(stage2, block[:]...)
		}
		stage1 = append(stage1, offset)
	}

	var b bytes.Buffer
	b.WriteString("// Code generated by go generate; DO NOT EDIT.\n\n")
	b.WriteString("package uucode\n\n")
	writeUint32Slice(&b, "runtimeStage1", stage1)
	writeUint16Slice(&b, "runtimeStage2", stage2)
	for _, name := range []string{
		"runtimeWordBreakNames",
		"runtimeSentenceBreakNames",
		"runtimeLineBreakNames",
		"runtimeEastAsianWidthNames",
		"runtimeGeneralCategoryNames",
	} {
		writeStringSlice(&b, name, generatedNames[name])
	}
	b.WriteString("var runtimeStage3 = [...]runtimeRow{\n")
	for _, r := range stage3 {
		fmt.Fprintf(
			&b,
			"{gb:%d,width:%d,wb:%d,sb:%d,lb:%d,eaw:%d,gc:%d,flags:%d,flags2:%d,upperDelta:%d,lowerDelta:%d,titleDelta:%d,foldDelta:%d},",
			r.gb,
			r.width,
			r.wb,
			r.sb,
			r.lb,
			r.eaw,
			r.gc,
			r.flags,
			r.flags2,
			r.upperDelta,
			r.lowerDelta,
			r.titleDelta,
			r.foldDelta,
		)
	}
	b.WriteString("\n}\n")

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		return nil, err
	}
	return formatted, nil
}

func blockKey(block [blockSize]uint32) string {
	var b strings.Builder
	for _, v := range block {
		b.WriteString(strconv.FormatUint(uint64(v), 36))
		b.WriteByte(',')
	}
	return b.String()
}

func writeUint32Slice(b *bytes.Buffer, name string, values []uint32) {
	fmt.Fprintf(b, "var %s = [...]uint32{\n", name)
	for _, v := range values {
		fmt.Fprintf(b, "%d,", v)
	}
	b.WriteString("\n}\n\n")
}

func writeStringSlice(b *bytes.Buffer, name string, values []string) {
	fmt.Fprintf(b, "var %s = [...]string{\n", name)
	for _, v := range values {
		fmt.Fprintf(b, "%q,", v)
	}
	b.WriteString("\n}\n\n")
}

func writeUint16Slice(b *bytes.Buffer, name string, values []uint32) {
	fmt.Fprintf(b, "var %s = [...]uint16{\n", name)
	for _, v := range values {
		if v > 0xffff {
			panic(fmt.Sprintf("%s value %d overflows uint16", name, v))
		}
		fmt.Fprintf(b, "%d,", v)
	}
	b.WriteString("\n}\n\n")
}
