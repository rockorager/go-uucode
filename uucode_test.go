package uucode

import (
	"strings"
	"testing"
)

func TestREADMEProperties(t *testing.T) {
	if got := GetAll(0x2200).GeneralCategory; got != SymbolMath {
		t.Fatalf("general category: got %q", got)
	}
	if got := GetAll(0x03c2).SimpleUppercaseMapping; got != 0x03a3 {
		t.Fatalf("simple uppercase: got U+%04X", got)
	}
	if got := GetAll(0x21c1).Name; got != "RIGHTWARDS HARPOON WITH BARB DOWNWARDS" {
		t.Fatalf("name: got %q", got)
	}
	if got := string(GetAll(0x00df).UppercaseMapping); got != "SS" {
		t.Fatalf("uppercase mapping: got %q", got)
	}
}

func TestUTF8Iterator(t *testing.T) {
	it := NewUTF8Iterator("😀😅😻👺")
	expect := []struct {
		cp CodePoint
		i  int
	}{
		{0x1f600, 4},
		{0x1f605, 8},
		{0x1f63b, 12},
		{0x1f47a, 16},
	}
	for _, want := range expect {
		got, ok := it.Next()
		if !ok || got != want.cp || it.I != want.i {
			t.Fatalf("Next = U+%04X,%v i=%d; want U+%04X i=%d", got, ok, it.I, want.cp, want.i)
		}
	}
	if _, ok := it.Next(); ok {
		t.Fatal("expected EOF")
	}
}

func TestGraphemeIteratorReadmeExample(t *testing.T) {
	s := "👩🏽‍🚀🇨🇭👨🏻‍🍼"
	it := NewUTF8GraphemeIterator(s)
	g, ok := it.PeekGrapheme()
	if !ok || g.Start != 0 || g.End != 15 {
		t.Fatalf("peek first grapheme = %+v,%v", g, ok)
	}
	g, ok = it.NextGrapheme()
	if !ok || s[g.Start:g.End] != "👩🏽‍🚀" {
		t.Fatalf("first grapheme = %q,%v", s[g.Start:g.End], ok)
	}
	g, ok = it.NextGrapheme()
	if !ok || s[g.Start:g.End] != "🇨🇭" {
		t.Fatalf("second grapheme = %q,%v", s[g.Start:g.End], ok)
	}
	g, ok = it.NextGrapheme()
	if !ok || s[g.Start:g.End] != "👨🏻‍🍼" {
		t.Fatalf("third grapheme = %q,%v", s[g.Start:g.End], ok)
	}
}

func TestGraphemeBreakTestFile(t *testing.T) {
	content, err := ucdFS.ReadFile("ucd/auxiliary/GraphemeBreakTest.txt")
	if err != nil {
		t.Fatal(err)
	}
	lines := 0
	for lineNo, raw := range strings.Split(string(content), "\n") {
		line := trimComment(raw)
		if line == "" {
			continue
		}
		lines++
		tokens := strings.Fields(line)
		if len(tokens) < 4 || tokens[0] != "÷" {
			t.Fatalf("bad test line %d: %q", lineNo+1, raw)
		}
		state := BreakStateDefault
		cp1, err := parseCP(tokens[1])
		if err != nil {
			t.Fatal(err)
		}
		for i := 2; i+1 < len(tokens); i += 2 {
			expectedBreak := tokens[i] == "÷"
			cp2, err := parseCP(tokens[i+1])
			if err != nil {
				t.Fatal(err)
			}
			got := IsBreak(cp1, cp2, &state)
			gb1 := GetAll(cp1).GraphemeBreak
			gb2 := GetAll(cp2).GraphemeBreak
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

func TestWcwidth(t *testing.T) {
	if got := UTF8Wcwidth("ò👨🏻‍❤️‍👨🏿_"); got != 4 {
		t.Fatalf("readme wcwidth: got %d", got)
	}
	tests := map[string]int{
		"A\u0300B":             2,
		"😀AB":                  4,
		"\u200b":               0,
		"\u20e3":               2,
		"\u1f1e6":              2,
		"\u2601\ufe0f":         2,
		"\u2601\ufe0e":         1,
		"\U0001f1fa\U0001f1f8": 2,
		"\U0001f469\u200d\U0001f469\u200d\U0001f467\u200d\U0001f466_": 3,
		"\u1100\u1161":             2,
		"\u0915\u094d\u200d\u0937": 2,
	}
	for s, want := range tests {
		if got := UTF8Wcwidth(s); got != want {
			t.Fatalf("UTF8Wcwidth(%q) = %d, want %d", s, got, want)
		}
	}
}
