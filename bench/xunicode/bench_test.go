package xunicodebench

import (
	"testing"

	"charm.land/xunicode/grapheme"
	uucode "github.com/rockorager/go-uucode"
)

var benchmarkStrings = []struct {
	name string
	s    string
}{
	{"ASCII", "The quick brown fox jumps over the lazy dog. 0123456789"},
	{"Combining", "A\u0300B e\u0301 o\u0300 n\u0303 c\u0327 A\u0300B e\u0301 o\u0300 n\u0303 c\u0327"},
	{"Emoji", "😀😅😻👺👩🏽‍🚀🇨🇭👨🏻‍🍼👨🏻‍❤️‍👨🏿"},
	{"Mixed", "ASCII A\u0300 👩🏽‍🚀 🇨🇭 क्‍ष 한글 😀 _ end"},
}

func BenchmarkGraphemeIterator(b *testing.B) {
	for _, tc := range benchmarkStrings {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				it := uucode.NewGraphemeIterator(tc.s)
				for {
					_, ok := it.Next()
					if !ok {
						break
					}
					count++
				}
			}
			sinkInt = count
		})
		b.Run("xunicode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			input := []byte(tc.s)
			var count int
			for i := 0; i < b.N; i++ {
				seg := grapheme.NewSegmenter(input)
				for seg.Next() {
					count++
				}
			}
			sinkInt = count
		})
	}
}

var sinkInt int
