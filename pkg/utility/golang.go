package utility

import "unicode"

func CheckGolangIdentifier(identifier string) bool {
	if identifier == "" {
		return false
	}

	for i, r := range identifier {
		if i == 0 && !unicode.IsLetter(r) && r != '_' {
			return false
		}
		if i > 0 && !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}

	return true
}
