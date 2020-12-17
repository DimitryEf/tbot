package tools

import "strings"

//'_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'
func EscapeMarkdownV2(s string) string {
	s = strings.ReplaceAll(s, "_", "\\_")
	s = strings.ReplaceAll(s, "*", "\\*")
	s = strings.ReplaceAll(s, "[", "\\[")
	s = strings.ReplaceAll(s, "]", "\\]")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	s = strings.ReplaceAll(s, "~", "\\~")
	s = strings.ReplaceAll(s, "`", "\\`")
	s = strings.ReplaceAll(s, ">", "\\>")
	s = strings.ReplaceAll(s, "#", "\\#")
	s = strings.ReplaceAll(s, "+", "\\+")
	s = strings.ReplaceAll(s, "-", "\\-")
	s = strings.ReplaceAll(s, "=", "\\=")
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "{", "\\{")
	s = strings.ReplaceAll(s, "}", "\\}")
	s = strings.ReplaceAll(s, ".", "\\.")
	s = strings.ReplaceAll(s, "!", "\\!")
	return s
}
