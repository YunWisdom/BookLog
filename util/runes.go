

package util

// IsNumOrLetter checks the specified rune is number or letter.
func IsNumOrLetter(r rune) bool {
	return IsLetter(r) || ('0' <= r && '9' >= r)
}

// IsLetter checks the specified run is letter.
func IsLetter(r rune) bool {
	return 'a' <= r && 'z' >= r || 'A' <= r && 'Z' >= r
}
