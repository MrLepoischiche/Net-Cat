package common

func IsAlphaNum(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return false
		}
	}

	return true
}
