package uucode

import "unicode/utf8"

// EqualFold reports whether s and t are equal under Unicode simple case folding.
func EqualFold(s, t string) bool {
	i := 0
	for n := min(len(s), len(t)); i < n; i++ {
		cs := s[i]
		ct := t[i]
		if cs|ct >= utf8.RuneSelf {
			goto hasUnicode
		}
		if cs == ct {
			continue
		}
		if ct < cs {
			cs, ct = ct, cs
		}
		if !('A' <= cs && cs <= 'Z' && ct == cs+'a'-'A') {
			return false
		}
	}
	return len(s) == len(t)

hasUnicode:
	s = s[i:]
	t = t[i:]
	for _, rs := range s {
		if len(t) == 0 {
			return false
		}
		rt, size := utf8.DecodeRuneInString(t)
		t = t[size:]
		if rt == rs {
			continue
		}
		if rt < rs {
			rs, rt = rt, rs
		}
		if rt < utf8.RuneSelf {
			if 'A' <= rs && rs <= 'Z' && rt == rs+'a'-'A' {
				continue
			}
			return false
		}
		if !equalFoldRuneOrdered(rs, rt) {
			return false
		}
	}
	return len(t) == 0
}

func equalFoldRuneOrdered(a, b rune) bool {
	r := SimpleFold(a)
	for r != a && r < b {
		r = SimpleFold(r)
	}
	return r == b
}
