package internal

func IsNilOrEmpty(s string) bool {
	return s == "" || s == "<nil>"
}