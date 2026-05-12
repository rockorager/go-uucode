# go-uucode

Go port of [`jacobsandlund/uucode`](https://github.com/jacobsandlund/uucode)'s core Unicode library behavior.

The package embeds Unicode 17 UCD files and provides:

- code point and UTF-8 iterators
- Unicode property lookup with `Get` and `GetAll`
- grapheme cluster iteration and `IsBreak`
- grapheme-aware terminal width helpers
- ASCII helpers matching the upstream `ascii` module

```go
package main

import (
	"fmt"

	uucode "tangled.org/rockorager.dev/go-uucode"
)

func main() {
	cp := rune(0x2200) // forall
	fmt.Println(uucode.GetAll(cp).GeneralCategory)

	it := uucode.NewUTF8GraphemeIterator("👩🏽‍🚀🇨🇭")
	g, _ := it.NextGrapheme()
	fmt.Println(g.Start, g.End)

	fmt.Println(uucode.UTF8Wcwidth("ò👨🏻‍❤️‍👨🏿_"))
}
```
