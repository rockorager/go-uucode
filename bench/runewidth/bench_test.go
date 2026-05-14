package runewidthbench

import (
	"testing"

	"github.com/mattn/go-runewidth"
	uucode "github.com/rockorager/go-uucode"
)

var benchmarkStrings = []struct {
	name string
	s    string
}{
	{"ASCII", "The quick brown fox jumps over the lazy dog. 0123456789"},
	{"Combining", "A\u0300B e\u0301 o\u0300 n\u0303 c\u0327 A\u0300B e\u0301 o\u0300 n\u0303 c\u0327"},
	{"HebrewMarks", "\u05d0\u0591 \u05d1\u05b0 \u05d2\u05b1 \u05d3\u05b2"},
	{"Indic", "\u0915\u093f \u0915\u094d\u200d\u0937 \u0c15\u0c4d\u200d\u0c37"},
	{"CJK", "界世界한글かなカナ"},
	{"MixedNoEmoji", "ASCII A\u0300 \u05d0\u0591 \u0915\u093f \u0915\u094d\u200d\u0937 한글 _ end"},
	{"RegionalIndicators", "\U0001f1fa\U0001f1f8 \U0001f1e8\U0001f1ed \U0001f1ef\U0001f1f5"},
	{"Emoji", "😀😅😻👺👩🏽‍🚀🇨🇭👨🏻‍🍼👨🏻‍❤️‍👨🏿"},
}

var benchmarkRunes = []struct {
	name string
	r    rune
}{
	{"ASCII", 'A'},
	{"Combining", '\u0300'},
	{"HebrewMark", '\u0591'},
	{"SpacingMark", '\u093f'},
	{"SoftHyphen", '\u00ad'},
	{"HangulFiller", '\u115f'},
	{"CJK", '界'},
	{"RegionalIndicator", '\U0001f1e6'},
	{"Emoji", '\U0001f600'},
}

func BenchmarkRuneWidth(b *testing.B) {
	cond := runewidth.NewCondition()
	cond.EastAsianWidth = false
	cond.StrictEmojiNeutral = false

	for _, tc := range benchmarkRunes {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var width int
			for i := 0; i < b.N; i++ {
				width += uucode.RuneWidth(tc.r)
			}
			sinkInt = width
		})
		b.Run("runewidth/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var width int
			for i := 0; i < b.N; i++ {
				width += cond.RuneWidth(tc.r)
			}
			sinkInt = width
		})
	}
}

func BenchmarkRuneWidthMixed(b *testing.B) {
	cond := runewidth.NewCondition()
	cond.EastAsianWidth = false
	cond.StrictEmojiNeutral = false
	runes := []rune{'A', '\u0300', '\u0591', '\u093f', '\u00ad', '\u115f', '界', '\U0001f1e6', '\U0001f600'}

	b.Run("uucode", func(b *testing.B) {
		b.ReportAllocs()
		var width int
		for i := 0; i < b.N; i++ {
			width += uucode.RuneWidth(runes[i%len(runes)])
		}
		sinkInt = width
	})
	b.Run("runewidth", func(b *testing.B) {
		b.ReportAllocs()
		var width int
		for i := 0; i < b.N; i++ {
			width += cond.RuneWidth(runes[i%len(runes)])
		}
		sinkInt = width
	})
}

func BenchmarkStringWidth(b *testing.B) {
	cond := runewidth.NewCondition()
	cond.EastAsianWidth = false
	cond.StrictEmojiNeutral = false

	for _, tc := range benchmarkStrings {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var width int
			for i := 0; i < b.N; i++ {
				width += uucode.StringWidth(tc.s)
			}
			sinkInt = width
		})
		b.Run("runewidth/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var width int
			for i := 0; i < b.N; i++ {
				width += cond.StringWidth(tc.s)
			}
			sinkInt = width
		})
	}
}

var sinkInt int
