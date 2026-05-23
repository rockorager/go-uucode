package uucode

import "testing"

func TestGraphemeWidthIteratorMatchesSplitStringWidth(t *testing.T) {
	cases := []string{
		"",
		"abc",
		"hello\nworld",
		"A\u0300B e\u0301",
		"😀😅😻👺",
		"👩🏽‍🚀🇨🇭👨🏻‍🍼👨🏻‍❤️‍👨🏿",
		"क्‍ष 한글 😀 _ end",
		"text\ufe0f text\ufe0e",
	}
	for _, s := range cases {
		it := NewGraphemeIterator(s)
		wit := NewGraphemeWidthIterator(s)
		var total int
		for {
			g, ok := it.Next()
			wg, wok := wit.Next()
			if ok != wok {
				t.Fatalf("%q: iterator mismatch ok=%v wok=%v", s, ok, wok)
			}
			if !ok {
				break
			}
			if g.Start != wg.Start || g.End != wg.End {
				t.Fatalf("%q: bounds mismatch got %d:%d want %d:%d", s, wg.Start, wg.End, g.Start, g.End)
			}
			want := StringWidth(s[g.Start:g.End])
			if wg.Width != want {
				t.Fatalf("%q[%d:%d]: width=%d want=%d cluster=%q", s, g.Start, g.End, wg.Width, want, s[g.Start:g.End])
			}
			total += wg.Width
		}
		if total != StringWidth(s) {
			t.Fatalf("%q: total width=%d want=%d", s, total, StringWidth(s))
		}
	}
}
