package reago

import (
	"strings"
)

type tplToken struct {
	IsBind bool
	Text   string
}

type TplParser struct {
	tokens []tplToken
	binds  []string
}

func NewTplParser(str string) *TplParser {
	tpl := TplParser{}
	tpl.parse(str)
	return &tpl
}

func (tpl *TplParser) parse(str string) {
	placeholderCount := strings.Count(str, "{{")
	tokens := make([]tplToken, 0, placeholderCount*2+1)
	binds := make([]string, 0, placeholderCount)

	start := 0
	for {
		open := strings.Index(str[start:], "{{")
		if open == -1 {
			// Append remaining static text.
			tokens = append(tokens, tplToken{IsBind: false, Text: str[start:]})
			break
		}
		open += start

		// Append static text before the placeholder.
		if open > start {
			tokens = append(tokens, tplToken{IsBind: false, Text: str[start:open]})
		}

		close := strings.Index(str[open:], "}}")
		if close == -1 {
			// If no closing delimiter, treat the rest as static text.
			tokens = append(tokens, tplToken{IsBind: false, Text: str[open:]})
			break
		}
		close += open

		// Extract key between delimiters, trimming whitespace.
		key := strings.TrimSpace(str[open+2 : close])
		tokens = append(tokens, tplToken{IsBind: true, Text: key})
		binds = append(binds, key)

		start = close + 2
	}

	tpl.tokens = tokens
	tpl.binds = binds
}

func (tpl *TplParser) GetBinds() []string {
	return tpl.binds
}

func (tpl *TplParser) Render(state *State) string {
	var builder strings.Builder

	for _, token := range tpl.tokens {
		if token.IsBind {
			builder.WriteString(state.GetString(token.Text).Get())
		} else {
			builder.WriteString(token.Text)
		}
	}

	return builder.String()
}
