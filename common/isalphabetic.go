package common

func IsAlphabetic(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return false
		}
	}

	return true
}
