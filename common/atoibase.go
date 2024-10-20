package common

// Returns the decimal interpretation of [s] being in base [base].
func AtoiBase(s, base string) int {
	if len(base) <= 1 || base[0] == '+' || base[0] == '-' {
		return 0
	}

	res := 0
	negative := false
	if s[0] == '-' {
		negative = true
		s = s[1:]
	}

	for i, pow := len(s)-1, 1; i >= 0; i, pow = i-1, pow*len(base) {
		for j := 0; j < len(base); j++ {
			if s[i] == base[j] {
				res += j * pow
				break
			}
		}
	}

	if negative {
		res = -res
		if res > 0 {
			return 0
		}
	}

	return res
}
