package common

func IsHexadecimal(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'A' || r > 'F') && (r < 'a' || r > 'f') {
			return false
		}
	}

	return true
}
