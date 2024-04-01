package tinybasic

import (
	"strconv"
	"strings"
	"unicode"
)

type LineScanner struct {
	line string
	pos  int
}

func NewLineScanner(line string) *LineScanner {
	return &LineScanner{
		line: line,
		pos:  0,
	}
}

func (ls *LineScanner) GetTail() (tail string) {
	if ls.IsEOL() {
		return tail
	}

	return ls.line[ls.pos:]
}

func (ls *LineScanner) IsBOL() bool {
	return ls.pos == 0
}

func (ls *LineScanner) IsEOL() bool {
	return ls.pos >= len(ls.line)
}

func (ls *LineScanner) PeekChar() (rune, bool) {
	if ls.IsEOL() {
		return 0, false
	}
	return rune(ls.line[ls.pos]), true
}

func (ls *LineScanner) TestChar(ch rune) bool {
	if ls.IsEOL() {
		return false
	}
	return rune(ls.line[ls.pos]) == ch
}

func (ls *LineScanner) TestSpace() bool {
	if ls.IsEOL() {
		return false
	}
	return unicode.IsSpace(rune(ls.line[ls.pos]))
}

func (ls *LineScanner) TestCharNot(ch rune) bool {
	if ls.IsEOL() {
		return false
	}
	return rune(ls.line[ls.pos]) != ch
}

func (ls *LineScanner) TestChars(chars string) bool {
	if ls.IsEOL() {
		return false
	}
	return strings.ContainsRune(chars, rune(ls.line[ls.pos]))
}

func (ls *LineScanner) TestString(str string) bool {
	return strings.HasPrefix(ls.line[ls.pos:], str)
}

func (ls *LineScanner) TestStrings(strings []string) bool {
	for _, str := range strings {
		if ls.TestString(str) {
			return true
		}
	}
	return false
}

func (ls *LineScanner) TestNumber() bool {
	if ls.IsEOL() {
		return false
	}
	_, err := strconv.Atoi(string(ls.line[ls.pos]))
	return err == nil
}

func (ls *LineScanner) Shift(num int) bool {
	ls.pos += num
	if ls.pos > len(ls.line) {
		ls.pos = len(ls.line)
		return false
	}
	return true
}

func (ls *LineScanner) GetChar(ch rune) *rune {
	if ls.TestChar(ch) {
		ls.Shift(1)
		return &ch
	}
	return nil
}

func (ls *LineScanner) GetSpace() *rune {
	if ls.TestSpace() {
		space := ' '
		ls.Shift(1)
		return &space
	}
	return nil
}

func (ls *LineScanner) GetSpaces() *string {
	startPos := ls.pos

	for ls.TestSpace() {
		ls.Shift(1)
	}

	if startPos == ls.pos {
		return nil
	}

	spaces := ls.line[startPos:ls.pos]

	return &spaces
}

func (ls *LineScanner) GetChars(chars string) *rune {
	if ls.TestChars(chars) {
		ch := rune(ls.line[ls.pos])
		ls.Shift(1)
		return &ch
	}
	return nil
}

func (ls *LineScanner) GetCharNot(chars string) *rune {
	if ls.IsEOL() || strings.ContainsRune(chars, rune(ls.line[ls.pos])) {
		return nil
	}
	ch := rune(ls.line[ls.pos])
	ls.Shift(1)
	return &ch
}

func (ls *LineScanner) GetString(str string) *string {
	if ls.TestString(str) {
		startPos := ls.pos
		ls.Shift(len(str))
		line := ls.line[startPos:ls.pos]
		return &line
	}
	return nil
}

func (ls *LineScanner) GetNumber() *int {
	startPos := ls.pos
	for !ls.IsEOL() && unicode.IsDigit(rune(ls.line[ls.pos])) {
		ls.Shift(1)
	}

	if startPos == ls.pos {
		return nil
	}

	num, err := strconv.Atoi(ls.line[startPos:ls.pos])
	if err != nil {
		return nil
	}

	return &num
}

func (ls *LineScanner) GetStrings(strings []string) *string {
	for _, str := range strings {
		if ls.TestString(str) {
			ls.Shift(len(str))
			return &str
		}
	}
	return nil
}
