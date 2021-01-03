package model

type Resp struct {
	Text     string
	ParseMod string
}

const (
	ModeMarkdown   = "Markdown"
	ModeMarkdownV2 = "MarkdownV2"
	ModeHTML       = "HTML"
)
