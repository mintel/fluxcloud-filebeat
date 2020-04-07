package utils

// StringInSlice ... Returns whether a string is in a slice or not
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
