package utils

func PadZero(str string, length int) string {
	for len(str) < length {
		str = "0" + str
	}
	return str
}
