package xutils

import (
	"strings"

	"github.com/rainycape/unidecode"
)

func StringToArray(input string) []string {
	if input == "" || input == "[]" {
		return []string{}
	}

	items := strings.Split(input, ",")
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}
	return items
}

func StripHTML(input string) string {
	input = strings.ReplaceAll(input, "\u003C", "<")
	input = strings.ReplaceAll(input, "\u003E", ">")

	var result strings.Builder
	var inside bool
	for _, char := range input {
		if char == '<' {
			inside = true
			continue
		}
		if char == '>' {
			inside = false
			continue
		}
		if !inside {
			result.WriteRune(char)
		}
	}
	return strings.TrimSpace(result.String())
}

func RemoveDiacritics(s string) string {
	return unidecode.Unidecode(s)
}
