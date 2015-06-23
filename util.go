package main

// plural returns p if num is not 1, or an empty string otherwise.
func plural(num int, p string) string {
	if num == 1 {
		return ""
	}

	return p
}
