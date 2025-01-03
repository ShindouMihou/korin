package kslices

// RemoveString removes the first occurrence of a string from a slice.
func RemoveString(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
