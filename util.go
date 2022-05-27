package main

func IsHex(s string) bool {
	for _, ch := range []rune(s) {
		if '0' <= ch && ch <= '9' || 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F' {
			continue
		}
		return false
	}
	return true
}
