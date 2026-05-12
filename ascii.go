package uucode

func IsAlphanumeric(c CodePoint) bool    { return IsAlphabetic(c) || IsDigit(c) }
func IsAlphabeticASCII(c CodePoint) bool { return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') }
func IsAlphabetic(c CodePoint) bool      { return IsAlphabeticASCII(c) }
func IsControl(c CodePoint) bool         { return c <= 0x1f || c == 0x7f }
func IsDigit(c CodePoint) bool           { return c >= '0' && c <= '9' }
func IsLower(c CodePoint) bool           { return c >= 'a' && c <= 'z' }
func IsPrint(c CodePoint) bool           { return IsASCII(c) && !IsControl(c) }
func IsWhitespace(c CodePoint) bool      { return c == ' ' || (c >= '\t' && c <= '\r') }
func IsUpper(c CodePoint) bool           { return c >= 'A' && c <= 'Z' }
func IsHex(c CodePoint) bool             { return IsDigit(c) || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f') }
func IsASCII(c CodePoint) bool           { return c >= 0 && c < 128 }

func ToUpper(c CodePoint) CodePoint {
	if IsLower(c) {
		return c - ('a' - 'A')
	}
	return c
}

func ToLower(c CodePoint) CodePoint {
	if IsUpper(c) {
		return c + ('a' - 'A')
	}
	return c
}
