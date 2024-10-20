package common

func ContainsIP(s string) bool {
	if len(s) < 7 {
		return false
	}

	firstNum := -1
	numAmt := 0

	for i, r := range s {
		if r >= '0' && r <= '9' {
			if firstNum == -1 {
				firstNum = i
			}
		} else if r == '.' {
			if firstNum != -1 && i-firstNum <= 3 {
				firstNum = -1
				numAmt++
			}
		} else {
			if numAmt == 3 && firstNum != -1 && IsNumeric(s[firstNum:i-1]) {
				return true
			}

			firstNum = -1
			numAmt = 0
		}
	}

	return numAmt == 3 && IsNumeric(s[firstNum:])
}
