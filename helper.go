package talkback

// sliceContainsString returns true if the slice contains the string.
func sliceContainsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}
