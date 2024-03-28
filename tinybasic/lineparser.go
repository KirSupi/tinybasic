package tinybasic

import (
	"strings"
	"unicode"
)

type LineParser struct {
	scanner *LineScanner
}

func NewLineParser(line string) *LineParser {
	return &LineParser{
		scanner: NewLineScanner(line),
	}
}

func NewLineParserWithScanner(s *LineScanner) *LineParser {
	return &LineParser{
		scanner: s,
	}
}

func (lp *LineParser) GetVariable() *string {
	for !lp.scanner.IsEOL() && unicode.IsSpace(rune(lp.scanner.line[lp.scanner.pos])) {
		lp.scanner.Shift(1)
	}

	if lp.scanner.IsEOL() || !unicode.IsLetter(rune(lp.scanner.line[lp.scanner.pos])) {
		return nil
	}

	ch := rune(lp.scanner.line[lp.scanner.pos])
	if ch < 'A' || ch > 'Z' && ch < 'a' || ch > 'z' {
		return nil
	}

	varName := strings.ToUpper(string(ch))
	lp.scanner.Shift(1)
	return &varName
}
func (lp *LineParser) GetQuotedString() *string {
	for !lp.scanner.IsEOL() && unicode.IsSpace(rune(lp.scanner.line[lp.scanner.pos])) {
		lp.scanner.Shift(1)
	}

	if lp.scanner.IsEOL() || lp.scanner.line[lp.scanner.pos] != '"' {
		return nil
	}

	lp.scanner.Shift(1) // Пропуск открывающей кавычки
	startPos := lp.scanner.pos

	for !lp.scanner.IsEOL() && lp.scanner.line[lp.scanner.pos] != '"' {
		lp.scanner.Shift(1)
	}

	if lp.scanner.IsEOL() {
		panic("Quoted string not terminated")
	}

	quotedString := lp.scanner.line[startPos:lp.scanner.pos]
	lp.scanner.Shift(1) // Пропуск закрывающей кавычки
	return &quotedString
}
