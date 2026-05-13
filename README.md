# go-uucode

`go-uucode` is a small Go Unicode segmentation and width package inspired by
Jacob Sandlund's excellent [`uucode`](https://github.com/jacobsandlund/uucode).
Jacob's Zig implementation does the hard architectural work here: generated
Unicode tables, compact property rows, and a fast staged lookup strategy. This
package ports that table-first approach to Go.

It provides:

- extended grapheme cluster iteration over Go strings
- grapheme-aware terminal cell width with `StringWidth`
- typed lookup APIs for generated Unicode category, break, binary, emoji, width, and case properties
- no runtime UCD parser, cache, or fallback path

## Usage

```go
package main

import (
	"fmt"

	"github.com/rockorager/go-uucode"
)

func main() {
	s := "👩🏽‍🚀🇨🇭A\u0300"

	it := uucode.NewGraphemeIterator(s)
	for {
		g, ok := it.Next()
		if !ok {
			break
		}
		fmt.Printf("%q [%d:%d]\n", s[g.Start:g.End], g.Start, g.End)
	}

	fmt.Println(uucode.StringWidth("ò👨🏻‍❤️‍👨🏿_"))
	fmt.Println(uucode.IsLetter('界'), uucode.WordBreak('A'), uucode.LineBreak(' '))
}
```

## Benchmarks

Benchmarks below were run on an Apple M4 Max with Go 1.26.1. Both libraries
reported `0 B/op` and `0 allocs/op`.

| Public API benchmark | go-uucode ns/op | uniseg ns/op | Speedup |
|---|---:|---:|---:|
| Grapheme ASCII | 361.6 | 3326 | 9.20x |
| Grapheme Combining | 254.6 | 1810 | 7.11x |
| Grapheme Emoji | 184.7 | 1863 | 10.09x |
| Grapheme Mixed | 255.5 | 2452 | 9.60x |
| Width ASCII | 33.75 | 489.2 | 14.49x |
| Width Combining | 286.6 | 331.2 | 1.16x |
| Width Emoji | 217.3 | 444.1 | 2.04x |
| Width Mixed | 250.9 | 500.1 | 1.99x |

Predicate APIs are benchmarked against Go's `unicode` package on a rotating
32-rune corpus. These are speed comparisons against the public stdlib APIs; the
local Go toolchain reports `unicode.Version == "15.0.0"` while go-uucode ships
Unicode 17 data. The rows below show the mean of the six benchmark subcases:

| Predicate benchmark | go-uucode ns/op | stdlib ns/op | Speedup |
|---|---:|---:|---:|
| IsUpper | 1.67 | 5.73 | 3.42x |
| IsLower | 1.68 | 5.77 | 3.44x |
| IsTitle | 2.63 | 2.26 | 0.86x |
| IsLetter | 1.75 | 6.50 | 3.72x |
| IsNumber | 1.69 | 5.24 | 3.11x |
| IsDigit | 1.68 | 4.68 | 2.78x |
| IsMark | 1.75 | 6.52 | 3.72x |
| IsPunct | 2.52 | 6.24 | 2.47x |
| IsSymbol | 2.54 | 6.31 | 2.48x |
| IsGraphic | 2.56 | 22.23 | 8.69x |
| IsPrint | 2.62 | 22.26 | 8.50x |
| IsControl | 0.37 | 0.68 | 1.83x |
| IsSpace | 2.66 | 3.33 | 1.25x |

Generated binary property APIs are benchmarked against `unicode.Is` with the
matching stdlib range table on a property-focused 32-rune corpus:

| Binary property benchmark | go-uucode ns/op | stdlib ns/op | Speedup |
|---|---:|---:|---:|
| IsASCIIHexDigit | 2.44 | 2.58 | 1.06x |
| IsHexDigit | 2.45 | 3.45 | 1.41x |
| IsDash | 2.47 | 3.83 | 1.55x |
| IsDiacritic | 2.45 | 5.25 | 2.14x |
| IsQuotationMark | 2.44 | 3.85 | 1.57x |
| IsPatternSyntax | 2.45 | 4.40 | 1.80x |
| IsPatternWhiteSpace | 2.44 | 3.31 | 1.36x |
| IsVariationSelector | 2.42 | 3.04 | 1.26x |
| IsNoncharacter | 2.44 | 3.13 | 1.28x |
| IsUnifiedIdeograph | 2.47 | 2.95 | 1.20x |

Simple case mapping APIs are benchmarked against the matching `unicode`
functions on a case-focused 32-rune corpus:

| Case mapping benchmark | go-uucode ns/op | stdlib ns/op | Speedup |
|---|---:|---:|---:|
| ToUpper | 2.01 | 6.86 | 3.42x |
| ToLower | 1.92 | 6.83 | 3.56x |
| ToTitle | 1.92 | 6.85 | 3.57x |
| SimpleFold | 1.93 | 6.37 | 3.30x |

String case folding is benchmarked against `strings.EqualFold`:

| EqualFold benchmark | go-uucode ns/op | stdlib ns/op | Speedup |
|---|---:|---:|---:|
| ASCII equal | 12.38 | 13.00 | 1.05x |
| ASCII miss | 8.92 | 9.02 | 1.01x |
| Kelvin | 9.80 | 9.29 | 0.95x |
| Greek sigma | 18.72 | 28.79 | 1.54x |
| Mixed Unicode | 31.94 | 41.27 | 1.29x |
| Length miss | 2.68 | 2.63 | 0.98x |

Run the package benchmarks:

```sh
go test -run '^$' -bench . -benchmem
```

The comparison against `github.com/rivo/uniseg` lives in a separate nested
module so `uniseg` is not a dependency of this package:

```sh
cd bench/uniseg
go test -run '^$' -bench . -benchmem
```

## Generated Tables

The package ships Unicode 17 source files and generates packed runtime tables.
The hot path uses three stages:

- `runtimeStage1` indexes 256-code-point blocks by `cp >> 8`
- `runtimeStage2` indexes the low byte within deduplicated blocks
- `runtimeStage3` stores deduplicated packed property rows

Regenerate after changing UCD files or generator logic:

```sh
go generate ./...
```

The generated runtime rows store compact fields for grapheme segmentation,
terminal-width calculation, general category predicates, word/sentence/line
break properties, East Asian Width, PropList binary properties, simple case
mapping, simple case folding, and emoji properties used by the public lookup
functions.

## Attribution

The design is based on the real
[`jacobsandlund/uucode`](https://github.com/jacobsandlund/uucode). If you are
interested in the original implementation, Unicode table generation strategy,
or a Zig library for this problem space, start there.
