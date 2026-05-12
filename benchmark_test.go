package uucode

import (
	"strings"
	"testing"
	"unicode"
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
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				it := NewGraphemeIterator(tc.s)
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
	}
}

func BenchmarkStringWidth(b *testing.B) {
	for _, tc := range benchmarkStrings {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var width int
			for i := 0; i < b.N; i++ {
				width += StringWidth(tc.s)
			}
			sinkInt = width
		})
	}
}

var benchmarkRunes = []struct {
	name string
	r    rune
}{
	{"ASCIIUpper", 'A'},
	{"ASCIILower", 'a'},
	{"ASCIIDigit", '0'},
	{"CombiningMark", '\u0300'},
	{"CJK", '界'},
	{"Emoji", '\U0001f600'},
}

func BenchmarkIsUpper(b *testing.B) {
	for _, tc := range benchmarkRunes {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if IsUpper(r) {
					count++
				}
			}
			sinkInt = count
		})
		b.Run("stdlib/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if unicode.IsUpper(r) {
					count++
				}
			}
			sinkInt = count
		})
	}
}

func BenchmarkIsLower(b *testing.B) {
	for _, tc := range benchmarkRunes {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if IsLower(r) {
					count++
				}
			}
			sinkInt = count
		})
		b.Run("stdlib/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if unicode.IsLower(r) {
					count++
				}
			}
			sinkInt = count
		})
	}
}

func BenchmarkIsLetter(b *testing.B) {
	for _, tc := range benchmarkRunes {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if IsLetter(r) {
					count++
				}
			}
			sinkInt = count
		})
		b.Run("stdlib/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if unicode.IsLetter(r) {
					count++
				}
			}
			sinkInt = count
		})
	}
}

func BenchmarkIsNumber(b *testing.B) {
	for _, tc := range benchmarkRunes {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if IsNumber(r) {
					count++
				}
			}
			sinkInt = count
		})
		b.Run("stdlib/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if unicode.IsNumber(r) {
					count++
				}
			}
			sinkInt = count
		})
	}
}

func BenchmarkIsDigit(b *testing.B) {
	for _, tc := range benchmarkRunes {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if IsDigit(r) {
					count++
				}
			}
			sinkInt = count
		})
		b.Run("stdlib/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if unicode.IsDigit(r) {
					count++
				}
			}
			sinkInt = count
		})
	}
}

func BenchmarkIsMark(b *testing.B) {
	for _, tc := range benchmarkRunes {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if IsMark(r) {
					count++
				}
			}
			sinkInt = count
		})
		b.Run("stdlib/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if unicode.IsMark(r) {
					count++
				}
			}
			sinkInt = count
		})
	}
}

func BenchmarkIsControl(b *testing.B) {
	for _, tc := range benchmarkRunes {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if IsControl(r) {
					count++
				}
			}
			sinkInt = count
		})
		b.Run("stdlib/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if unicode.IsControl(r) {
					count++
				}
			}
			sinkInt = count
		})
	}
}

func BenchmarkIsTitle(b *testing.B) {
	benchmarkPredicate(b, IsTitle, unicode.IsTitle)
}

func BenchmarkIsPunct(b *testing.B) {
	benchmarkPredicate(b, IsPunct, unicode.IsPunct)
}

func BenchmarkIsSymbol(b *testing.B) {
	benchmarkPredicate(b, IsSymbol, unicode.IsSymbol)
}

func BenchmarkIsGraphic(b *testing.B) {
	benchmarkPredicate(b, IsGraphic, unicode.IsGraphic)
}

func BenchmarkIsPrint(b *testing.B) {
	benchmarkPredicate(b, IsPrint, unicode.IsPrint)
}

func BenchmarkIsSpace(b *testing.B) {
	benchmarkPredicate(b, IsSpace, unicode.IsSpace)
}

