package strings

import "strings"

func StripCtlFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r != 127 {
			return r
		}
		return -1
	}, str)
}

func TruncateText(str string, max int) string {
	if len(str) <= max {
		return str
	}
	return str[:max] + "…"
}
