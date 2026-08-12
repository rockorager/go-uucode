package uucode

import "testing"

func TestGraphemeWidthIteratorTagSequences(t *testing.T) {
	cases := map[string]string{
		"scotland": "\U0001f3f4\U000e0067\U000e0062\U000e0073\U000e0063\U000e0074\U000e007f",
		"england":  "\U0001f3f4\U000e0067\U000e0062\U000e0065\U000e006e\U000e0067\U000e007f",
		"wales":    "\U0001f3f4\U000e0067\U000e0062\U000e0077\U000e006c\U000e0073\U000e007f",
	}
	for name, s := range cases {
		it := NewGraphemeWidthIterator(s)
		g, ok := it.Next()
		if !ok {
			t.Fatalf("%s: no cluster", name)
		}
		if g.End != len(s) {
			t.Fatalf("%s: cluster is %d bytes, want %d", name, g.End, len(s))
		}
		if g.Width != 2 {
			t.Fatalf("%s: cluster width = %d, want 2", name, g.Width)
		}
		if _, ok := it.Next(); ok {
			t.Fatalf("%s: expected exactly one cluster", name)
		}
	}
}

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
		"a\U0001f3f4\U000e0067\U000e0062\U000e0073\U000e0063\U000e0074\U000e007fb",
		"\U000e0067\U000e0068",
		"\uff76\uff9e\uff77\uff9e",
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