func BenchmarkIsASCIIHexDigit(b *testing.B) {
	benchmarkPropertyPredicate(b, IsASCIIHexDigit, func(r rune) bool { return unicode.Is(unicode.ASCII_Hex_Digit, r) })
}

func BenchmarkIsHexDigit(b *testing.B) {
	benchmarkPropertyPredicate(b, IsHexDigit, func(r rune) bool { return unicode.Is(unicode.Hex_Digit, r) })
}

func BenchmarkIsDash(b *testing.B) {
	benchmarkPropertyPredicate(b, IsDash, func(r rune) bool { return unicode.Is(unicode.Dash, r) })
}

func BenchmarkIsDiacritic(b *testing.B) {
	benchmarkPropertyPredicate(b, IsDiacritic, func(r rune) bool { return unicode.Is(unicode.Diacritic, r) })
}

func BenchmarkIsQuotationMark(b *testing.B) {
	benchmarkPropertyPredicate(b, IsQuotationMark, func(r rune) bool { return unicode.Is(unicode.Quotation_Mark, r) })
}

func BenchmarkIsPatternSyntax(b *testing.B) {
	benchmarkPropertyPredicate(b, IsPatternSyntax, func(r rune) bool { return unicode.Is(unicode.Pattern_Syntax, r) })
}

func BenchmarkIsPatternWhiteSpace(b *testing.B) {
	benchmarkPropertyPredicate(b, IsPatternWhiteSpace, func(r rune) bool { return unicode.Is(unicode.Pattern_White_Space, r) })
}

func BenchmarkIsVariationSelector(b *testing.B) {
	benchmarkPropertyPredicate(b, IsVariationSelector, func(r rune) bool { return unicode.Is(unicode.Variation_Selector, r) })
}

func BenchmarkIsNoncharacter(b *testing.B) {
	benchmarkPropertyPredicate(b, IsNoncharacter, func(r rune) bool { return unicode.Is(unicode.Noncharacter_Code_Point, r) })
}

func BenchmarkIsUnifiedIdeograph(b *testing.B) {
	benchmarkPropertyPredicate(b, IsUnifiedIdeograph, func(r rune) bool { return unicode.Is(unicode.Unified_Ideograph, r) })
}

func BenchmarkToUpper(b *testing.B) {
	benchmarkRuneMap(b, ToUpper, unicode.ToUpper)
}

func BenchmarkToLower(b *testing.B) {
	benchmarkRuneMap(b, ToLower, unicode.ToLower)
}

func BenchmarkToTitle(b *testing.B) {
	benchmarkRuneMap(b, ToTitle, unicode.ToTitle)
}

func BenchmarkSimpleFold(b *testing.B) {
	benchmarkRuneMap(b, SimpleFold, unicode.SimpleFold)
}

func BenchmarkEqualFold(b *testing.B) {
	for _, tc := range benchmarkEqualFoldStrings {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				if EqualFold(tc.s, tc.t) {
					count++
				}
			}
			sinkInt = count
		})
		b.Run("stdlib/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				if strings.EqualFold(tc.s, tc.t) {
					count++
				}
			}
			sinkInt = count
		})
	}
}

func benchmarkPredicate(b *testing.B, ours func(rune) bool, stdlib func(rune) bool) {
	for _, tc := range benchmarkRunes {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if ours(r) {
					count++
				}
			}
			sinkInt = count
		})
		b.Run("stdlib/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkCategoryRunes[(i+int(tc.r))&benchmarkCategoryMask]
				if stdlib(r) {
					count++
				}
			}
			sinkInt = count
		})
	}
}

func benchmarkPropertyPredicate(b *testing.B, ours func(rune) bool, stdlib func(rune) bool) {
	for _, tc := range benchmarkRunes {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkPropertyRunes[(i+int(tc.r))&benchmarkPropertyMask]
				if ours(r) {
					count++
				}
			}
			sinkInt = count
		})
		b.Run("stdlib/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var count int
			for i := 0; i < b.N; i++ {
				r := benchmarkPropertyRunes[(i+int(tc.r))&benchmarkPropertyMask]
				if stdlib(r) {
					count++
				}
			}
			sinkInt = count
		})
	}
}

