package unisegbench

import (
	"testing"

	"github.com/rivo/uniseg"
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
		b.Run("uniseg/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				g := uniseg.NewGraphemes(tc.s)
				for g.Next() {
					count++
				}
			}
			sinkInt = count
		})
	}
}

func BenchmarkStringWidth(b *testing.B) {
	for _, tc := range benchmarkStrings {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var width int
			for i := 0; i < b.N; i++ {
				width += uucode.StringWidth(tc.s)
			}
			sinkInt = width
		})
		b.Run("uniseg/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var width int
			for i := 0; i < b.N; i++ {
				width += uniseg.StringWidth(tc.s)
			}
			sinkInt = width
		})
	}
}

var sinkInt int
