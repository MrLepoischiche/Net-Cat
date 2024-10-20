package common

/*
	Chunk separates the string of data into sub-strings of N elements, all stored inside the returning array.

	... Given the string is not empty, and the number of elements per substring is neither Zero nor bigger than the initial string's length.
*/
func StrChunk(str string, n int) []string {
	if len(str) == 0 || n == 0 || n > len(str) {
		return []string{}
	}
	if n == len(str) {
		return []string{str}
	}

	res := make([]string, 0)
	tmp := ""

	for i := 0; i < len(str); i++ {
		if i%n == 0 && i != 0 {
			res = append(res, tmp)
			tmp = ""
		}
		tmp += string(str[i])
	}

	if len(tmp) > 0 {
		res = append(res, tmp)
	}

	return res
}