func benchmarkRuneMap(b *testing.B, ours func(rune) rune, stdlib func(rune) rune) {
	for _, tc := range benchmarkRunes {
		b.Run("uucode/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var sum rune
			for i := 0; i < b.N; i++ {
				r := benchmarkMappingRunes[(i+int(tc.r))&benchmarkMappingMask]
				sum += ours(r)
			}
			sinkRune = sum
		})
		b.Run("stdlib/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			var sum rune
			for i := 0; i < b.N; i++ {
				r := benchmarkMappingRunes[(i+int(tc.r))&benchmarkMappingMask]
				sum += stdlib(r)
			}
			sinkRune = sum
		})
	}
}

var benchmarkCategoryRunes = [...]rune{
	'A',
	'z',
	'0',
	' ',
	'\u0300',
	'\u03a9',
	'\u05d0',
	'\u3042',
	'\u3000',
	'\U0001f600',
	'\u200d',
	'\U0001f1fa',
	'\U00010400',
	'\U0001d7ce',
	'\u20ac',
	'\uff21',
	'\u01c5',
	'.',
	'\u2014',
	'\u00a0',
	'\u0085',
	'\u2028',
	'\u2029',
	'\u2003',
	'\t',
	'\n',
	'\u00b2',
	'\u2167',
	'\u203f',
	'\u2118',
	'\u2665',
	'\u037e',
}

const benchmarkCategoryMask = len(benchmarkCategoryRunes) - 1

var benchmarkPropertyRunes = [...]rune{
	'A',
	'F',
	'f',
	'g',
	'0',
	'\uff19',
	'\uff26',
	'-',
	'\u2014',
	'^',
	'\u0300',
	'"',
	'\u201c',
	'+',
	'\u2192',
	' ',
	'\n',
	'\u200e',
	'\ufe0f',
	'\U000e0100',
	'\ufffe',
	'\U0010ffff',
	'界',
	'\U00020000',
	'z',
	'_',
	'\u00a0',
	'\U0001f600',
	'\u2028',
	'\u037e',
	'\u2665',
	'\u3000',
}

const benchmarkPropertyMask = len(benchmarkPropertyRunes) - 1

var benchmarkMappingRunes = [...]rune{
	'A',
	'a',
	'K',
	'k',
	'\u212a',
	'\u00b5',
	'\u039c',
	'\u03bc',
	'\u0130',
	'i',
	'\u01c4',
	'\u01c5',
	'\u01c6',
	'\u03a3',
	'\u03c3',
	'\u03c2',
	'0',
	'界',
	'\U00010400',
	'\U00010428',
	'\U0001e900',
	'\U0001e922',
	'\U0001f600',
	'\uff21',
	'\uff41',
	'\u017f',
	'\u1e9e',
	'\u00df',
	'\u2126',
	'\u03a9',
	'\u03c9',
	'\u1f88',
}

const benchmarkMappingMask = len(benchmarkMappingRunes) - 1

var benchmarkEqualFoldStrings = []struct {
	name string
	s    string
	t    string
}{
	{"ASCIIEqual", "The Quick Brown Fox 0123456789", "the quick brown fox 0123456789"},
	{"ASCIIMiss", "The Quick Brown Fox 0123456789", "the quick brown fax 0123456789"},
	{"Kelvin", "Kilo Kelvin \u212a", "kilo kelvin k"},
	{"GreekSigma", "\u03a3\u03c3\u03c2", "\u03c2\u03a3\u03c3"},
	{"Mixed", "Stra\u00dfe \u212a \u00b5 \u03a3 \U00010400", "stra\u00dfe k \u039c \u03c2 \U00010428"},
	{"LengthMiss", "abcdef", "abcde"},
}

var sinkInt int
var sinkRune rune
