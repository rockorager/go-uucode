package uucode

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGeneratedRuntimeRows(t *testing.T) {
	if got := runtimeLookup(0x1f469).graphemeBreak(); got != GraphemeEmojiModifierBase {
		t.Fatalf("woman grapheme break: got %q", got)
	}
	if got := GeneralCategory('A'); got != GeneralCategoryLu {
		t.Fatalf("A general category: got %q", got)
	}
	if got := EastAsianWidth('A'); got != EastAsianWidthNa {
		t.Fatalf("A east asian width: got %q", got)
	}
	if got := WordBreak('A'); got != WordBreakALetter {
		t.Fatalf("A word break: got %q", got)
	}
	if got := SentenceBreak('A'); got != SentenceBreakUpper {
		t.Fatalf("A sentence break: got %q", got)
	}
	if got := LineBreak('A'); got != LineBreakAL {
		t.Fatalf("A line break: got %q", got)
	}
	if got := runtimeLookup(0x0300).wcwidthZeroInGrapheme(); !got {
		t.Fatal("combining grave should be zero width inside grapheme")
	}
	if got := runtimeLookup(0x3000).wcwidthStandalone(); got != 2 {
		t.Fatalf("ideographic space width: got %d", got)
	}
	if got := runtimeLookup(0x1f600).isEmojiPresentation(); !got {
		t.Fatal("grinning face should have emoji presentation")
	}
	if got := runtimeLookup(0x1f469).isExtendedPictographic(); !got {
		t.Fatal("woman should be extended pictographic")
	}
	if ToUpper('a') != 'A' || ToUpper('\u00b5') != '\u039c' || ToUpper('A') != 'A' {
		t.Fatal("ToUpper mismatch")
	}
	if ToLower('A') != 'a' || ToLower('\u0130') != 'i' || ToLower('a') != 'a' {
		t.Fatal("ToLower mismatch")
	}
	if ToTitle('\u01c6') != '\u01c5' || ToTitle('a') != 'A' || ToTitle('A') != 'A' {
		t.Fatal("ToTitle mismatch")
	}
	if SimpleFold('K') != 'k' || SimpleFold('k') != '\u212a' || SimpleFold('\u212a') != 'K' || SimpleFold('1') != '1' {
		t.Fatal("SimpleFold mismatch")
	}
	if !IsUpper('A') || IsUpper('a') {
		t.Fatal("IsUpper mismatch")
	}
	if !IsLower('a') || IsLower('A') {
		t.Fatal("IsLower mismatch")
	}
	if !IsTitle('\u01c5') || IsTitle('A') {
		t.Fatal("IsTitle mismatch")
	}
	if !IsLetter('A') || !IsLetter('界') || IsLetter('0') {
		t.Fatal("IsLetter mismatch")
	}
	if !IsNumber('0') || !IsDigit('0') || IsDigit('四') {
		t.Fatal("number predicate mismatch")
	}
	if !IsMark('\u0300') || IsMark('A') {
		t.Fatal("IsMark mismatch")
	}
	if !IsPunct('.') || !IsPunct('\u2014') || IsPunct('A') {
		t.Fatal("IsPunct mismatch")
	}
	if !IsSymbol('\u20ac') || !IsSymbol('\U0001f600') || IsSymbol('A') {
		t.Fatal("IsSymbol mismatch")
	}
	if !IsGraphic('A') || !IsGraphic('\u3000') || IsGraphic('\n') {
		t.Fatal("IsGraphic mismatch")
	}
	if !IsPrint('A') || !IsPrint(' ') || IsPrint('\u3000') || IsPrint('\n') {
		t.Fatal("IsPrint mismatch")
	}
	if !IsControl('\n') || IsControl('A') {
		t.Fatal("IsControl mismatch")
	}
	if !IsSpace(' ') || !IsSpace('\t') || !IsSpace('\u00a0') || !IsSpace('\u3000') || IsSpace('A') {
		t.Fatal("IsSpace mismatch")
	}
	if !IsASCIIHexDigit('f') || !IsASCIIHexDigit('F') || !IsASCIIHexDigit('9') || IsASCIIHexDigit('\uff26') {
		t.Fatal("IsASCIIHexDigit mismatch")
	}
	if !IsHexDigit('f') || !IsHexDigit('\uff26') || IsHexDigit('g') {
		t.Fatal("IsHexDigit mismatch")
	}
	if !IsDash('-') || !IsDash('\u2014') || IsDash('A') {
		t.Fatal("IsDash mismatch")
	}
	if !IsDiacritic('\u0300') || !IsDiacritic('^') || IsDiacritic('A') {
		t.Fatal("IsDiacritic mismatch")
	}
	if !IsQuotationMark('"') || !IsQuotationMark('\u201c') || IsQuotationMark('A') {
		t.Fatal("IsQuotationMark mismatch")
	}
	if !IsPatternSyntax('+') || !IsPatternSyntax('\u2192') || IsPatternSyntax('A') {
		t.Fatal("IsPatternSyntax mismatch")
	}
	if !IsPatternWhiteSpace('\u200e') || !IsPatternWhiteSpace('\n') || IsPatternWhiteSpace('A') {
		t.Fatal("IsPatternWhiteSpace mismatch")
	}
	if !IsVariationSelector('\ufe0f') || IsVariationSelector('A') {
		t.Fatal("IsVariationSelector mismatch")
	}
	if !IsNoncharacter('\ufffe') || !IsNoncharacter('\U0010ffff') || IsNoncharacter('A') {
		t.Fatal("IsNoncharacter mismatch")
	}
	if !IsUnifiedIdeograph('界') || !IsUnifiedIdeograph('\U00020000') || IsUnifiedIdeograph('A') {
		t.Fatal("IsUnifiedIdeograph mismatch")
	}
}

