package cmd

func checkForEmptyString(s string) (r string) {
	r = s
	if s == "" {
		r = "<empty>"
	}
	return r
}
