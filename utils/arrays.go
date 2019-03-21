package utils

// Contains : function to lookup a value in an array of strings
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