func TestGraphemeIteratorReadmeExample(t *testing.T) {
	s := "👩🏽‍🚀🇨🇭👨🏻‍🍼"
	it := NewGraphemeIterator(s)
	g, ok := it.Peek()
	if !ok || g.Start != 0 || g.End != 15 {
		t.Fatalf("peek first grapheme = %+v,%v", g, ok)
	}
	g, ok = it.Next()
	if !ok || s[g.Start:g.End] != "👩🏽‍🚀" {
		t.Fatalf("first grapheme = %q,%v", s[g.Start:g.End], ok)
	}
	g, ok = it.Next()
	if !ok || s[g.Start:g.End] != "🇨🇭" {
		t.Fatalf("second grapheme = %q,%v", s[g.Start:g.End], ok)
	}
	g, ok = it.Next()
	if !ok || s[g.Start:g.End] != "👨🏻‍🍼" {
		t.Fatalf("third grapheme = %q,%v", s[g.Start:g.End], ok)
	}
}

func TestEqualFold(t *testing.T) {
	tests := []struct {
		s    string
		t    string
		want bool
	}{
		{"hello", "HELLO", true},
		{"K", "\u212a", true},
		{"\u212a", "k", true},
		{"\u00b5", "\u039c", true},
		{"\u03a3", "\u03c2", true},
		{"stra\u00dfe", "STRASSE", false},
		{"abc", "ab", false},
		{"abc", "abd", false},
		{"界", "界", true},
		{"界", "畍", false},
	}
	for _, tc := range tests {
		if got := EqualFold(tc.s, tc.t); got != tc.want {
			t.Fatalf("EqualFold(%q, %q) = %v, want %v", tc.s, tc.t, got, tc.want)
		}
	}
}

func TestGraphemeBreakTestFile(t *testing.T) {
	content, err := os.ReadFile("ucd/auxiliary/GraphemeBreakTest.txt")
	if err != nil {
		t.Fatal(err)
	}
	lines := 0
	for lineNo, raw := range strings.Split(string(content), "\n") {
		line := testTrimComment(raw)
		if line == "" {
			continue
		}
		lines++
		tokens := strings.Fields(line)
		if len(tokens) < 4 || tokens[0] != "÷" {
			t.Fatalf("bad test line %d: %q", lineNo+1, raw)
		}
		state := BreakStateDefault
		cp1, err := testParseCP(tokens[1])
		if err != nil {
			t.Fatal(err)
		}
		for i := 2; i+1 < len(tokens); i += 2 {
			expectedBreak := tokens[i] == "÷"
			cp2, err := testParseCP(tokens[i+1])
			if err != nil {
				t.Fatal(err)
			}
			got := IsBreak(cp1, cp2, &state)
			gb1 := runtimeLookup(cp1).graphemeBreak()
			gb2 := runtimeLookup(cp2).graphemeBreak()
			if gb2 == GraphemeEmojiModifier && gb1 != GraphemeEmojiModifierBase {
				expectedBreak = true
			}
			if got != expectedBreak {
				t.Fatalf("line %d U+%04X/%s U+%04X/%s: got break=%v want %v", lineNo+1, cp1, gb1, cp2, gb2, got, expectedBreak)
			}
			cp1 = cp2
		}
	}
}

func testTrimComment(s string) string {
	if i := strings.IndexByte(s, '#'); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}

func testParseCP(s string) (rune, error) {
	n, err := strconv.ParseInt(strings.TrimSpace(s), 16, 32)
	return rune(n), err
}

func TestStringWidth(t *testing.T) {
	if got := StringWidth("ò👨🏻‍❤️‍👨🏿_"); got != 4 {
		t.Fatalf("readme wcwidth: got %d", got)
	}
	tests := map[string]int{
		"A\u0300B":             2,
		"😀AB":                  4,
		"\u200b":               0,
		"\u20e3":               0,
		"1\ufe0f\u20e3":        2,
		"\U0001f1e6":           1,
		"\u2601\ufe0f":         2,
		"\u2601\ufe0e":         1,
		"\U0001f1fa\U0001f1f8": 2,
		"\U0001f469\u200d\U0001f469\u200d\U0001f467\u200d\U0001f466_": 3,
		"\u1100\u1161":             2,
		"\u0915\u093f":             1,
		"\u0915\u094d\u200d\u0937": 2,
	}
	for s, want := range tests {
		if got := StringWidth(s); got != want {
			t.Fatalf("StringWidth(%q) = %d, want %d", s, got, want)
		}
	}
}

func TestRuneWidth(t *testing.T) {
	tests := map[rune]int{
		-1:       0,
		'A':      1,
		'\u00ad': 0,
		'\u0300': 0,
		'\u0591': 1,
		'\u070f': 0,
		'\u093f': 1,
		'\u0cf3': 0,
		'\u115f': 2,
		'\u20e3': 0,
		'\u4e00': 2,
		0x1f1e6:  1,
		0x2a6e0:  2,
		0x110000: 0,
	}
	for r, want := range tests {
		if got := RuneWidth(r); got != want {
			t.Fatalf("RuneWidth(%U) = %d, want %d", r, got, want)
		}
	}
}

func TestEastAsianWidthMissingRanges(t *testing.T) {
	if got := EastAsianWidth(0x2a6e0); got != EastAsianWidthW {
		t.Fatalf("EastAsianWidth(U+2A6E0) = %s, want W", got)
	}
}
