package convert

import "strconv"

// S2Int convert string to int.
func S2Int(s string, d ...int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		}

		return 0
	}

	return i
}

// I2S convert int to string.
func I2S(i int) string {
	return strconv.Itoa(i)
}
