package ui

import (
	"bytes"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func HighlightCode(code, language string) string {
	if language == "" {
		return code
	}
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	it, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}
	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, it); err != nil {
		return code
	}
	return buf.String()
}
